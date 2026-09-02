package app

import (
	"book-inventory/models"
	"net/http"

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
	h.DB.Find(&books)
	c.JSON(http.StatusOK, gin.H{"data": books})
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