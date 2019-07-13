package main

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// readCSV returns a slice of maps corresponding to the rows
// in a CSV provided in `data`.
func readCSV(data io.Reader) ([]map[string]string, error) {
	r := csv.NewReader(data)
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	var records []map[string]string
	fields := lines[0]
	for _, line := range lines[1:] {
		rec := make(map[string]string)
		for i := 0; i < len(fields); i++ {
			rec[fields[i]] = line[i]
		}
		records = append(records, rec)

	}
	return records, nil
}

// Truck holds information about a food truck.
type Truck struct {
	DisplayName string
	Names       []string
	Twitter     string
}

// keyName returns a string suitable for use as a Firestore document name.
func keyName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return strings.ToLower(re.ReplaceAllString(name, ""))
}

// readTrucks takes a CSV and returns an array of Trucks.
func readTrucks(file io.Reader) ([]Truck, error) {
	recs, err := readCSV(file)
	if err != nil {
		return []Truck{}, err
	}
	trucks := make(map[string]Truck)
	for _, rec := range recs {
		name := rec["display_name"]
		if truck, ok := trucks[name]; ok {
			truck.Names = append(truck.Names, keyName(rec["business_name"]))
			trucks[name] = truck
		} else {
			trucks[name] = Truck{
				DisplayName: rec["display_name"],
				Names: []string{
					keyName(rec["display_name"]),
					keyName(rec["business_name"]),
				},
				Twitter: rec["twitter"],
			}
		}
	}
	var result []Truck
	for _, v := range trucks {
		result = append(result, v)
	}
	return result, nil
}

// Name holds the truck ID corresponding to a truck name.
type Name struct {
	ID string `firestore:"id"`
}

// getTruckID returns the ID of a food truck given a name,
// or an empty string if the truck has no existing ID.
func getTruckID(ctx context.Context, names []string, client *firestore.Client) (string, error) {
	for _, name := range names {
		if name == "" {
			continue
		}
		doc, err := client.Collection("truckNames").Doc(name).Get(ctx)
		if err != nil && grpc.Code(err) != codes.NotFound {
			return "", err
		}
		if doc.Exists() {
			var name Name
			doc.DataTo(&name)
			return name.ID, nil
		}
	}
	return "", nil
}

// uploadTruck adds or updates a truck in the database,
// and adds or updates truck name to truck ID mappings.
func uploadTruck(ctx context.Context, truck Truck, client *firestore.Client) error {
	data := map[string]interface{}{
		"displayName": truck.DisplayName,
		"twitter":     truck.Twitter,
	}
	id, err := getTruckID(ctx, truck.Names, client)
	if err != nil {
		return err
	}

	trucksRef := client.Collection("trucks")
	batch := client.Batch()
	if id == "" {
		// New truck
		// TODO: Check for and handle ID collisions.
		docRef := trucksRef.NewDoc()
		batch.Set(docRef, data, firestore.MergeAll)
		id = docRef.ID
	} else {
		// Existing truck
		// TODO: Avoid overwriting data, e.g. Twitter, with empty strings.
		docRef := trucksRef.Doc(id)
		batch.Set(docRef, data, firestore.MergeAll)
	}

	namesRef := client.Collection("truckNames")
	for _, name := range truck.Names {
		if name == "" {
			continue
		}
		docRef := namesRef.Doc(name)
		batch.Set(docRef, map[string]string{
			"id": id,
		})
	}
	_, err = batch.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

// uploadTrucks uploads a collection of trucks to the database.
func uploadTrucks(trucks []Truck, project string) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, project)
	if err != nil {
		return err
	}
	defer client.Close()

	for _, truck := range trucks {
		log.Print(truck.DisplayName)
		err = uploadTruck(ctx, truck, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	file, err := os.Open("trucks.csv")
	if err != nil {
		log.Fatal(err)
	}
	trucks, err := readTrucks(file)
	if err != nil {
		log.Fatal(err)
	}
	project := os.Getenv("PROJECT")
	err = uploadTrucks(trucks, project)
	if err != nil {
		log.Fatal(err)
	}
}
