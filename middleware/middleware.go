package middleware

import (
	"book-inventory/models"
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid"})
		c.Abort()
		return
	}

	tokenString := parts[1]
	claims := &models.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, valid := t.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("algoritma tidak valid")
		}
		return []byte(os.Getenv("SUPER_SECRET")), nil
	})

	if err == nil && token.Valid {
		// Menyimpan data user ke context agar bisa dibaca oleh fungsi handler
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kadaluarsa"})
		c.Abort()
	}
}

// RequireAdmin memblokir akses jika role bukan 'admin'
func RequireAdmin(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Hanya admin yang diizinkan"})
		c.Abort()
		return
	}
	c.Next()
}