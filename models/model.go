package models

import "time"

type Books struct {
	ID          int        `json:"id" form:"id" gorm:"primaryKey"`
	Title       string     `json:"title" form:"title" binding:"required"`
	Author      string     `json:"author" form:"author" binding:"required"`
	Description string     `json:"description" form:"description" binding:"required"`
	Stock       int        `json:"stock" form:"stock" binding:"required"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}