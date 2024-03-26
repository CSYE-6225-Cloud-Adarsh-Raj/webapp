package health

import (
	"net/http"
	"webapp/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HealthCheckHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.ContentLength > 0 {
			logger.Logger.Error("HealthCheckHandler() - Bad Request")
			c.Status(http.StatusBadRequest)
			return
		}

		if len(c.Request.URL.RawQuery) > 0 {
			logger.Logger.Error("HealthCheckHandler() - Bad Request")
			c.Status(http.StatusBadRequest)
			return
		}

		// c.Header("Cache-Control", "no-cache")

		postgresDB, err := db.DB()
		if err != nil {
			logger.Logger.Error("HealthCheckHandler() - ServiceUnavailable, cannot connect to db")
			c.Status(http.StatusServiceUnavailable)
			return
		}

		if err := postgresDB.Ping(); err != nil {
			logger.Logger.Error("HealthCheckHandler() - ServiceUnavailable, cannot ping db")
			c.Status(http.StatusServiceUnavailable)
			return
		}

		logger.Logger.Info("HealthCheckHandler() - Database health is OK")
		c.Status(http.StatusOK)
	}
}
