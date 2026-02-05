package config

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErpUsername    string
	ErpPassword    string
	ErpBaseUrl     string
	KoboAuthToken  string
	KoboBaseURL    string
	DBFilePath     string
	DB             *gorm.DB
	MigrationsList []any
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

	if DBFilePath != "" {

		fmt.Fprintln(os.Stderr, DBFilePath)
		DB, err = gorm.Open(sqlite.Open(DBFilePath), &gorm.Config{})
		if err != nil {
			return
		}

		// disable logging
		DB.Logger = logger.Discard
		err = DB.AutoMigrate(MigrationsList...)
		if err != nil {
			return
		}
	}

	return nil
}
