package user

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"-"`
	CreatedAt time.Time `json:"-" readOnly:"true"`
	UpdatedAt time.Time `json:"-" readOnly:"true"`
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

func CreateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user UserModel
		if err := c.ShouldBindJSON(&user); err != nil {
			fmt.Println("CreateUserHandler() - Error in json body")
			c.Status(http.StatusBadRequest)
			return
		}

		var validate = validator.New()
		if validationErr := validate.Struct(user); validationErr != nil {
			fmt.Println("CreateUserHandler() - Validation Error:", validationErr.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		user.ID = uuid.New()

		if strings.Contains(user.Username, ":") {
			fmt.Println("CreateUserHandler() -Error: Username cannot contain ':' ")
			c.Status(http.StatusBadRequest)
			return
		}

		// Check for existing email
		var count int64
		db.Model(&UserModel{}).Where("username = ?", user.Username).Count(&count)
		if count > 0 {
			fmt.Println("CreateUserHandler() - Error: Email-id already exists")
			c.Status(http.StatusBadRequest)
			return
		}

		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		// Hash the password using bcrypt
		if err := user.HashPassword(); err != nil {
			fmt.Println("CreateUserHandler() - Error hashing password")
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := db.Create(&user).Error; err != nil {
			fmt.Println("CreateUserHandler() - Error saving user to database")
			c.Status(http.StatusInternalServerError)
			return
		}

		createdAtFormatted := user.CreatedAt.UTC().Format(time.RFC3339Nano)
		createdAtFormatted = strings.Replace(createdAtFormatted, "+00:00", "Z", 1)

		updatedAtformatted := user.UpdatedAt.UTC().Format(time.RFC3339Nano)
		updatedAtformatted = strings.Replace(updatedAtformatted, "+00:00", "Z", 1)

		c.JSON(http.StatusCreated, gin.H{
			"id":              user.ID,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"username":        user.Username,
			"account_created": createdAtFormatted,
			"account_updated": updatedAtformatted,
		})
	}
}

func (u *UserModel) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func GetUserDetails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > 0 || len(c.Request.URL.Query()) > 0 {
			c.Status(http.StatusBadRequest)
			return
		}

		username, exists := c.Get("username")
		if !exists {
			fmt.Println("GetUserDetails() -Error:: User not authenticated")
			c.Status(http.StatusUnauthorized)
			return
		}

		var user UserModel
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			fmt.Println("GetUserDetails() - Error: Failed to retrieve user details")
			c.Status(http.StatusInternalServerError)
			return
		}

		createdAtFormatted := user.CreatedAt.UTC().Format(time.RFC3339Nano)
		createdAtFormatted = strings.Replace(createdAtFormatted, "+00:00", "Z", 1)

		updatedAtformatted := user.UpdatedAt.UTC().Format(time.RFC3339Nano)
		updatedAtformatted = strings.Replace(updatedAtformatted, "+00:00", "Z", 1)

		c.JSON(http.StatusOK, gin.H{
			"id":              user.ID,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"username":        user.Username,
			"account_created": createdAtFormatted,
			"account_updated": updatedAtformatted,
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
		fmt.Println("GetUserID() - Error: Not able to get user ID")
		return uuid.Nil
	}
	return userID.(uuid.UUID)
}

func updateUserDetails(db *gorm.DB, userID uuid.UUID, firstName, lastName, password string) error {
	var user UserModel
	updateFields := make(map[string]interface{})
	if firstName != "" {
		updateFields["FirstName"] = firstName
	}
	if lastName != "" {
		updateFields["LastName"] = lastName
	}
	if password != "" {
		user.Password = password
		if err := user.HashPassword(); err != nil {
			return err
		}
		updateFields["Password"] = user.Password
	}

	result := db.Model(&UserModel{}).Where("id = ?", userID).Updates(updateFields)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userDetails := make(map[string]interface{})
		if err := c.BindJSON(&userDetails); err != nil {
			fmt.Println("UpdateUserHandler() - Error: Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		allowedFields := map[string]bool{
			"first_name": true,
			"last_name":  true,
			"password":   true,
		}

		providedFieldsCount := 0
		for key, value := range userDetails {
			if allowedFields[key] {
				if valueStr, ok := value.(string); ok && valueStr != "" {
					providedFieldsCount++
				} else {
					fmt.Printf("UpdateUserHandler() - Error: Field '%s' cannot be empty\n", key)
					c.Status(http.StatusBadRequest)
					return
				}
			} else {
				fmt.Printf("UpdateUserHandler() - Error: Field '%s' not allowed\n", key)
				c.Status(http.StatusBadRequest)
				return
			}
		}

		if providedFieldsCount < 3 {
			fmt.Println("UpdateUserHandler() - Error: At least three fields must be provided for update with non-empty values")
			c.Status(http.StatusBadRequest)
			return
		}

		userID := GetUserID(c)

		firstName, _ := userDetails["first_name"].(string)
		lastName, _ := userDetails["last_name"].(string)
		password, _ := userDetails["password"].(string)

		if err := updateUserDetails(db, userID, firstName, lastName, password); err != nil {
			fmt.Println("UpdateUserHandler() - Error: Failed to update user details")
			c.Status(http.StatusInternalServerError)
			return
		}

		fmt.Println("UpdateUserHandler() - User details updated successfully")
		c.Status(http.StatusNoContent)
	}
}
