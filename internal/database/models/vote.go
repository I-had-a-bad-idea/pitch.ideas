package models



type IdeaVote struct {
	ID uint `gorm:"primaryKey"`

	IdeaID uint `gorm:"uniqueIndex:vote_unique"`
	UserID uint `gorm:"uniqueIndex:vote_unique"`

	Value int `gorm:"default:1"`

	Idea Idea
	User User
}