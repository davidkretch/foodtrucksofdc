// Package p contains a Google Cloud Storage Cloud Function.
package p

import (
	"context"
	"log"
	"os"

	"foodtrucks/dcgov/loaddb"
)

// GCSEvent is the payload of a GCS event. Please refer to the docs for
// additional information regarding GCS events.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// LoadDB loads a CSV into the database upon being written to Cloud Storage.
func LoadDB(ctx context.Context, e GCSEvent) error {
	log.Printf("Processing file: %s", e.Name)
	project := os.Getenv("PROJECT")
	err := loaddb.LoadDB(e.Name, e.Bucket, project)
	if err != nil {
		return err
	}
	return nil
}
