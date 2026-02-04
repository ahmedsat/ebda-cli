package config

import (
	"fmt"
	"os"
)

var (
	ErpUsername   string
	ErpPassword   string
	ErpBaseUrl    string
	KoboAuthToken string
	KoboBaseURL   string
)

func Configure() error {
	var ok bool
	ErpUsername, ok = os.LookupEnv("ERP_USERNAME")
	if !ok {
		return fmt.Errorf("ERP_USERNAME is not set")
	}
	ErpPassword, ok = os.LookupEnv("ERP_PASSWORD")
	if !ok {
		return fmt.Errorf("ERP_PASSWORD is not set")
	}
	ErpBaseUrl, ok = os.LookupEnv("ERP_BASE_URL")
	if !ok {
		return fmt.Errorf("ERP_BASE_URL is not set")
	}

	KoboAuthToken, ok = os.LookupEnv("KOBO_AUTH_TOKEN")
	if !ok {
		return fmt.Errorf("KOBO_AUTH_TOKEN is not set")
	}
	KoboBaseURL, ok = os.LookupEnv("KOBO_BASE_URL")
	if !ok {
		return fmt.Errorf("KOBO_BASE_URL is not set")
	}

	return nil
}
