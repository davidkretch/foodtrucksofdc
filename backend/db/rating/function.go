// Package p contains a Cloud Function triggered by a Firestore event.
package p

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	// Fields is the data for this value. The type depends on the format of your
	// database. Log the interface{} value and inspect the result to see a JSON
	// representation of your database fields.
	Fields     Rating    `json:"fields"`
	Name       string    `json:"name"`
	UpdateTime time.Time `json:"updateTime"`
}

// Rating holds a single user rating.
type Rating struct {
	Rating struct {
		IntegerValue string `json:"integerValue"`
	} `json:"rating"`
}

// Truck holds aggregate rating info for a truck.
type Truck struct {
	AvgRating  float64 `firestore:"avgRating"`
	NumRatings int     `firestore:"numRatings"`
}

// GCLOUD_PROJECT is automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("GCLOUD_PROJECT")

// client is a Firestore client, reused between function invocations.
// Source: https://cloud.google.com/functions/docs/calling/cloud-firestore
var client *firestore.Client

func init() {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
}

// SetAvgRating updates a truck's average rating whenever a user
// enters a new rating.
func SetAvgRating(ctx context.Context, e FirestoreEvent) error {
	path := strings.Split(e.Value.Name, "/documents/")[1]
	truckName := strings.Split(path, "/")[1]

	newRating, err := strconv.ParseFloat(e.Value.Fields.Rating.IntegerValue, 64)
	if err != nil {
		return err
	}
	oldRating, err := strconv.ParseFloat(e.OldValue.Fields.Rating.IntegerValue, 64)
	newUserRating := false
	if err != nil {
		newUserRating = true
	}

	truck := client.Doc("trucks/" + truckName)
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(truck)
		if err != nil && grpc.Code(err) != codes.NotFound {
			return err
		}

		var newNumRatings float64
		var newAvgRating float64

		if doc.Exists() {
			var data Truck
			err = doc.DataTo(&data)
			if err != nil {
				return err
			}
			oldAvgRating := data.AvgRating
			oldNumRatings := float64(data.NumRatings)
			oldSumRatings := oldAvgRating * oldNumRatings

			if newUserRating {
				newNumRatings = oldNumRatings + 1
				newAvgRating = (oldSumRatings + newRating) / newNumRatings
			} else {
				newNumRatings = oldNumRatings
				newAvgRating = (oldSumRatings - oldRating + newRating) / newNumRatings
			}
		} else {
			newNumRatings = 1
			newAvgRating = newRating
		}

		return tx.Set(truck, map[string]interface{}{
			"avgRating":  newAvgRating,
			"numRatings": int(newNumRatings),
		}, firestore.MergeAll)
	})
	if err != nil {
		return err
	}
	return nil
}
