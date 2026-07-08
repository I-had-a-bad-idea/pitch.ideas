package models

import "time"


type Comment struct {
	ID uint `gorm:"primaryKey"`

	IdeaID uint
	UserID uint

	Idea Idea
	User User

	CreatedAt time.Time
	Content string

	Votes int
}


func (c Comment) ToDict() map[string]interface{} {

	return map[string]interface{}{
		"id": c.ID,
		"idea_id": c.IdeaID,
		"user_id": c.UserID,
		"user_name": c.User.Username,
		"content": c.Content,
		"votes": c.Votes,
		"created_at": c.CreatedAt.Format(time.RFC3339),
		"created_at_pretty": c.CreatedAt.Format("02 Jan 2006, 15:04"),
	}
}