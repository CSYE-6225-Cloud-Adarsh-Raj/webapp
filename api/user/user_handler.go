package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// type User struct {
// 	gorm.Model
// 	FirstName      string    `json:"first_name" validate:"required"`
// 	LastName       string    `json:"last_name" validate:"required"`
// 	Password       string    `json:"password" validate:"required" writeOnly:"true"`
// 	Username       string    `json:"username" validate:"required,email"`
// }

type UserModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"` //Do we need to change it to string
	CreatedAt time.Time `json:"account_created" readOnly:"true"`
	UpdatedAt time.Time `json:"account_updated" readOnly:"true"`
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	Password  string    `json:"password" validate:"required" writeOnly:"true"`
	Username  string    `json:"username" validate:"required,email"`
}

// HashPassword hashes the user's password.
func (u *UserModel) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// // generateUUID generates a new UUID string.
// func generateUUID() string {
// 	return uuid.New().String()
// }

func CreateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user UserModel
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error in json body": err.Error()})
			return
		}

		user.ID = uuid.New()
		// user.ID = user.generateUUID()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		// Hash the password using bcrypt
		if err := user.HashPassword(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		// Insert user into database
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user to database"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":              user.ID,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"username":        user.Username,
			"account_created": user.CreatedAt.Format(time.RFC3339Nano),
			"account_updated": user.UpdatedAt.Format(time.RFC3339Nano),
		})
	}
}
