package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("x") // Considere mover isso para a configuração

type contextKey string

const UserIDKey contextKey = "user_id"

func ProtectedRoutes() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			if userIDFloat, ok := claims["user_id"].(float64); ok {
				userID = fmt.Sprintf("%.0f", userIDFloat)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
				return
			}
		}

		// Adicionar o userID ao contexto do Gin
		c.Set(string(UserIDKey), userID)

		// Continuar para o próximo handler
		c.Next()
	}
}