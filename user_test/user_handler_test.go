package usertest

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"webapp/api/user"
	"webapp/logger"
	"webapp/router"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func setupTestDatabase() *gorm.DB {
	host := os.Getenv("TEST_DB_HOST")
	user := os.Getenv("TEST_DB_USER")
	password := os.Getenv("TEST_DB_PASSWORD")
	dbName := os.Getenv("TEST_DB_NAME")
	port := "5432"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func teardownDatabase() {
	db.Migrator().DropTable(&user.UserModel{})
}

func TestMain(m *testing.M) {
	db = setupTestDatabase()
	err := db.AutoMigrate(&user.UserModel{})
	if err != nil {
		fmt.Println("Failed to migrate testtable schema")
		logger.Logger.Error("TestMain() - Failed to migrate testtable schema")
		os.Exit(1)
	}
	code := m.Run()
	teardownDatabase()
	os.Exit(code)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func unmarshalResponseBody(t *testing.T, body *bytes.Buffer) map[string]interface{} {
	var responseBody map[string]interface{}
	err := json.Unmarshal(body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	return responseBody
}

func TestCreateAndGetUser(t *testing.T) {
	r := router.InitRouter(db)

	defer func() {
		db.Where("username = ?", "john.doe@example.com").Delete(&user.UserModel{})
	}()

	// Test 1: Create a new user with valid data
	userData := map[string]string{
		"first_name": "John",
		"last_name":  "Doe",
		"username":   "john.doe@example.com",
		"password":   "password123",
	}
	userDataBytes, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(userDataBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d for valid user creation", http.StatusCreated, w.Code)
	}

	// Test 2: Get user details with valid authentication
	req, _ = http.NewRequest("GET", "/v1/user/self", nil)
	req.Header.Set("Authorization", "Basic "+basicAuth("john.doe@example.com", "password123"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d for valid user retrieval", http.StatusOK, w.Code)
	}

	responseBody := unmarshalResponseBody(t, w.Body)
	if responseBody["username"] != userData["username"] {
		t.Fatalf("Expected username %s, got %s", userData["username"], responseBody["username"])
	}

	// fmt.Println("ALL TESTS PASSED in TestCreateAndGetUser() !!!")
	logger.Logger.Error("TestCreateAndGetUser() - ALL TESTS PASSED !!!")
}

func TestUpdateAndGetUser(t *testing.T) {
	r := router.InitRouter(db)
	// Step 1: Create a user directly in the database for testing
	testUser := user.UserModel{
		FirstName: "Johny",
		LastName:  "Doe",
		Username:  "john.update@example.com",
		Password:  "password@123",
	}

	if err := testUser.HashPassword(); err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	defer func() {
		db.Where("username = ?", "john.update@example.com").Delete(&user.UserModel{})
	}()

	// Step 2: Update the user's details through an HTTP request
	updatedUserData := map[string]string{
		"first_name": "Jane",
		"last_name":  "Doe",
		"password":   "newpassword123",
	}
	updatedUserDataBytes, _ := json.Marshal(updatedUserData)
	req, _ := http.NewRequest("PUT", "/v1/user/self", bytes.NewBuffer(updatedUserDataBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth("john.update@example.com", "password@123"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("Expected status code %d, got %d for valid user update", http.StatusNoContent, w.Code)
	}

	// Step 3: Validate the account was updated using a GET call
	req, _ = http.NewRequest("GET", "/v1/user/self", nil)
	req.Header.Set("Authorization", "Basic "+basicAuth("john.update@example.com", "newpassword123"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d for valid user retrieval after update", http.StatusOK, w.Code)
	}

	responseBody := unmarshalResponseBody(t, w.Body)
	if responseBody["first_name"] != updatedUserData["first_name"] {
		t.Fatalf("Expected first name %s, got %s", updatedUserData["first_name"], responseBody["first_name"])
	}

	// fmt.Println("ALL TESTS PASSED in TestUpdateAndGetUser() !!!")
	logger.Logger.Error("TestUpdateAndGetUser() - ALL TESTS PASSED!!!")
}
