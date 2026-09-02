package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthValid(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Header Authorization diperlukan"})
		c.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid. Gunakan: Bearer <token>"})
		c.Abort()
		return
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, valid := t.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("algoritma token tidak valid: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("SUPER_SECRET")), nil
	})

	if token != nil && err == nil {
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
	}
}