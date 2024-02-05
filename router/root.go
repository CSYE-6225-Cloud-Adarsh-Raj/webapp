package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"webapp/api/health"
	"webapp/api/user"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.URL.Path == "/healthz" {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}

	})

	//Check DB health
	r.GET("/healthz", health.HealthCheckHandler(db))

	//Create User
	r.POST("/v1/user", user.CreateUserHandler(db))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{})
	})

	return r
}
