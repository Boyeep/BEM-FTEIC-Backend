package config

import (
	"fmt"
	"os"
	"strings"
)

func Required(name string) (string, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return "", fmt.Errorf("required environment variable %q is not set", name)
	}
	return value, nil
}

func Optional(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}
