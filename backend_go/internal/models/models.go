package models

import (
	"database/sql"
	"time"
)

type Word struct {
	ID       int64           `json:"id"`
	Japanese string         `json:"japanese"`
	Romaji   string         `json:"romaji"`
	English  string         `json:"english"`
	Parts    sql.NullString `json:"-"`
	PartsMap map[string]interface{} `json:"parts,omitempty"`
}

type Group struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	WordCount int    `json:"word_count,omitempty"`
}

type StudyActivity struct {
	ID             int64     `json:"id"`
	StudySessionID int64     `json:"study_session_id"`
	GroupID        int64     `json:"group_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int64     `json:"study_activity_id"`
	GroupName       string    `json:"group_name,omitempty"`
	ActivityName    string    `json:"activity_name,omitempty"`
	ReviewItemCount int       `json:"review_items_count,omitempty"`
}

type WordReviewItem struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

type WordStats struct {
	CorrectCount int `json:"correct_count"`
	WrongCount   int `json:"wrong_count"`
}

type StudyProgress struct {
	TotalWordsStudied    int `json:"total_words_studied"`
	TotalAvailableWords int `json:"total_available_words"`
}

type QuickStats struct {
	SuccessRate        float64 `json:"success_rate"`
	TotalStudySessions int     `json:"total_study_sessions"`
	TotalActiveGroups  int     `json:"total_active_groups"`
	StudyStreakDays    int     `json:"study_streak_days"`
}

type Pagination struct {
	CurrentPage   int `json:"current_page"`
	TotalPages    int `json:"total_pages"`
	TotalItems    int `json:"total_items"`
	ItemsPerPage  int `json:"items_per_page"`
}

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}
