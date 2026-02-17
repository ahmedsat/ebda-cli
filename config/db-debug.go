//go:build !release

package config

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func initDB() (err error) {
	if DBFilePath != "" {

		fmt.Fprintln(os.Stderr, DBFilePath)
		DB, err = gorm.Open(sqlite.Open(DBFilePath), &gorm.Config{})
		if err != nil {
			return err
		}

		// disable logging
		DB.Logger = logger.Discard
		err = DB.AutoMigrate(MigrationsList...)
		if err != nil {
			return err
		}
	}

	return nil
}
