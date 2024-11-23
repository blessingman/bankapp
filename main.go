package main

import (
	"bankapp/models"
	"bankapp/routes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors" // Импортируем пакет для CORS
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Загрузка переменных окружения из файла .env, если он есть
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Получение переменных окружения
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Строка подключения к базе данных
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Миграция моделей базы данных
	db.AutoMigrate(&models.User{}, &models.Account{}, &models.Transaction{}, &models.Budget{})

	// Инициализация роутера Gin
	r := gin.Default()

	// Настройка CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:80"}, // Укажите URL вашего фронтенда
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware для предоставления подключения к базе данных контроллерам
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Настройка маршрутов
	routes.SetupRoutes(r)

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
