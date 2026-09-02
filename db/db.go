package db

import (
	"book-inventory/models"
	_ "database/sql"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	db.AutoMigrate(&models.Books{})

	data := models.Books{}
	if db.Find(&data).RecordNotFound() {
		seederBook(db)
	}
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
	}}

	for _, v := range data {
		db.Create(&v)
	}

}
