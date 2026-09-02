package models

import "time"

type Borrowing struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	BookID       int        `json:"book_id" binding:"required"`
	BorrowerName string     `json:"borrower_name" binding:"required"`
	BorrowDate   time.Time  `json:"borrow_date"`
	ReturnDate   *time.Time `json:"return_date"`
	Status       string     `json:"status"`
}