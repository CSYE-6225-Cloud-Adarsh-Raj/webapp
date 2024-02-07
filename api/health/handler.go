package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HealthCheckHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.ContentLength > 0 {
			c.Status(http.StatusBadRequest)
			return
		}

		if len(c.Request.URL.RawQuery) > 0 {
			c.Status(http.StatusBadRequest)
			return
		}

		// c.Header("Cache-Control", "no-cache")

		postgresDB, err := db.DB()
		if err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		if err := postgresDB.Ping(); err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		c.Status(http.StatusOK)
	}
}
