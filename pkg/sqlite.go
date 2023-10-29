package pkg

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetSqliteDb(file string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error open sqlite db: %s \n", err)
	}
	return db, nil
}
