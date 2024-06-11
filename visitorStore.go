package main

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createTable(dbFile string) error {
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) { // create the table
		db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Visitor{})
		if err != nil {
			return err
		}
	}
	return nil
}

func saveVisitor(v Visitor) error {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return err
	}
	result := db.Create(&v)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func listVisitors(dbFile string) ([]Visitor, error) {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	var visitors []Visitor
	result := db.Find(&visitors)
	if result.Error != nil {
		return nil, result.Error
	}
	return visitors, nil
}
