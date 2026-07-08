package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"pitch.ideas/internal/models"
)

var DB *gorm.DB

func Init() error {
	url := os.Getenv("DATABASE_URL")

	var err error

	if url != "" {
		DB, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	} else {
		DB, err = gorm.Open(sqlite.Open("local_development.db"), &gorm.Config{})
	}
	return  err
}

func Migrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Idea{},
		&models.Comment{},
		&models.IdeaVote{},
		&models.Session{},
	)
}