package controllers_test

import (
	"bankapp/controllers"
	"bankapp/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTransactionRouter() (*gin.Engine, *gorm.DB, uint) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Migrate models
	db.AutoMigrate(&models.User{}, &models.Account{}, &models.Transaction{})

	// Create a user and accounts
	user := models.User{Username: "testuser", Password: "hashedpassword"}
	db.Create(&user)
	fromAccount := models.Account{UserID: user.ID, Balance: 1000.0}
	toAccount := models.Account{UserID: user.ID, Balance: 500.0}
	db.Create(&fromAccount)
	db.Create(&toAccount)

	// Middleware to inject db and userID
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("userID", user.ID)
		c.Next()
	})

	// Routes
	router.POST("/transfer", controllers.Transfer)
	router.GET("/transactions", controllers.GetTransactions)

	return router, db, user.ID
}

func TestTransfer(t *testing.T) {
	router, db, _ := setupTransactionRouter()

	// Get accounts
	var accounts []models.Account
	db.Find(&accounts)

	// Test data
	transferData := map[string]interface{}{
		"from_account_id": accounts[0].ID,
		"to_account_id":   accounts[1].ID,
		"amount":          200.0,
		"category":        "Utilities",
	}
	jsonData, _ := json.Marshal(transferData)

	req, _ := http.NewRequest("POST", "/transfer", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Transfer successful")

	// Verify balances
	db.First(&accounts[0])
	db.First(&accounts[1])
	assert.Equal(t, 800.0, accounts[0].Balance)
	assert.Equal(t, 700.0, accounts[1].Balance)
}

func TestGetTransactions(t *testing.T) {
	router, db, _ := setupTransactionRouter()

	// Create transactions
	db.Create(&models.Transaction{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        200.0,
		Category:      "Utilities",
		Timestamp:     time.Now(),
	})

	req, _ := http.NewRequest("GET", "/transactions", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"amount":200.0`)
}
