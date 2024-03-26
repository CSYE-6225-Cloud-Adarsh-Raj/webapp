package user

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"webapp/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"cloud.google.com/go/pubsub"
)

var projectID = "csye6225-dev-414220"
var topicName = "verify_email"

type UserModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"-"`
	CreatedAt  time.Time `json:"-" readOnly:"true"`
	UpdatedAt  time.Time `json:"-" readOnly:"true"`
	FirstName  string    `json:"first_name" validate:"required"`
	LastName   string    `json:"last_name" validate:"required"`
	Password   string    `json:"password" validate:"required" writeOnly:"true"`
	Username   string    `json:"username" validate:"required,email"`
	IsVerified bool      `json:"is_verified" gorm:"default:false"`
}

type VerificationMessage struct {
	Email             string    `json:"email"`
	VerificationToken uuid.UUID `json:"verificationToken"`
}

type EmailVerification struct {
	Email    string    `gorm:"primaryKey;type:varchar(100)"`
	UUID     uuid.UUID `gorm:"type:uuid;unique"`
	TimeSent time.Time
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
			// fmt.Println("CreateUserHandler() - Error in json body")
			logger.Logger.Error("CreateUserHandler() - Error in json body")
			c.Status(http.StatusBadRequest)
			return
		}

		var validate = validator.New()
		if validationErr := validate.Struct(user); validationErr != nil {
			// fmt.Println("CreateUserHandler() - Validation Error:", validationErr.Error())
			logger.Logger.Error("CreateUserHandler() - Validation Error")
			c.Status(http.StatusBadRequest)
			return
		}

		user.ID = uuid.New()
		user.IsVerified = false

		var verMsg VerificationMessage
		// user.VerificationToken = user.ID
		verMsg.VerificationToken = user.ID

		if strings.Contains(user.Username, ":") {
			// fmt.Println("CreateUserHandler() -Error: Username cannot contain ':' ")
			logger.Logger.Error("CreateUserHandler() - Username cannot contain ':'")
			c.Status(http.StatusBadRequest)
			return
		}

		// Check for existing email
		var count int64
		db.Model(&UserModel{}).Where("username = ?", user.Username).Count(&count)
		if count > 0 {
			// fmt.Println("CreateUserHandler() - Error: Email-id already exists")
			logger.Logger.Error("CreateUserHandler() - Email-id already exists")
			c.Status(http.StatusBadRequest)
			return
		}

		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		// Hash the password using bcrypt
		if err := user.HashPassword(); err != nil {
			// fmt.Println("CreateUserHandler() - Error hashing password")
			logger.Logger.Error("CreateUserHandler() - Error hashing password")
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := db.Create(&user).Error; err != nil {
			// fmt.Println("CreateUserHandler() - Error saving user to database")
			logger.Logger.Error("CreateUserHandler() - Error saving user to database")
			c.Status(http.StatusInternalServerError)
			return
		}

		createdAtFormatted := user.CreatedAt.UTC().Format(time.RFC3339Nano)
		createdAtFormatted = strings.Replace(createdAtFormatted, "+00:00", "Z", 1)

		updatedAtformatted := user.UpdatedAt.UTC().Format(time.RFC3339Nano)
		updatedAtformatted = strings.Replace(updatedAtformatted, "+00:00", "Z", 1)

		logger.Logger.WithFields(logrus.Fields{
			"id":              user.ID,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"username":        user.Username,
			"account_created": createdAtFormatted,
			"account_updated": updatedAtformatted,
		}).Info("User created successfully")

		// Initialize Pub/Sub client
		ctx := context.Background()
		if os.Getenv("SKIP_PUBSUB") != "true" {
			pubsubClient, err := pubsub.NewClient(ctx, projectID)
			if err != nil {
				logger.Logger.Errorf("Failed to create Pub/Sub client: %v", err)
				c.Status(http.StatusInternalServerError)
				return
			}
			defer pubsubClient.Close()

			// Get the Pub/Sub topic
			topic := pubsubClient.Topic(topicName)

			// Create an instance of VerificationMessage with the necessary data
			vMessage := VerificationMessage{
				Email:             user.Username,
				VerificationToken: verMsg.VerificationToken,
			}

			// Marshal the VerificationMessage into JSON
			jsonData, err := json.Marshal(vMessage)
			if err != nil {
				logger.Logger.Errorf("Error marshaling verification message: %v", err)
				c.Status(http.StatusInternalServerError)
				return
			}

			// Prepare and publish the message with jsonData
			msg := &pubsub.Message{
				Data: jsonData,
				Attributes: map[string]string{
					"email": user.Username,
				},
			}

			result := topic.Publish(ctx, msg)

			// Wait for the result
			id, err := result.Get(ctx)
			if err != nil {
				logger.Logger.Errorf("Failed to publish to Pub/Sub: %v", err)
				c.Status(http.StatusInternalServerError)
				return
			}
			logger.Logger.Infof("Published message with ID: %s", id)

		}

		c.JSON(http.StatusCreated, gin.H{
			"id":              user.ID,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"username":        user.Username,
			"account_created": createdAtFormatted,
			"account_updated": updatedAtformatted,
		})

		logger.Logger.Debug("Completed Execution of CreateUserHandler")

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
			// fmt.Println("GetUserDetails() -Error:: User not authenticated")
			logger.Logger.Error("GetUserDetails() - User not authenticated")
			c.Status(http.StatusUnauthorized)
			return
		}

		var user UserModel
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			// fmt.Println("GetUserDetails() - Error: Failed to retrieve user details")
			logger.Logger.Error("GetUserDetails() - Failed to retrieve user details")
			c.Status(http.StatusInternalServerError)
			return
		}

		createdAtFormatted := user.CreatedAt.UTC().Format(time.RFC3339Nano)
		createdAtFormatted = strings.Replace(createdAtFormatted, "+00:00", "Z", 1)

		updatedAtformatted := user.UpdatedAt.UTC().Format(time.RFC3339Nano)
		updatedAtformatted = strings.Replace(updatedAtformatted, "+00:00", "Z", 1)

		logger.Logger.WithFields(logrus.Fields{
			"id":              user.ID,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"username":        user.Username,
			"account_created": createdAtFormatted,
			"account_updated": updatedAtformatted,
		}).Info("User details retrieved successfully")

		c.JSON(http.StatusOK, gin.H{
			"id":              user.ID,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"username":        user.Username,
			"account_created": createdAtFormatted,
			"account_updated": updatedAtformatted,
		})
		logger.Logger.Debug("Completed Execution of GetUserDetails")
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
		// fmt.Println("GetUserID() - Error: Not able to get user ID")
		logger.Logger.Error("GetUserID() - Not able to get user ID")
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
			// fmt.Println("UpdateUserHandler() - Error: Invalid request")
			logger.Logger.Error("UpdateUserHandler() - Invalid request")
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
					// fmt.Printf("UpdateUserHandler() - Error: Field '%s' cannot be empty\n", key)
					logger.Logger.Errorf("UpdateUserHandler() - Field '%s' cannot be empty", key)
					c.Status(http.StatusBadRequest)
					return
				}
			} else {
				// fmt.Printf("UpdateUserHandler() - Error: Field '%s' not allowed\n", key)
				logger.Logger.Errorf("UpdateUserHandler() - Field '%s' not allowed", key)
				c.Status(http.StatusBadRequest)
				return
			}
		}

		if providedFieldsCount < 3 {
			// fmt.Println("UpdateUserHandler() - Error: At least three fields must be provided for update with non-empty values")
			logger.Logger.Error("UpdateUserHandler() -  At least three fields must be provided for update with non-empty values")
			c.Status(http.StatusBadRequest)
			return
		}

		userID := GetUserID(c)

		firstName, _ := userDetails["first_name"].(string)
		lastName, _ := userDetails["last_name"].(string)
		password, _ := userDetails["password"].(string)

		if err := updateUserDetails(db, userID, firstName, lastName, password); err != nil {
			// fmt.Println("UpdateUserHandler() - Error: Failed to update user details")
			logger.Logger.Error("UpdateUserHandler() - Failed to update user details")
			c.Status(http.StatusInternalServerError)
			return
		}

		// fmt.Println("UpdateUserHandler() - User details updated successfully")
		logger.Logger.Info("UpdateUserHandler() - User details updated successfully")
		c.Status(http.StatusNoContent)
		logger.Logger.Debug("Completed Execution of UpdateUserHandler")
	}
}

func VerifyUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenParam := c.Query("token")
		if tokenParam == "" {
			logger.Logger.Error("VerifyUserHandler() - Missing token")
			c.Status(http.StatusBadRequest)
			return
		}

		token, err := uuid.Parse(tokenParam)
		if err != nil {
			logger.Logger.Error("VerifyUserHandler() - Invalid token format")
			c.Status(http.StatusBadRequest)
			return
		}

		var user UserModel
		if err := db.Where("id = ?", token).First(&user).Error; err != nil {
			logger.Logger.Error("VerifyUserHandler() - ID not found")
			c.Status(http.StatusNotFound)
			return
		}

		// // Check if the request is within 2 minutes of the user's creation time
		// if time.Since(user.CreatedAt) > 2*time.Minute {
		// 	logger.Logger.Error("VerifyUserHandler() - Verification link expired")
		// 	c.Status(http.StatusBadRequest)
		// 	return
		// }

		// Fetch the EmailVerification record associated with the user
		var emailVerification EmailVerification
		if err := db.Where("uuid = ?", user.ID).First(&emailVerification).Error; err != nil {
			logger.Logger.Error("VerifyUserHandler() - Failed to retrieve email verification details")
			c.Status(http.StatusUnauthorized)
			return
		}

		// Check if the request is within 2 minutes of the email's TimeSent
		if time.Since(emailVerification.TimeSent) > 2*time.Minute {
			logger.Logger.Error("VerifyUserHandler() - Verification link expired")
			c.Status(http.StatusBadRequest)
			return
		}

		// Update the isVerified flag
		user.IsVerified = true
		if err := db.Save(&user).Error; err != nil {
			logger.Logger.Error("VerifyUserHandler() - Failed to verify user")
			c.Status(http.StatusUnauthorized)
			return
		}

		logger.Logger.Info("User verified successfully")
		c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})

		logger.Logger.Debug("Completed Execution of VerifyUserHandler")
	}
}
