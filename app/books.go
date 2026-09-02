package app

import (
	"book-inventory/models"
	"fmt"
	"net/http"
	"strings"

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
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "Home Page",
		"payload": books,
		"auth":    c.Query("auth"),
	})
}

func (h *Handler) GetBookById(c *gin.Context) {
	bookId := c.Param("id")
	var book models.Books

	if h.DB.Find(&book, bookId).RecordNotFound() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.HTML(http.StatusOK, "book.html", gin.H{
		"title":   book.Title,
		"payload": book,
		"auth":    c.Query("auth"),
	})
}

func (h *Handler) AddBook(c *gin.Context) {
	c.HTML(http.StatusOK, "formBook.html", gin.H{
		"title": "Add Book",
		"auth":  c.Query("auth"),
	})
}

func (h *Handler) PostBook(c *gin.Context) {
	var books models.Books
	c.Bind(&books)
	h.DB.Create(&books)
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/books?auth=%s", c.PostForm("auth")))
}

func (h *Handler) UpdateBook(c *gin.Context) {
	var book models.Books
	bookId := c.Param("id")

	if h.DB.Find(&book, bookId).RecordNotFound() {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Book not found"})
		return
	}

	c.HTML(http.StatusOK, "formBook.html", gin.H{
		"title":   "Update Book",
		"payload": book,
		"auth":    c.Query("auth"),
	})
}

func (h *Handler) PutBook(c *gin.Context) {
	var book models.Books
	bookId := c.Param("id")

	if h.DB.Find(&book, bookId).RecordNotFound() {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Book not found"})
		return
	}

	// Deteksi HTTP PUT method via form[cite: 17]
	if method := strings.ToLower(c.PostForm("_method")); method == "put" {
		c.Bind(&book)
		h.DB.Model(&book).Update(book)
	}

	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/book/%s?auth=%s", bookId, c.PostForm("auth")))
}

func (h *Handler) DeleteBook(c *gin.Context) {
	var book models.Books
	bookId := c.Param("id")
	h.DB.Delete(&book, bookId)
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/books?auth=%s", c.PostForm("auth")))
}