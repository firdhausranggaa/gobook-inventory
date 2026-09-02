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
	r.LoadHTMLGlob("templates/*")

	handler := app.New(database)

	r.GET("/", auth.HomeHandler)
	r.GET("/login", auth.LoginGetHandler)
	r.POST("/login", auth.LoginPostHandler)

	r.GET("/books", middleware.AuthValid, handler.GetBooks)
	r.GET("/book/:id", middleware.AuthValid, handler.GetBookById)

	r.GET("/addBook", middleware.AuthValid, handler.AddBook)
	r.POST("/book", middleware.AuthValid, handler.PostBook)

	r.GET("/updateBook/:id", middleware.AuthValid, handler.UpdateBook)
	
	// HTML Forms do not support PUT and DELETE directly, so we map them via POST
	r.POST("/updateBook/:id", middleware.AuthValid, handler.PutBook)
	r.POST("/deleteBook/:id", middleware.AuthValid, handler.DeleteBook)

	r.Run(":8080")
}