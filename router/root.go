package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"webapp/api/health"
	"webapp/api/user"
	"webapp/logger"
)

func AuthenticationMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			// fmt.Println("AuthenticationMiddleware() - Error: Basic authentication required")
			logger.Logger.Error("AuthenticationMiddleware() - Basic authentication required")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}

		if !user.ValidateCredentials(db, username, password) {
			// fmt.Println("AuthenticationMiddleware() - Error: Invalid credentials")
			logger.Logger.Error("AuthenticationMiddleware() - Invalid credentials")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}

		var user user.UserModel
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			// fmt.Println("AuthenticationMiddleware() - Error: Failed to retrieve user ID details")
			logger.Logger.Error("AuthenticationMiddleware() - Failed to retrieve user ID details")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.Set("username", username)
		c.Set("userID", user.ID)
		c.Next()
	}
}

func CacheControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		c.Next()
	}
}

func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(CacheControlMiddleware())

	r.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		authHeader := c.GetHeader("Authorization")

		nonAuthEndpoints := map[string]bool{
			"/healthz": true,
			"/v1/user": true,
			"/verify":  true,
		}

		if nonAuthEndpoints[path] && authHeader != "" {
			// fmt.Println("Error: Non-authenticated endpoint should not include Authorization header")
			logger.Logger.Error("InitRouter() -Non-authenticated endpoint should not include Authorization header")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		allowedMethods := map[string][]string{
			"/healthz":      {"GET"},
			"/v1/user":      {"POST"},
			"/v1/user/self": {"GET", "PUT"},
			"/verify":       {"GET"},
		}

		if methods, exists := allowedMethods[path]; exists {
			methodAllowed := false
			for _, m := range methods {
				if method == m {
					methodAllowed = true
					break
				}
			}
			if !methodAllowed {
				c.AbortWithStatus(http.StatusMethodNotAllowed)
				return
			}
		}

	})

	//Check DB health
	r.GET("/healthz", health.HealthCheckHandler(db))

	//Create User
	r.POST("/v1/user", user.CreateUserHandler(db))

	//Verify User
	r.GET("/verify", user.VerifyUserHandler(db))

	authGroup := r.Group("/")
	authGroup.Use(AuthenticationMiddleware(db))
	{
		// Routes that also require email verification
		verifiedGroup := authGroup.Group("/")
		verifiedGroup.Use(EmailVerificationMiddleware(db))
		{
			verifiedGroup.GET("/v1/user/self", user.GetUserDetails(db))
			verifiedGroup.PUT("/v1/user/self", user.UpdateUserHandler(db))
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{})
	})

	return r
}

func EmailVerificationMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			logger.Logger.Error("EmailVerificationMiddleware() - User ID not found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var user user.UserModel
		if err := db.First(&user, "id = ?", userID).Error; err != nil {
			logger.Logger.Error("EmailVerificationMiddleware() - Failed to retrieve user")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		if !user.IsVerified {
			logger.Logger.Error("EmailVerificationMiddleware() - User's email address is not verified")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Email address is not verified"})
			return
		}

		c.Next()
	}
}
