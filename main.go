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
		// Routing Publik (Tidak memerlukan token)
		api.POST("/register", authHandler.RegisterHandler)
		api.POST("/login", authHandler.LoginHandler)

		// Routing Privat (memerlukan token)
		protected := api.Group("/")
		protected.Use(middleware.AuthValid)
		{
			protected.GET("/books", appHandler.GetBooks)
			protected.GET("/books/:id", appHandler.GetBookById)
			protected.POST("/books", appHandler.PostBook)
			protected.PUT("/books/:id", appHandler.PutBook)
			protected.DELETE("/books/:id", appHandler.DeleteBook)

			protected.POST("/borrow", appHandler.BorrowBook)
			protected.POST("/return/:id", appHandler.ReturnBook)
		}
	}

	r.Run(":8080")
}