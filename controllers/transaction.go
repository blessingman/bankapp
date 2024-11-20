package controllers

import (
	"bankapp/models"
	"bankapp/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTransactions(c *gin.Context) {
	// Retrieve user ID from context
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// Get user's account IDs
	var accounts []models.Account
	if err := db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accounts"})
		return
	}

	if len(accounts) == 0 {
		c.JSON(http.StatusOK, gin.H{"transactions": []models.Transaction{}})
		return
	}

	accountIDs := make([]uint, len(accounts))
	for i, account := range accounts {
		accountIDs[i] = account.ID
	}

	// Build query
	query := db.Where("from_account_id IN ? OR to_account_id IN ?", accountIDs, accountIDs)

	// Apply filters
	category := c.Query("category")
	if category != "" {
		query = query.Where("category = ?", category)
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	if startDateStr != "" && endDateStr != "" {
		layout := "2006-01-02"
		startDate, err1 := time.Parse(layout, startDateStr)
		endDate, err2 := time.Parse(layout, endDateStr)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		query = query.Where("timestamp BETWEEN ? AND ?", startDate, endDate)
	}

	// Execute query
	var transactions []models.Transaction
	if err := query.Order("timestamp desc").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
