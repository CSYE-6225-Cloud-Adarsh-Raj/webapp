package main

import (
	"time"
	"webapp/db"
	"webapp/router"
	"webapp/setup"

	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

var dbConnection *gorm.DB
var err error

var logger = logrus.New()

func main() {
	config := setup.GetDBConfig()
	dbConnection, err = db.ConnectToDB(config.DSN)
	if err != nil {
		logger.Error("main() - Database connection failed. Starting retry sequence...")
		go func() {
			for {
				time.Sleep(3 * time.Second)
				dbConnection, err = db.ConnectToDB(config.DSN)
				if err == nil {
					// log.Println("Database connection established")
					logger.Error("main() - Database connection established")
					break
				}
				logger.Error("main() - Retrying to connect to database...")
			}
		}()
	}

	r := router.InitRouter(dbConnection)
	r.Run()
}
