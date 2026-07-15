package database

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"pitch.ideas/internal/database/models"
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
	return err
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

// InitDB initializes the database and creates tables
func InitDB() error {
	if err:= Init(); err != nil {
		return fmt.Errorf("Failed to initialize DB: %w", err)
	}
	// Create tables
	if err := Migrate(); err != nil {
		return fmt.Errorf("Failed to migrate DB: %w", err)
	}
	add_test_pitches()
	fmt.Println("Database and tables created!")
	return nil
}


func add_test_user() models.User {
	user, err := CreateUser("testuser", "hashedpassword")
	if err != nil {
		panic("Failed to create test user")
	}
	return *user
}

func add_test_pitch(title string, topic string, description string, vote_amount int, user_id uint) {
	id, err := CreateIdea(title, topic, description, 1)
	if err != nil || id == nil {
		panic("Failed to create test pitch")
	}
	VoteIdea(*id, user_id, vote_amount)
}

func add_test_pitches() {
	count, err := IdeaCount()
	if err != nil {
		panic("Failed to count ideas")
	}
	if count > 0 {
		return // only fill if empty
	}
	test_user := add_test_user()

	add_test_pitch(
		"AI Meeting Assistant",
		"AI",
		"AI that joins meetings, creates summaries, and automatically generates tasks.",
		421,
		test_user.ID,
	)
	add_test_pitch(
        "BudgetFlow",
        "FinTech",
        "Modern financial planning platform built for freelancers and creators.",
        312,
		test_user.ID,
    )
    add_test_pitch(
        "MedTrack",
        "HealthTech",
        "Patient monitoring system that helps clinics reduce administrative work.",
        198,
		test_user.ID,
    )
}





// IdeaCount returns the total number of ideas
func IdeaCount() (int64, error) {
	var count int64
	if err := DB.Model(&models.Idea{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CreateUser creates a new user
func CreateUser(username, passwordHash string) (*models.User, error) {
	// Check if user already exists
	var existing models.User
	if err := DB.Where("username = ?", username).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("user already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	user := models.User{
		Username:     username,
		PasswordHash: passwordHash,
	}

	if err := DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) *models.User {
	var user models.User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil
	}
	return &user
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID uint) *models.User {
	var user models.User
	if err := DB.First(&user, userID).Error; err != nil {
		return nil
	}
	return &user
}

// CreateIdea creates a new idea/pitch
func CreateIdea(title, topic, description string, userID uint) (*uint, error) {
	idea := models.Idea{
		Title:       title,
		Topic:       topic,
		Description: description,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	if err := DB.Create(&idea).Error; err != nil {
		return nil, err
	}

	return &idea.ID, nil
}

// EditIdea edits an existing idea
func EditIdea(ideaID, userID uint, title, topic, description string) (bool, error) {
	var idea models.Idea
	if err := DB.First(&idea, ideaID).Error; err != nil {
		return false, err
	}

	if idea.UserID != userID {
		return false, fmt.Errorf("user is not the owner of this idea")
	}

	if err := DB.Model(&idea).Updates(map[string]interface{}{
		"title":       title,
		"topic":       topic,
		"description": description,
	}).Error; err != nil {
		return false, err
	}

	return true, nil
}

// DeleteIdea deletes an idea
func DeleteIdea(ideaID, userID uint) (bool, error) {
	var idea models.Idea
	if err := DB.First(&idea, ideaID).Error; err != nil {
		return false, err
	}

	if idea.UserID != userID {
		return false, fmt.Errorf("user is not the owner of this idea")
	}

	if err := DB.Delete(&idea).Error; err != nil {
		return false, err
	}

	return true, nil
}

type OrderBy int

const (
	OrderByLikesDESC = iota
	OrderByLikesASC
	OrderByCreationDESC
	OrderByCreationASC
)

var orderByName = map[OrderBy]string{
	OrderByLikesDESC:  "likes-desc",
	OrderByLikesASC: "likes-asc",
	OrderByCreationDESC: "creation-desc",
	OrderByCreationASC: "creation-asc",
	OrderByCommentsDESC: "comments-desc",
	OrderByCommentsASC: "comments-asc"
}

var orderByLookup = map[string]OrderBy{
	"likes-desc":    OrderByLikesDESC,
	"likes-asc":     OrderByLikesASC,
	"creation-desc": OrderByCreationDESC,
	"creation-asc":  OrderByCreationASC,
	"comments-desc": OrderByCommentsDESC,
	"comments-asc":  OrderByCommentsASC,
}

func ParseOrderBy(s string) OrderBy {
	if ob, ok := orderByLookup[s]; ok {
		return ob
	}
	return OrderByCreationDESC
}

func (ob OrderBy) String() string {
    return orderByName[ob]
}

// GetAllIdeasAsDicts retrieves all ideas as dictionaries
func GetAllIdeasAsDicts(limit int, orderBy OrderBy) ([]map[string]interface{}, error) {
	query := DB.
		Model(&models.Idea{}).
		Preload("User").
		Preload("Comments").
		Preload("VoteRecords")

	switch orderBy {

	case OrderByCreationDESC:
		query = query.Order("created_at DESC")

	case OrderByCreationASC:
		query = query.Order("created_at ASC")

	case OrderByLikesDESC:
		query = query.
			Joins("LEFT JOIN idea_votes ON idea_votes.idea_id = ideas.id").
			Group("ideas.id").
			Order("COALESCE(SUM(idea_votes.value),0) DESC")

	case OrderByLikesASC:
		query = query.
			Joins("LEFT JOIN idea_votes ON idea_votes.idea_id = ideas.id").
			Group("ideas.id").
			Order("COALESCE(SUM(idea_votes.value),0) ASC")

	case OrderByCommentsDESC:
		query = query.
			Joins("LEFT JOIN comments ON comments.idea_id = ideas.id").
			Group("ideas.id").
			Order("COUNT(comments.id) DESC")

	case OrderByCommentsASC:
		query = query.
			Joins("LEFT JOIN comments ON comments.idea_id = ideas.id").
			Group("ideas.id").
			Order("COUNT(comments.id) ASC")
	}

	var ideas []models.Idea
	if err := query.Limit(limit).Find(&ideas).Error; err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(ideas))
	for _, idea := range ideas {
		result = append(result, idea.ToDict())
	}

	return result, nil
}

// GetIdeaDict retrieves a single idea as a dictionary
func GetIdeaDict(ideaID uint) (map[string]interface{}, error) {
	var idea models.Idea
	if err := DB.Preload("User").Preload("Comments").Preload("VoteRecords").
		First(&idea, ideaID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return idea.ToDict(), nil
}

// VoteIdea adds or updates a vote on an idea
func VoteIdea(ideaID, userID uint, value int) (int, error) {
	var existingVote models.IdeaVote
	err := DB.Where("idea_id = ? AND user_id = ?", ideaID, userID).First(&existingVote).Error

	if err == nil {
		// Vote exists
		if existingVote.Value == value {
			// Undo vote
			if err := DB.Delete(&existingVote).Error; err != nil {
				return 0, err
			}
		} else {
			// Update vote
			if err := DB.Model(&existingVote).Update("value", value).Error; err != nil {
				return 0, err
			}
		}
	} else if err == gorm.ErrRecordNotFound {
		// Create new vote
		newVote := models.IdeaVote{
			IdeaID: ideaID,
			UserID: userID,
			Value:  value,
		}
		if err := DB.Create(&newVote).Error; err != nil {
			return 0, err
		}
	} else {
		return 0, err
	}

	// Calculate total votes for the idea
	var totalVotes int
	if err := DB.Model(&models.IdeaVote{}).
		Where("idea_id = ?", ideaID).
		Select("COALESCE(SUM(value), 0)").
		Row().Scan(&totalVotes); err != nil {
		return 0, err
	}

	return totalVotes, nil
}

// CreateComment creates a new comment
func CreateComment(ideaID, userID uint, content string) error {
	comment := models.Comment{
		IdeaID:    ideaID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		Votes:     0,
	}

	return DB.Create(&comment).Error
}

// EditComment edits an existing comment
func EditComment(commentID, userID uint, content string) error {
	var comment models.Comment
	if err := DB.First(&comment, commentID).Error; err != nil {
		return err
	}

	if comment.UserID != userID {
		return fmt.Errorf("user is not the owner of this comment")
	}

	if err := DB.Model(&comment).Update("content", content).Error; err != nil {
		return err
	}

	return nil
}

// DeleteComment deletes a comment
func DeleteComment(commentID, userID uint) error {
	var comment models.Comment
	if err := DB.First(&comment, commentID).Error; err != nil {
		return err
	}

	if comment.UserID != userID {
		return fmt.Errorf("user is not the owner of this comment")
	}

	if err := DB.Delete(&comment).Error; err != nil {
		return err
	}

	return nil
}

// GetCommentCount returns the number of comments for an idea
func GetCommentCount(ideaID uint) (int64, error) {
	var count int64
	if err := DB.Model(&models.Comment{}).Where("idea_id = ?", ideaID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetCommentsDict retrieves comments for an idea as dictionaries
func GetCommentsDict(ideaID uint, limit int) ([]map[string]interface{}, error) {
	var comments []models.Comment
	if err := DB.Preload("User").
		Where("idea_id = ?", ideaID).
		Order("created_at ASC").
		Limit(limit).
		Find(&comments).Error; err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, comment := range comments {
		result = append(result, comment.ToDict())
	}

	return result, nil
}

// UpdateCommentVotes updates the vote count for a comment
func UpdateCommentVotes(commentID uint, amount int) error {
	return DB.Model(&models.Comment{}).
		Where("id = ?", commentID).
		Update("votes", gorm.Expr("votes + ?", amount)).Error
}

// CreateSession creates a new session token
func CreateSession(userID uint, days int) (string, error) {
	sessionID := uuid.New().String()

	session := models.Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().AddDate(0, 0, days),
	}

	if err := DB.Create(&session).Error; err != nil {
		return "", err
	}

	return sessionID, nil
}

// GetUserBySession retrieves a user by session token
func GetUserBySession(sessionID string) *models.User {
	var session models.Session
	if err := DB.First(&session, "id = ?", sessionID).Error; err != nil {
		return nil
	}

	// Check if session has expired
	if session.ExpiresAt.Before(time.Now()) {
		return nil
	}

	// Get the user associated with the session
	var user models.User
	if err := DB.First(&user, session.UserID).Error; err != nil {
		return nil
	}

	return &user
}

// DeleteSession deletes a session
func DeleteSession(sessionID string) error {
	return DB.Delete(&models.Session{}, "id = ?", sessionID).Error
}

// CleanupSessions removes all expired sessions
func CleanupSessions() error {
	return DB.Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
}
