package controllers

import (
	"bankapp/models"
	"bankapp/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetBalance(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	var accounts []models.Account
	db := c.MustGet("db").(*gorm.DB)

	// Fetch accounts belonging to the user
	if err := db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accounts"})
		return
	}

	totalBalance := 0.0
	for _, account := range accounts {
		totalBalance += account.Balance
	}

	c.JSON(http.StatusOK, gin.H{
		"total_balance": totalBalance,
		"accounts":      accounts,
	})
}

func Transfer(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	var input struct {
		FromAccountID uint    `json:"from_account_id" binding:"required"`
		ToAccountID   uint    `json:"to_account_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required,gt=0"`
		Category      string  `json:"category"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var fromAccount, toAccount models.Account

	// Verify ownership of the fromAccount
	if err := db.Where("id = ? AND user_id = ?", input.FromAccountID, userID).First(&fromAccount).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Source account not found or does not belong to the user"})
		return
	}

	// Retrieve the destination account
	if err := db.First(&toAccount, input.ToAccountID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination account not found"})
		return
	}

	if fromAccount.Balance < input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	// Database transaction for atomicity
	err := db.Transaction(func(tx *gorm.DB) error {
		fromAccount.Balance -= input.Amount
		if err := tx.Save(&fromAccount).Error; err != nil {
			return err
		}

		toAccount.Balance += input.Amount
		if err := tx.Save(&toAccount).Error; err != nil {
			return err
		}

		transaction := models.Transaction{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Amount:        input.Amount,
			Category:      input.Category,
			Timestamp:     tx.NowFunc(),
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}

func CreateAccount(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	var input struct {
		InitialBalance float64 `json:"initial_balance" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Ошибка привязки JSON:", err) // Добавьте эту строку для отладки
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account := models.Account{
		UserID:  userID,
		Balance: input.InitialBalance,
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account created successfully",
		"account": account,
	})
}

// GetAccounts retrieves all accounts belonging to the authenticated user
func GetAccounts(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return
	}

	var accounts []models.Account
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}
