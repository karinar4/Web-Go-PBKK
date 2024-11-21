package main

import (
	"fmt"
	"log"
    "os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
    "github.com/joho/godotenv"
)

type Album struct {
	ID     uint    `gorm:"primaryKey"`
	Title  string 
	Artist string 
	Price  float32
}

func main() {
    err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
    
    dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Open a connection to the database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	fmt.Println("Connected!")

	err = db.AutoMigrate(&Album{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	// Insert a new album
	newAlbum := Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	}
	err = addAlbum(db, &newAlbum)
	if err != nil {
		log.Fatal("Failed to add album:", err)
	}
	fmt.Printf("New album added with ID: %d\n", newAlbum.ID)

	// Fetch albums by artist
	albums, err := albumsByArtist(db, "Betty Carter")
	if err != nil {
		log.Fatal("Failed to fetch albums:", err)
	}
	fmt.Printf("Albums found: %+v\n", albums)

	// Fetch album by ID
	alb, err := albumByID(db, newAlbum.ID)
	if err != nil {
		log.Fatal("Failed to fetch album by ID:", err)
	}
	fmt.Printf("Album found: %+v\n", alb)
}

func albumsByArtist(db *gorm.DB, artist string) ([]Album, error) {
	var albums []Album
	result := db.Where("artist = ?", artist).Find(&albums)
	if result.Error != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", artist, result.Error)
	}
	return albums, nil
}

func albumByID(db *gorm.DB, id uint) (Album, error) {
	var album Album
	result := db.First(&album, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return album, fmt.Errorf("albumByID %d: no such album", id)
		}
		return album, fmt.Errorf("albumByID %d: %v", id, result.Error)
	}
	return album, nil
}

func addAlbum(db *gorm.DB, album *Album) error {
	result := db.Create(album)
	if result.Error != nil {
		return fmt.Errorf("addAlbum: %v", result.Error)
	}
	return nil
}
