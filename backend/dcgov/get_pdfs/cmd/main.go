package main

import (
	"log"

	"foodtrucks/dcgov/getpdfs"
)

func main() {
	url := "https://dcra.dc.gov/mrv"
	bucket := "davidkretch-test"
	project := "serene-foundry-234813"
	err := getpdfs.GetPDFs(url, bucket, project)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
