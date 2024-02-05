package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"webapp/api/health"
	"webapp/api/user"
)

func AuthenticationMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Basic authentication required"})
			return
		}

		if !user.ValidateCredentials(db, username, password) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		var user user.UserModel
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID details"})
			return
		}

		c.Set("username", username)
		c.Set("userID", user.ID)
		c.Next()
	}
}

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

	authGroup := r.Group("/")
	authGroup.Use(AuthenticationMiddleware(db))
	{
		//Get User
		authGroup.GET("/v1/user/self", user.GetUserDetails(db))
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{})
	})

	return r
}
