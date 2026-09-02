package db

import (
	"book-inventory/models"
	_ "database/sql"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func InitDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error load env")
	}
	conn := os.Getenv("POSTGRES_URL")
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(db)
	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Books{}, &models.Borrowing{}, &models.User{})

	var book models.Books
	if db.Find(&book).RecordNotFound() {
		seederBook(db)
	}

	var user models.User
	if db.Find(&user).RecordNotFound() {
		seederUser(db)
	}
}

func seederUser(db *gorm.DB) {
	adminUser := os.Getenv("SUPER_USER")
	adminPass := os.Getenv("SUPER_PASS")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)

	admin := models.User{
		Username: adminUser,
		Password: string(hashedPassword),
		Role:     "admin",
	}
	db.Create(&admin)
}

func seederBook(db *gorm.DB) {
	data := []models.Books{{
		Title:       "Eternal Giants: Unveiling the Secrets of Dinosaurs",
		Author:      "Dr. Victoria Sinclair, Paleontologist Extraordinaire",
		Description: "Embark on a thrilling journey through time as Dr. Sinclair unravels the mysteries of dinosaurs and their incredible discoveries.",
		Stock:       10,
	}, {
		Title:       "Ingenious Ingenuity: A Chronicle of Engineering Marvels",
		Author:      "Professor Benjamin Inventorius",
		Description: "Dive into the world of awe-inspiring engineering feats and groundbreaking innovations that have shaped our modern civilization.",
		Stock:       5,
	}, {
		Title:       "Cosmic Odyssey: Exploring the Wonders of Space",
		Author:      "Captain Celeste Stargazer, Renowned Astronaut",
		Description: "Embark on a cosmic adventure with Captain Stargazer as she unveils the marvels of space, celestial bodies, and the mysteries beyond the stars.",
		Stock:       7,
	}, {
		Title:       "Timeless Chronicles: A Journey Through History's Tapestry",
		Author:      "Historian Extraordinaire, Dr. Eleanor Timewalker",
		Description: "Join Dr. Timewalker on an enthralling exploration of historical events, cultures, and the captivating evolution of civilizations throughout the ages.",
		Stock:       5,
	}, {
		Title:       "Generative AI & Neural Networks: A Practical Approach",
		Author:      "Dr. Alan Diffuser",
		Description: "Mastering image generation workflows, inpainting, and outpainting using PyTorch and Streamlit.",
		Stock:       15,
	}, {
		Title:       "Advanced Network Administration",
		Author:      "Router Guru",
		Description: "Comprehensive guide to configuring MikroTik WinBox, NAT, Trunk Ports, and GNS3 network topologies.",
		Stock:       8,
	}, {
		Title:       "Modern Frontend Mastery",
		Author:      "Evan Vite",
		Description: "Building scalable web applications and responsive UIs using Vue.js, React.js, and modern CSS.",
		Stock:       12,
	}, {
		Title:       "Object-Oriented Architecture for Enterprise",
		Author:      "James C. Java",
		Description: "Designing robust Point of Sale desktop applications using Java, C++, and relational SQL databases.",
		Stock:       6,
	}, {
		Title:       "Stylized 3D Aesthetics",
		Author:      "Render Artist",
		Description: "Techniques for achieving smooth matte plastic materials and cartoon toy figure design in 3D rendering.",
		Stock:       9,
	}, {
		Title:       "The Titanic Dataset: Machine Learning Fundamentals",
		Author:      "Data Science Academic",
		Description: "A deep dive into data analysis and academic research publishing using classic datasets.",
		Stock:       4,
	}}

	for _, v := range data {
		db.Create(&v)
	}
}