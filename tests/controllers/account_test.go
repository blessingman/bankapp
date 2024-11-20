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
	"gorm.io/gorm"
)

func setupAccountRouter() (*gin.Engine, *gorm.DB, uint) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Use the pure Go SQLite driver
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Migrate the models
	db.AutoMigrate(&models.User{}, &models.Account{})

	// Create a test user
	user := models.User{Username: "testuser", Password: "hashedpassword"}
	db.Create(&user)

	// Middleware to inject db and userID into context
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("userID", user.ID)
		c.Next()
	})

	// Routes
	router.POST("/accounts", controllers.CreateAccount)
	router.GET("/accounts", controllers.GetAccounts)
	router.GET("/balance", controllers.GetBalance)

	return router, db, user.ID
}

func TestCreateAccount(t *testing.T) {
	router, _, _ := setupAccountRouter()

	// Prepare test data
	accountData := map[string]float64{
		"initial_balance": 1000.0,
	}
	jsonData, _ := json.Marshal(accountData)

	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Account created successfully")
}
