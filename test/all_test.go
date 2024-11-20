package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bankapp/models"
	"bankapp/routes"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Set up the router and mock environment
func setupRouter() (*gin.Engine, *gorm.DB, uint) {
	os.Setenv("JWT_SECRET", "secret")

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	db.AutoMigrate(&models.User{}, &models.Account{}, &models.Transaction{}, &models.Budget{})

	// Hash the password before storing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{Username: "testuser", Password: string(hashedPassword)}
	db.Create(&user)

	// Mock Accounts
	db.Create(&models.Account{UserID: user.ID, Balance: 2000})
	db.Create(&models.Account{UserID: user.ID, Balance: 0})

	// Mock Budget
	db.Create(&models.Budget{UserID: user.ID, Category: "Utilities", Amount: 1500})

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("userID", user.ID)
		c.Next()
	})

	routes.SetupRoutes(router)

	return router, db, user.ID
}

// Generate a mock JWT token
func generateMockToken() string {
	var signingKey = []byte("secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,          // Mock user ID
		"exp":     9999999999, // Far future expiration for testing
	})

	tokenString, _ := token.SignedString(signingKey)

	return "Bearer " + tokenString
}

// Tests for accounts
func TestAccounts(t *testing.T) {
	router, _, _ := setupRouter()

	t.Run("TestCreateAccount", func(t *testing.T) {
		// Prepare test data
		accountData := map[string]float64{
			"initial_balance": 1000.0,
		}
		jsonData, _ := json.Marshal(accountData)

		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", generateMockToken()) // Add Authorization header

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Account created successfully")
	})

	t.Run("TestGetAccounts", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/accounts", nil)
		req.Header.Set("Authorization", generateMockToken())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"accounts"`)
	})
}

// Tests for authentication
func TestAuth(t *testing.T) {
	router, _, _ := setupRouter()

	t.Run("TestRegister", func(t *testing.T) {
		reqBody := []byte(`{"username":"newuser", "password":"password123"}`)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Registration successful")
	})

	t.Run("TestLogin", func(t *testing.T) {
		reqBody := []byte(`{"username":"testuser", "password":"password123"}`)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})
}
