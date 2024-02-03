package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"CSYE_6225_01/api/health"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// httpMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	// for _, method := range httpMethods {
	// 	r.Handle(method, "/healthz", health.HealthCheckHandler(db))
	// }

	r.Use(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.URL.Path == "/healthz" {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}

	})

	r.GET("/healthz", health.HealthCheckHandler(db))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{})
	})

	return r
}
