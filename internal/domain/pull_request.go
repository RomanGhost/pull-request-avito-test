package domain

import "time"

type PullRequest struct {
	PullRequestID   string       `gorm:"primaryKey;size:255"`
	PullRequestName string       `gorm:"size:255;not null"`
	AuthorID        string       `gorm:"size:255;not null"`
	Status          string       `gorm:"size:10;default:OPEN"`
	CreatedAt       time.Time    `json:"created_at"`
	MergedAt        *time.Time   `json:"merged_at,omitempty"`
	Reviewers       []PRReviewer `gorm:"foreignKey:PullRequestID;references:PullRequestID"`
}

type PRReviewer struct {
	PullRequestID string `gorm:"primaryKey;size:255"`
	UserID        string `gorm:"primaryKey;size:255"`
}
