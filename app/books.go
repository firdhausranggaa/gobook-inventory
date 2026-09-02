package app

import (
	"book-inventory/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Handler {
	return Handler{DB: db}
}

func (h *Handler) GetBooks(c *gin.Context) {
	var books []models.Books

	// 1. Ambil query parameter dari URL
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	search := c.Query("search")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// 2. Membangun Query GORM
	query := h.DB.Model(&models.Books{})

	// Menggunakan ILIKE (PostgreSQL) agar case-insensitive
	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("title ILIKE ? OR author ILIKE ?", searchKeyword, searchKeyword)
	}

	// 3. Menghitung total data keseluruhan (untuk metadata)
	var total int
	query.Count(&total)

	// 4. Eksekusi query final dengan Limit dan Offset
	query.Limit(limit).Offset(offset).Order("id asc").Find(&books)

	// 5. Mengembalikan JSON beserta metadata pagination
	c.JSON(http.StatusOK, gin.H{
		"data": books,
		"meta": gin.H{
			"total_data":  total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

func (h *Handler) GetBookById(c *gin.Context) {
	bookId := c.Param("id")
	var book models.Books

	if h.DB.Find(&book, bookId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (h *Handler) PostBook(c *gin.Context) {
	var book models.Books
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	h.DB.Create(&book)
	c.JSON(http.StatusCreated, gin.H{"message": "Buku berhasil ditambahkan", "data": book})
}

func (h *Handler) PutBook(c *gin.Context) {
	var book models.Books
	bookId := c.Param("id")

	if h.DB.Find(&book, bookId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	var reqBook models.Books
	if err := c.ShouldBindJSON(&reqBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Model(&book).Update(reqBook)
	c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil diperbarui", "data": book})
}

func (h *Handler) DeleteBook(c *gin.Context) {
	var book models.Books
	bookId := c.Param("id")

	if h.DB.Find(&book, bookId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	h.DB.Delete(&book, bookId)
	c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dihapus"})
}