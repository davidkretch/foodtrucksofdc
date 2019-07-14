package loaddb

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
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
	re := regexp.MustCompile(`^(.+) ?(\d{4})`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 3 {
		return time.Month(0), 0, errors.New("Invalid string")
	}
	m, y := matches[1], matches[2]
	month, ok := months[strings.ToLower(m[0:3])]
	if !ok {
		return 0, 0, errors.New("Cannot convert month to int")
	}
	year, err := strconv.Atoi(y)
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

// Records holds a slice of key/value pairs representing the records in a CSV.
type Records = []map[string]string

// ReadCSV returns a slice of maps corresponding to the rows
// in a CSV provided in `data`.
func ReadCSV(data io.Reader) (Records, error) {
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

// CheckData returns whether the data has the expected columns.
func CheckData(data Records) bool {
	rec := data[0]
	expected := []string{"Business Name"}
	for d := time.Monday; d <= time.Friday; d++ {
		expected = append(expected, d.String())
	}
	for _, col := range expected {
		if _, ok := rec[col]; !ok {
			return false
		}
	}
	return true
}

// Set holds distinct string values, e.g. a list of distinct trucks.
type Set = map[string]bool

// DailySchedule holds all stops and their trucks for a day.
type DailySchedule = map[string][]string

// MonthlySchedule holds the trucks and daily schedules for a month.
type MonthlySchedule struct {
	Trucks Set
	Days   map[string]DailySchedule
}

// Process returns a MonthlySchedule from data taken from a CSV.
func Process(data Records, month time.Month, year int) (MonthlySchedule, error) {
	if ok := CheckData(data); !ok {
		return MonthlySchedule{}, errors.New("Data in wrong format")
	}

	days := make(map[string]map[string][]string)
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	for d := start; d.Month() == start.Month(); d = d.AddDate(0, 0, 1) {
		date := d.Format("2006-01-02")
		days[date] = make(map[string][]string)
	}

	trucks := make(map[string]bool)

	for _, rec := range data {
		truck := rec["Business Name"]
		if truck == "" {
			continue
		}
		trucks[truck] = true

		for d := start; d.Month() == start.Month(); d = d.AddDate(0, 0, 1) {
			date := d.Format("2006-01-02")
			weekday := d.Weekday().String()
			if stop, ok := rec[weekday]; ok {
				if stop == "OFF" {
					continue
				}
				if t, ok := days[date][stop]; ok {
					days[date][stop] = append(t, truck)
				} else {
					days[date][stop] = []string{truck}
				}
			}
		}
	}
	result := MonthlySchedule{
		Trucks: trucks,
		Days:   days,
	}
	return result, nil
}

// KeyName returns a string suitable for use as a Firestore document name.
func KeyName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return strings.ToLower(re.ReplaceAllString(name, ""))
}

// ID holds an ID corresponding to a name, retrieved from the database.
type ID struct {
	ID string `firestoreValue:"id"`
}

// GetExistingTruckIDs returns all existing truck IDs in the database.
func GetExistingTruckIDs(ctx context.Context, client *firestore.Client) (map[string]string, error) {
	truckIDs := make(map[string]string)
	docs, err := client.Collection("truckNames").Documents(ctx).GetAll()
	if err != nil {
		return map[string]string{}, err
	}
	for _, doc := range docs {
		key := doc.Ref.ID
		var id ID
		err = doc.DataTo(&id)
		if err != nil {
			return map[string]string{}, err
		}
		truckIDs[key] = id.ID
	}
	return truckIDs, nil
}

// AddTruckID adds a new truck to the database and returns its ID.
func AddTruckID(batch *firestore.WriteBatch, truck string, client *firestore.Client) string {
	truckRef := client.Collection("trucks").NewDoc()
	truckID := truckRef.ID
	nameRef := client.Collection("truckNames").Doc(KeyName(truck))
	batch.Set(truckRef, map[string]string{"displayName": truck})
	batch.Set(nameRef, map[string]string{"id": truckID})
	return truckID
}

// GetTruckIDs returns a truck's ID, creating one if it does not exist.
func GetTruckIDs(ctx context.Context, trucks Set, client *firestore.Client) (map[string]string, error) {
	truckIDs, err := GetExistingTruckIDs(ctx, client)
	if err != nil {
		return map[string]string{}, err
	}
	batch := client.Batch()
	newTrucks := 0
	for truck := range trucks {
		if truckID, ok := truckIDs[KeyName(truck)]; !ok {
			newTrucks = newTrucks + 1
			truckID = AddTruckID(batch, truck, client)
			truckIDs[KeyName(truck)] = truckID
		}
	}
	if newTrucks > 0 {
		_, err = batch.Commit(ctx)
		if err != nil {
			return map[string]string{}, err
		}
	}
	return truckIDs, nil
}

// Upload uploads a dataset to the database.
func Upload(schedule MonthlySchedule, project string, file string) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, project)
	if err != nil {
		return err
	}
	defer client.Close()

	truckIDs, err := GetTruckIDs(ctx, schedule.Trucks, client)
	if err != nil {
		return err
	}
	batch := client.Batch()
	for date, stops := range schedule.Days {
		docRef := client.Collection("schedules").Doc(date)
		data := make(map[string][]map[string]string)
		for stop, trucks := range stops {
			data[stop] = []map[string]string{}
			for _, truck := range trucks {
				entry := map[string]string{
					"id": truckIDs[KeyName(truck)],
				}
				data[stop] = append(data[stop], entry)
			}
		}
		batch.Set(docRef, data)
	}
	_, err = batch.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// SetFileStatus sets a file ok or not ok in the database,
// based on whether there is an error.
func SetFileStatus(name string, project string, status error) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, project)
	if err != nil {
		return err
	}
	defer client.Close()
	fileNoExt := strings.TrimSuffix(name, path.Ext(name))
	fileRef := client.Collection("dcGovFiles").Doc(fileNoExt)
	ok := true
	if status != nil {
		ok = false
	}
	fileRef.Set(ctx, map[string]bool{"ok": ok})
	return nil
}

// LoadDB extracts a month's data from a CSV, transforms it into one
// observation per day, and then loads it into the database.
func LoadDB(name string, bucket string, project string) (err error) {
	defer func() {
		SetFileStatus(name, project, err)
	}()
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
	processed, err := Process(data, month, year)
	if err != nil {
		return err
	}
	err = Upload(processed, project, name)
	if err != nil {
		return err
	}
	return nil
}
