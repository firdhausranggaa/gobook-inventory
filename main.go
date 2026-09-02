package main

import (
	"book-inventory/app"
	"book-inventory/auth"
	"book-inventory/db"
	"book-inventory/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	r := gin.Default()
	
	// Injeksi database ke masing-masing handler
	appHandler := app.New(database)
	authHandler := auth.New(database)

	api := r.Group("/api")
	{
		api.POST("/register", authHandler.RegisterHandler)
		api.POST("/login", authHandler.LoginHandler)

		protected := api.Group("/")
		protected.Use(middleware.AuthValid)
		{
			// Semua user yang login (Admin & Member) bisa melihat buku dan meminjam
			protected.GET("/books", appHandler.GetBooks)
			protected.GET("/books/:id", appHandler.GetBookById)
			protected.POST("/borrow", appHandler.BorrowBook)
			protected.POST("/return/:id", appHandler.ReturnBook)

			// Khusus Admin yang bisa manipulasi data buku
			adminOnly := protected.Group("/")
			adminOnly.Use(middleware.RequireAdmin)
			{
				adminOnly.POST("/books", appHandler.PostBook)
				adminOnly.PUT("/books/:id", appHandler.PutBook)
				adminOnly.DELETE("/books/:id", appHandler.DeleteBook)
			}
		}
	}

	r.Run(":8080")
}