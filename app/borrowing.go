package app

import (
	"book-inventory/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) BorrowBook(c *gin.Context) {
	var input struct {
		BookID       int    `json:"book_id" binding:"required"`
		BorrowerName string `json:"borrower_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	// Memulai Database Transaction
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Cek ketersediaan buku
	var book models.Books
	if tx.First(&book, input.BookID).RecordNotFound() {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	if book.Stock <= 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok buku habis"})
		return
	}

	// 2. Kurangi stok buku
	book.Stock -= 1
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate stok buku"})
		return
	}

	// 3. Mencatat transaksi peminjaman
	borrowing := models.Borrowing{
		BookID:       input.BookID,
		BorrowerName: input.BorrowerName,
		BorrowDate:   time.Now(),
		Status:       "BORROWED",
	}

	if err := tx.Create(&borrowing).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat data peminjaman"})
		return
	}

	tx.Commit() // Menyimpan semua perubahan permanen ke database
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

	// 1. Mengecek data peminjaman
	var borrowing models.Borrowing
	if tx.First(&borrowing, borrowingId).RecordNotFound() {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Data peminjaman tidak ditemukan"})
		return
	}

	if borrowing.Status == "RETURNED" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Buku ini sudah dikembalikan sebelumnya"})
		return
	}

	// 2. Update status dan tanggal kembali
	now := time.Now()
	borrowing.Status = "RETURNED"
	borrowing.ReturnDate = &now

	if err := tx.Save(&borrowing).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status peminjaman"})
		return
	}

	// 3. Mengembalikan stok buku
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