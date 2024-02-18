package setup

import (
	"fmt"
	"os"
)

type DBConfig struct {
	DSN string
}

func GetDBConfig() DBConfig {
	//GORM
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := "host=" + dbHost + " user=" + dbUser + " dbname=" + dbName + " sslmode=disable password=" + dbPassword
	fmt.Println("DSN:", dsn)

	return DBConfig{DSN: dsn}
}
