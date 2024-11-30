package controllers

import (
	"bankapp/models" // Импорт моделей проекта
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"  // Gin-фреймворк для обработки запросов
	"github.com/golang-jwt/jwt" // JWT-библиотека для работы с токенами

	"gorm.io/gorm" // GORM ORM для работы с базой данных
)

// AuthMiddleware — Middleware для проверки авторизации пользователя с использованием JWT токена
func AuthMiddleware(c *gin.Context) {
	// Извлекаем заголовок Authorization из HTTP-запроса
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" { // Проверяем, что заголовок существует
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Удаляем префикс "Bearer" и лишние пробелы из заголовка
	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	if tokenString == "" { // Проверяем, что токен не пустой
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
		return
	}

	// Парсим и проверяем токен
	token, err := parseJWTToken(tokenString)
	if err != nil || !token.Valid { // Если токен недействителен, возвращаем ошибку
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Извлекаем user_id из Claims (данные, закодированные в токене)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok { // Проверяем, что Claims являются корректной картой
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Проверяем, что user_id доступен и является числом
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}
	userID := uint(userIDFloat) // Приводим user_id к типу uint

	// Проверяем, существует ли пользователь с данным user_id в базе данных
	db := c.MustGet("db").(*gorm.DB) // Извлекаем подключение к базе данных из контекста Gin
	var user models.User
	if err := db.First(&user, userID).Error; err != nil { // Проверяем, что пользователь существует
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Устанавливаем userID в контекст запроса, чтобы другие контроллеры могли его использовать
	c.Set("userID", user.ID)
	c.Next() // Продолжаем выполнение следующего обработчика
}

// parseJWTToken — функция для парсинга и проверки JWT токена
func parseJWTToken(tokenString string) (*jwt.Token, error) {
	// Извлекаем секретный ключ для подписи токена из переменных окружения
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" { // Если ключ отсутствует, возвращаем ошибку
		return nil, fmt.Errorf("JWT secret key not set")
	}

	// Парсим токен с помощью секретного ключа
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Возвращаем секретный ключ для проверки подписи токена
		return []byte(secretKey), nil
	})
}
