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
	handler := app.New(database)

	// Grup routing API
	api := r.Group("/api")
	{
		api.POST("/login", auth.LoginHandler)

		// Grup routing yang dilindungi token JWT
		protected := api.Group("/")
		protected.Use(middleware.AuthValid)
		{
			protected.GET("/books", handler.GetBooks)
			protected.GET("/books/:id", handler.GetBookById)
			protected.POST("/books", handler.PostBook)
			protected.PUT("/books/:id", handler.PutBook)
			protected.DELETE("/books/:id", handler.DeleteBook)
		}
	}

	r.Run(":8080")
}