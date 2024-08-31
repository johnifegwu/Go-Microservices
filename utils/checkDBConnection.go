package utils

import (
	"fmt"

	"gorm.io/gorm"
)

func checkDBConnection(db *gorm.DB) error {
	sqlDB, err := db.DB() // Get the underlying sql.DB object
	if err != nil {
		return fmt.Errorf("failed to get database object: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("database connection lost: %v", err)
	}

	fmt.Println("Database is connected")
	return nil
}
