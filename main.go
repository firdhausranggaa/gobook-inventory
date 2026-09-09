package main

import (
	"book-inventory/app"
	"book-inventory/auth"
	"book-inventory/db"
	"book-inventory/middleware"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	// gin.SetMode(gin.ReleaseMode) // Buka komentar ini saat aplikasi di-deploy ke server asli

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	appHandler := app.New(database)
	authHandler := auth.New(database)

	api := r.Group("/api")
	{
		api.POST("/register", authHandler.RegisterHandler)
		api.POST("/login", authHandler.LoginHandler)

		protected := api.Group("/")
		protected.Use(middleware.AuthValid)
		{
			protected.GET("/books", appHandler.GetBooks)
			protected.GET("/books/:id", appHandler.GetBookById)
			protected.POST("/borrow", appHandler.BorrowBook)
			protected.POST("/return/:id", appHandler.ReturnBook)
			protected.GET("/borrowings/me", appHandler.GetMyBorrowings) // Relasi transaksi user

			adminOnly := protected.Group("/")
			adminOnly.Use(middleware.RequireAdmin)
			{
				adminOnly.POST("/books", appHandler.PostBook)
				adminOnly.PUT("/books/:id", appHandler.PutBook)
				adminOnly.DELETE("/books/:id", appHandler.DeleteBook)
			}
		}
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gagal menjalankan server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Mematikan server SIMB secara perlahan...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server dipaksa mati karena timeout:", err)
	}

	log.Println("Server berhasil dimatikan dengan aman tanpa merusak data.")
}