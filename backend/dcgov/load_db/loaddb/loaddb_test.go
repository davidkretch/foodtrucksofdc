package loaddb

import (
	"strings"
	"testing"
	"time"
)

func TestGetMonthAndYear(t *testing.T) {
	month, year, err := GetMonthAndYear("Jul 2019 foo bar")
	if err != nil {
		t.Fatalf("GetMonthAndYear returned error: %v", err)
	}
	if month != time.July || year != 2019 {
		t.Fatal("Month or year wrong")
	}

	month, year, err = GetMonthAndYear("December2019foobar")
	if err != nil {
		t.Fatalf("GetMonthAndYear returned error: %v", err)
	}
	if month != time.December || year != 2019 {
		t.Fatal("Month or year wrong")
	}

	month, year, err = GetMonthAndYear("FooBar")
	if err == nil {
		t.Fatal("GetMonthAndYear failed to return an error on invalid data")
	}
}

func TestReadCSV(t *testing.T) {
	data := "a,b,c\n1,2,3\n4,5,6"
	recs, err := ReadCSV(strings.NewReader(data))
	if err != nil {
		t.Fatalf("ReadCSV returned error: %v", err)
	}
	if len(recs) != 2 {
		t.Fatal("ReadCSV returned the wrong number of rows")
	}
	if rec := recs[0]; rec["a"] != "1" || rec["b"] != "2" || rec["c"] != "3" {
		t.Fatal("ReadCSV returned incorrect data")
	}
}

func TestCheckData(t *testing.T) {
	recs1 := make([]map[string]string, 1)
	recs1[0] = map[string]string{
		"Business Name": "foo",
		"Monday":        "foo",
		"Tuesday":       "foo",
		"Wednesday":     "foo",
		"Thursday":      "foo",
		"Friday":        "foo",
	}
	ok := CheckData(recs1)
	if !ok {
		t.Fatal("CheckData returned not ok on valid data")
	}

	recs2 := make([]map[string]string, 2)
	recs2[0] = map[string]string{"a": "1", "b": "2", "c": "3"}
	recs2[1] = map[string]string{"a": "4", "b": "5", "c": "6"}
	ok = CheckData(recs2)
	if ok {
		t.Fatal("CheckData returned ok on invalid data")
	}
}
