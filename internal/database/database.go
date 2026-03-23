// Package database provides GORM database setup and migrations.
package database

import (
	"calculator-api/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDB opens a SQLite database at dbPath and runs auto-migrations.
func NewDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run auto-migrations to create/update the schema.
	if err := db.AutoMigrate(&model.Calculation{}); err != nil {
		return nil, err
	}

	return db, nil
}
