// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"os"

	"foodtrucks/dcgov/getpdfs"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// GetPDFs gets new PDFs from the given URL and stores them in a bucket.
func GetPDFs(ctx context.Context, m PubSubMessage) error {
	url := os.Getenv("URL")
	bucket := os.Getenv("BUCKET")
	project := os.Getenv("PROJECT")
	err := getpdfs.GetPDFs(url, bucket, project)
	return err
}
