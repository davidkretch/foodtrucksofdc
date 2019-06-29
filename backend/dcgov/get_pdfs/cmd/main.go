package main

import (
	"log"

	"foodtrucks/dcgov/getpdfs"
)

func main() {
	err := getpdfs.GetPDFs("https://dcra.dc.gov/mrv", "davidkretch-test")
	if err != nil {
		log.Fatalf("%s", err)
	}
}
