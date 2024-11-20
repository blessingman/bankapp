// routes/routes.go
package routes

import (
	"bankapp/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Authentication routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(controllers.AuthMiddleware)
	{
		authorized.GET("/balance", controllers.GetBalance)
		authorized.POST("/transfer", controllers.Transfer)
		authorized.GET("/transactions", controllers.GetTransactions)
		authorized.POST("/budget", controllers.SetBudget)
		authorized.GET("/budget", controllers.GetBudget)

		// Add the /accounts route
		authorized.POST("/accounts", controllers.CreateAccount)
		authorized.GET("/accounts", controllers.GetAccounts)
	}
}
