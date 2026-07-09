package models

import "time"


type Session struct {
	ID string `gorm:"primaryKey"`

	UserID uint
	User User

	CreatedAt time.Time
	ExpiresAt time.Time
}