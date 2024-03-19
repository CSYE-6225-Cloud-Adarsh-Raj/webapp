package main

import (
	"time"
	"webapp/db"
	"webapp/router"
	"webapp/setup"

	// "cloud.google.com/go/logging"
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
		// log.Println("Database connection failed. Starting retry sequence...")
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
				// log.Println("Retrying to connect to database...")
				logger.Error("main() - Retrying to connect to database...")
			}
		}()
	}

	r := router.InitRouter(dbConnection)
	r.Run()

	// // Log an example message
	// log.WithFields(logrus.Fields{
	// 	"event": "event_name",
	// 	"topic": "topic_name",
	// 	"key":   "some_value",
	// }).Info("A structured logging message")

}

// func createLogger(projectID, logID string) (*logging.Logger, error) {
// 	ctx := context.Background()
// 	client, err := logging.NewClient(ctx, projectID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return client.Logger(logID), nil
// }
