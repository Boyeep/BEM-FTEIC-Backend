package models

import (
	"testing"
	"time"
)

func TestDateDatabaseContract(t *testing.T) {
	var date Date
	if err := date.Scan(time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	if date.String() != "2026-08-01" {
		t.Fatalf("unexpected scanned date %q", date)
	}
	value, err := date.Value()
	if err != nil {
		t.Fatal(err)
	}
	if value != "2026-08-01" {
		t.Fatalf("unexpected database value %#v", value)
	}
}

func TestDateRejectsInvalidDatabaseValues(t *testing.T) {
	var date Date
	if err := date.Scan("not-a-date"); err == nil {
		t.Fatal("expected invalid string date to fail")
	}
	if _, err := Date("not-a-date").Value(); err == nil {
		t.Fatal("expected invalid date value to fail")
	}
}
