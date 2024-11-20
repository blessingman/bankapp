package controllers

import (
	"bankapp/models"
	"bankapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetBudget(c *gin.Context) {
	// Retrieve user ID from context
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var input struct {
		Category string  `json:"category" binding:"required"`
		Amount   float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var budget models.Budget

	// Check if budget exists
	if err := db.Where("user_id = ? AND category = ?", userID, input.Category).First(&budget).Error; err == nil {
		// Update existing budget
		budget.Amount = input.Amount
		if err := db.Save(&budget).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
			return
		}
	} else if err == gorm.ErrRecordNotFound {
		// Create new budget
		budget = models.Budget{
			UserID:   userID,
			Category: input.Category,
			Amount:   input.Amount,
		}
		if err := db.Create(&budget).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create budget"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget set successfully"})
}

func GetBudget(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	var budgets []models.Budget
	db := c.MustGet("db").(*gorm.DB)

	// Fetch budgets belonging to the user
	if err := db.Where("user_id = ?", userID).Find(&budgets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve budgets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"budgets": budgets})
}
