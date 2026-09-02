package models

import (
	"time"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null" binding:"required"`
	Password  string    `json:"-" gorm:"not null"`
	Role      string    `json:"role" gorm:"default:'member'"`
	CreatedAt time.Time `json:"created_at"`
}

type JWTClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}