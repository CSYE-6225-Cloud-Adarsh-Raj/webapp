package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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

		// Check for existing email
		var count int64
		db.Model(&UserModel{}).Where("username = ?", user.Username).Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
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

func (u *UserModel) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func GetUserDetails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var user UserModel
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user details"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":              user.ID,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"username":        user.Username,
			"account_created": user.CreatedAt.Format(time.RFC3339Nano),
			"account_updated": user.UpdatedAt.Format(time.RFC3339Nano),
		})
	}
}

func ValidateCredentials(db *gorm.DB, username, password string) bool {
	var user UserModel
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}

	if err := user.CheckPassword(password); err != nil {
		return false
	}

	return true
}

func GetUserID(c *gin.Context) uuid.UUID {
	userID, exists := c.Get("userID")
	if !exists {
		fmt.Println("GetUserID() - Not able to get user ID")
		return uuid.Nil
	}
	return userID.(uuid.UUID)
}

func UpdateUserDetails(db *gorm.DB, userID uuid.UUID, firstName, lastName, password string) error {
	var user UserModel
	user.Password = password
	if err := user.HashPassword(); err != nil {
		fmt.Println("UpdateUserDetails() - Could not hash password")
		return err
	}

	hashedPassword := user.Password
	result := db.Model(&UserModel{}).Where("id = ?", userID).Updates(UserModel{FirstName: firstName, LastName: lastName, Password: hashedPassword})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userDetails struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Password  string `json:"password"`
		}
		if err := c.BindJSON(&userDetails); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		userID := GetUserID(c)

		// Update user details
		if err := UpdateUserDetails(db, userID, userDetails.FirstName, userDetails.LastName, userDetails.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user details"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User details updated successfully"})
	}
}
