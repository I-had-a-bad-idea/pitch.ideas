package models

import "time"

type Idea struct {
	ID          uint `gorm:"primaryKey"`

	Title       string
	Topic       string
	Description string `gorm:"type:text"`

	UserID uint
	User   User

	CreatedAt time.Time

	Comments    []Comment `gorm:"constraint:OnDelete:CASCADE;"`
	VoteRecords []IdeaVote `gorm:"constraint:OnDelete:CASCADE;"`
} 

func (i Idea) ToDict() map[string]interface{} {

	votes := 0
	votedBy := map[uint]int{}

	for _, v := range i.VoteRecords {
		votes += v.Value
		votedBy[v.UserID] = v.Value
	}

	return map[string]interface{}{
		"id": i.ID,
		"title": i.Title,
		"topic": i.Topic,
		"description": i.Description,
		"user_id": i.UserID,
		"user_name": i.User.Username,
		"created_at": i.CreatedAt.Format(time.RFC3339),
		"created_at_pretty": i.CreatedAt.Format("02 Jan 2006, 15:04"),
		"votes": votes,
		"voted_by_user": votedBy,
		"comment_count": len(i.Comments),
	}
}