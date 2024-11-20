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

func setupBudgetRouter() (*gin.Engine, *gorm.DB, uint) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Migrate models
	db.AutoMigrate(&models.User{}, &models.Budget{})

	// Create a test user
	user := models.User{Username: "testuser", Password: "hashedpassword"}
	db.Create(&user)

	// Middleware to inject db and userID
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("userID", user.ID)
		c.Next()
	})

	// Routes
	router.POST("/budget", controllers.SetBudget)
	router.GET("/budget", controllers.GetBudget)

	return router, db, user.ID
}

func TestSetBudget(t *testing.T) {
	router, db, userID := setupBudgetRouter()

	// Test data
	budgetData := map[string]interface{}{
		"category": "Utilities",
		"amount":   500.0,
	}
	jsonData, _ := json.Marshal(budgetData)

	req, _ := http.NewRequest("POST", "/budget", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Budget set successfully")

	// Verify budget in the database
	var budgets []models.Budget
	db.Where("user_id = ?", userID).Find(&budgets)
	assert.Equal(t, 1, len(budgets))
	assert.Equal(t, "Utilities", budgets[0].Category)
	assert.Equal(t, 500.0, budgets[0].Amount)
}
