package db

import (
	"fmt"
	"webapp/api/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DSN string
}

func ConnectToDB(DSN string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("Not able to connect to the Postgres Database")
	}

	err = db.AutoMigrate(&user.UserModel{})
	if err != nil {
		fmt.Println("Failed to migrate table schema")
		return nil, err
	}
	return db, err
}
