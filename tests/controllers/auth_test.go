package controllers_test

import (
	"bankapp/controllers"
	"bankapp/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupAuthRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Migrate user model
	db.AutoMigrate(&models.User{})

	// Middleware to inject the database
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Setup routes
	router.POST("/auth/register", controllers.Register)
	router.POST("/auth/login", controllers.Login)

	return router, db
}

func TestRegister(t *testing.T) {
	router, db := setupAuthRouter()

	// Test data
	user := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	jsonData, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Registration successful")

	// Check database
	var createdUser models.User
	result := db.First(&createdUser, "username = ?", "testuser")
	assert.Nil(t, result.Error)
	assert.Equal(t, "testuser", createdUser.Username)
}

func TestLogin(t *testing.T) {
	router, db := setupAuthRouter()

	// Create a test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	db.Create(&models.User{Username: "testuser", Password: string(hashedPassword)})

	// Test data
	credentials := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	jsonData, _ := json.Marshal(credentials)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Login successful")

	// Check token in response
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["token"])
}
