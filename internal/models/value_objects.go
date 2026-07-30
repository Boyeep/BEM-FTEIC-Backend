package models

import (
	"database/sql/driver"
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

func (date *Date) Scan(value any) error {
	switch raw := value.(type) {
	case nil:
		*date = ""
		return nil
	case time.Time:
		*date = Date(raw.Format(time.DateOnly))
		return nil
	case string:
		parsed, err := ParseDate(raw)
		if err != nil {
			return err
		}
		*date = parsed
		return nil
	case []byte:
		return date.Scan(string(raw))
	default:
		return fmt.Errorf("scan date: unsupported value type %T", value)
	}
}

func (date Date) Value() (driver.Value, error) {
	if date == "" {
		return nil, nil
	}
	parsed, err := ParseDate(string(date))
	if err != nil {
		return nil, err
	}
	return parsed.String(), nil
}
