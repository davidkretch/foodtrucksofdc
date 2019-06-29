package main

import (
	"log"

	"foodtrucks/dcgov/loaddb"
)

func main() {
	file := "Apr 2017 - MRV Lottery Results.csv"
	bucket := "davidkretch-test"
	project := "serene-foundry-234813"
	err := loaddb.LoadDB(file, bucket, project)
	if err != nil {
		log.Fatalf("Error processing file %s: %s", file, err)
	}
}
