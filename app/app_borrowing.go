package app

import (
	"book-inventory/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) BorrowBook(c *gin.Context) {
	var input struct {
		BookID int `json:"book_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Buku diperlukan"})
		return
	}

	loggedInUser, _ := c.Get("username")

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var book models.Books
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&book, input.BookID).RecordNotFound() {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	if book.Stock <= 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok buku habis"})
		return
	}

	book.Stock -= 1
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate stok"})
		return
	}

	borrowing := models.Borrowing{
		BookID:       input.BookID,
		BorrowerName: loggedInUser.(string),
		BorrowDate:   time.Now(),
		Status:       "BORROWED",
	}

	if err := tx.Create(&borrowing).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat peminjaman"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"message": "Buku berhasil dipinjam", "data": borrowing})
}

func (h *Handler) ReturnBook(c *gin.Context) {
	borrowingId := c.Param("id")

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var borrowing models.Borrowing
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&borrowing, borrowingId).RecordNotFound() {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Data peminjaman tidak ditemukan"})
		return
	}

	if borrowing.Status == "RETURNED" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Buku ini sudah dikembalikan sebelumnya"})
		return
	}

	now := time.Now()
	borrowing.Status = "RETURNED"
	borrowing.ReturnDate = &now

	if err := tx.Save(&borrowing).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status peminjaman"})
		return
	}

	var book models.Books
	if tx.First(&book, borrowing.BookID).RecordNotFound() {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data buku tidak sinkron"})
		return
	}

	book.Stock += 1
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulihkan stok buku"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dikembalikan", "data": borrowing})
}

func (h *Handler) GetMyBorrowings(c *gin.Context) {
	loggedInUser, _ := c.Get("username")
	var borrowings []models.Borrowing

	if err := h.DB.Preload("Book").Where("borrower_name = ?", loggedInUser).Order("borrow_date desc").Find(&borrowings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat peminjaman"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil riwayat peminjaman",
		"data":    borrowings,
	})
}