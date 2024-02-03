package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DSN string
}

func ConnectToDB(DSN string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	return db, err
}
