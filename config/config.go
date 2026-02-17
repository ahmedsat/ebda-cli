package config

import (
	"fmt"
	"os"
)

var (
	ErpUsername    string
	ErpPassword    string
	ErpBaseUrl     string
	KoboAuthToken  string
	KoboBaseURL    string
	DBFilePath     string
	MigrationsList []any
	DisableNotify  bool
)

func Configure() (err error) {
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

	DBFilePath, ok = os.LookupEnv("DB_EBDA_CLI_FILE_PATH")
	if !ok {
		return fmt.Errorf("DB_EBDA_CLI_FILE_PATH is not set")
	}

	disableNotify, ok := os.LookupEnv("DISABLE_NOTIFY")
	if !ok {
		return fmt.Errorf("DISABLE_NOTIFY is not set")
	}
	DisableNotify = disableNotify == "true"

	err = initDB()
	if err != nil {
		return
	}

	return nil
}
