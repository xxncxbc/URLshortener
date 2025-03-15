package main

import (
	"URLshortener/internal/link"
	"URLshortener/internal/stat"
	"URLshortener/internal/user"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&user.User{}, &link.Link{}, &stat.Stat{})
	if err != nil {
		panic(err)
	}
}
