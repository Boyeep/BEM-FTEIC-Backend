package models

import (
	"fmt"
	"regexp"
	"time"
)

type UUID string

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-8][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

func ParseUUID(value string) (UUID, error) {
	if !uuidPattern.MatchString(value) {
		return "", fmt.Errorf("invalid UUID")
	}
	return UUID(value), nil
}

func (id UUID) String() string { return string(id) }

type Date string

func ParseDate(value string) (Date, error) {
	if _, err := time.Parse(time.DateOnly, value); err != nil {
		return "", fmt.Errorf("invalid date: %w", err)
	}
	return Date(value), nil
}

func (date Date) String() string { return string(date) }
