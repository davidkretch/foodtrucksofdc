package loaddb

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
)

var months = map[string]int{
	"jan": 1,
	"feb": 2,
	"mar": 3,
	"apr": 4,
	"may": 5,
	"jun": 6,
	"jul": 7,
	"aug": 8,
	"sep": 9,
	"oct": 10,
	"nov": 11,
	"dec": 12,
}

// GetMonthAndYear returns the month and year from a string
// of the form "Jan 2006...".
func GetMonthAndYear(s string) (time.Month, int, error) {
	el := strings.Split(s, " ")
	month, ok := months[strings.ToLower(el[0][0:3])]
	if !ok {
		return 0, 0, errors.New("Cannot convert month to int")
	}
	year, err := strconv.Atoi(el[1])
	if err != nil {
		return 0, 0, errors.New("Cannot convert year to int")
	}
	return time.Month(month), year, nil
}

// GetFile returns an array of bytes for `file` in `bucket`.
func GetFile(file string, bucket string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	rc, err := client.Bucket(bucket).Object(file).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ReadCSV returns a slice of maps corresponding to the rows
// in a CSV provided in `data`.
func ReadCSV(data io.Reader) ([]map[string]string, error) {
	r := csv.NewReader(data)
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

// ProcessData returns a nested map with:
// Each day of a given month
//   -> The stops being served
//      -> Which trucks serve the stop on the given day
func ProcessData(data []map[string]string, month time.Month, year int) (map[string]map[string][]string, error) {
	days := make(map[string]map[string][]string)
	for d := time.Sunday; d <= time.Saturday; d++ {
		day := d.String()
		days[day] = make(map[string][]string)
	}
	for _, rec := range data {
		truck := rec["Business Name"]
		for d := time.Sunday; d <= time.Saturday; d++ {
			day := d.String()
			if stop, ok := rec[day]; ok {
				if stop == "OFF" {
					continue
				}
				if trucks, ok := days[day][stop]; ok {
					days[day][stop] = append(trucks, truck)
				} else {
					days[day][stop] = []string{truck}
				}
			}
		}
	}
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	dates := make(map[string]map[string][]string)
	for d := start; d.Month() == start.Month(); d = d.AddDate(0, 0, 1) {
		date := d.Format("2006-01-02")
		weekday := d.Weekday().String()
		dates[date] = days[weekday]
	}
	return dates, nil
}

// UploadData uploads a dataset to the database.
func UploadData(data map[string]map[string][]string, project string) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, project)
	if err != nil {
		return err
	}
	defer client.Close()
	coll := client.Collection("dates")
	for date, stops := range data {
		_, err := coll.Doc(date).Set(ctx, stops)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadDB extracts a month's data from a CSV, transforms it into one
// observation per day, and then loads it into the database.
func LoadDB(name string, bucket string, project string) error {
	if ext := filepath.Ext(name); ext != ".csv" {
		return nil
	}
	month, year, err := GetMonthAndYear(name)
	if err != nil {
		return err
	}
	file, err := GetFile(name, bucket)
	if err != nil {
		return err
	}
	data, err := ReadCSV(bytes.NewReader(file))
	if err != nil {
		return err
	}
	processed, err := ProcessData(data, month, year)
	if err != nil {
		return err
	}
	err = UploadData(processed, project)
	if err != nil {
		return err
	}
	return nil
}
