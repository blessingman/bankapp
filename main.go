package main

import (
	"bankapp/models" // Пакет с моделями базы данных
	"bankapp/routes" // Пакет с маршрутизацией
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors" // Пакет для настройки CORS
	"github.com/gin-gonic/gin"    // Основной веб-фреймворк
	"github.com/joho/godotenv"    // Для загрузки переменных окружения из .env файла
	"gorm.io/driver/postgres"     // Драйвер PostgreSQL для GORM
	"gorm.io/gorm"                // ORM-библиотека GORM
)

func main() {
	// Загрузка переменных окружения из файла .env, если он существует
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found") // Логируем, если файл не найден
	}

	// Получение переменных окружения (данные для подключения к базе данных)
	dbHost := os.Getenv("DB_HOST")         // Хост базы данных
	dbUser := os.Getenv("DB_USER")         // Пользователь базы данных
	dbPassword := os.Getenv("DB_PASSWORD") // Пароль
	dbName := os.Getenv("DB_NAME")         // Имя базы данных
	dbPort := os.Getenv("DB_PORT")         // Порт базы данных

	// Формирование строки подключения к базе данных
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err) // Завершаем выполнение, если подключение не удалось
	}

	// Автоматическая миграция схемы базы данных для всех моделей
	db.AutoMigrate(&models.User{}, &models.Account{}, &models.Transaction{}, &models.Budget{})

	// Инициализация веб-сервера Gin
	r := gin.Default()

	// Настройка CORS middleware для разрешения запросов с определенных источников
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},                        // Укажите адрес фронтенда
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Разрешенные HTTP-методы
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Разрешенные заголовки
		ExposeHeaders:    []string{"Content-Length"},                          // Заголовки для экспорта
		AllowCredentials: true,                                                // Разрешение на передачу cookies
		MaxAge:           12 * time.Hour,                                      // Максимальное время кеширования
	}))

	// Middleware для передачи подключения к базе данных в контекст запроса
	r.Use(func(c *gin.Context) {
		c.Set("db", db) // Добавляем объект базы данных в контекст
		c.Next()        // Передаем управление следующему middleware или маршруту
	})

	// Инициализация маршрутов (API)
	routes.SetupRoutes(r)

	// Запуск веб-сервера на порту 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err) // Логируем ошибку при запуске сервера
	}
}
