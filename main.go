package main

import (
	"CSYE_6225_01/db"
	"CSYE_6225_01/router"
	"CSYE_6225_01/setup"
	"log"
	"time"

	"gorm.io/gorm"
)

var dbConnection *gorm.DB
var err error

func main() {
	config := setup.GetDBConfig()
	dbConnection, err = db.ConnectToDB(config.DSN)
	if err != nil {
		log.Println("Database connection failed. Starting retry sequence...")
		go func() {
			for {
				time.Sleep(3 * time.Second)
				dbConnection, err = db.ConnectToDB(config.DSN)
				if err == nil {
					log.Println("Database connection established")
					break
				}
				log.Println("Retrying to connect to database...")
			}
		}()
	}

	r := router.InitRouter(dbConnection)
	r.Run()
}
