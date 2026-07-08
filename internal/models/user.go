package models

type User struct {
	ID           uint `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex"`
	PasswordHash string

	Ideas    []Idea
	Comments []Comment
}