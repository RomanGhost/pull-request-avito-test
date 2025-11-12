package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/RomanGhost/pull-request-avito-test/internal/domain"
	"gorm.io/gorm"
)

type PRRepository struct{ db *gorm.DB }

func NewPRRepository(db *gorm.DB) *PRRepository { return &PRRepository{db: db} }

func (r *PRRepository) CreatePR(ctx context.Context, pr *domain.PullRequest) error {
	if err := r.db.WithContext(ctx).Create(pr).Error; err != nil {
		slog.Error("CreatePR db error", "err", err)
		return domain.ErrDB
	}
	return nil
}

func (r *PRRepository) GetPR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	if err := r.db.WithContext(ctx).Preload("Reviewers").First(&pr, "pull_request_id = ?", prID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("GetPR db error", "err", err)
		return nil, domain.ErrDB
	}
	return &pr, nil
}

func (r *PRRepository) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	if err := r.db.WithContext(ctx).First(&pr, "pull_request_id = ?", prID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("MergePR find error", "err", err)
		return nil, domain.ErrDB
	}
	if pr.Status != "MERGED" {
		now := time.Now()
		pr.Status = "MERGED"
		pr.MergedAt = &now
		if err := r.db.WithContext(ctx).Save(&pr).Error; err != nil {
			slog.Error("MergePR save error", "err", err)
			return nil, domain.ErrDB
		}
	}
	return &pr, nil
}

func (r *PRRepository) AssignReviewers(ctx context.Context, prID string, userIDs []string) error {
	for _, uid := range userIDs {
		re := domain.PRReviewer{PullRequestID: prID, UserID: uid}
		if err := r.db.WithContext(ctx).Create(&re).Error; err != nil {
			slog.Error("AssignReviewers db error", "err", err)
			return domain.ErrDB
		}
	}
	return nil
}

func (r *PRRepository) GetReviewers(ctx context.Context, prID string) ([]domain.PRReviewer, error) {
	var res []domain.PRReviewer
	if err := r.db.WithContext(ctx).Where("pull_request_id = ?", prID).Find(&res).Error; err != nil {
		slog.Error("GetReviewers db error", "err", err)
		return nil, domain.ErrDB
	}
	return res, nil
}

func (r *PRRepository) IsReviewerAssigned(ctx context.Context, prID, userID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.PRReviewer{}).Where("pull_request_id = ? AND user_id = ?", prID, userID).Count(&count).Error; err != nil {
		slog.Error("IsReviewerAssigned db error", "err", err)
		return false, domain.ErrDB
	}
	return count > 0, nil
}

func (r *PRRepository) RemoveReviewer(ctx context.Context, prID, userID string) error {
	if err := r.db.WithContext(ctx).Where("pull_request_id = ? AND user_id = ?", prID, userID).Delete(&domain.PRReviewer{}).Error; err != nil {
		slog.Error("RemoveReviewer db error", "err", err)
		return domain.ErrDB
	}
	return nil
}

func (r *PRRepository) AddReviewer(ctx context.Context, prID, userID string) error {
	if err := r.db.WithContext(ctx).Create(&domain.PRReviewer{PullRequestID: prID, UserID: userID}).Error; err != nil {
		slog.Error("AddReviewer db error", "err", err)
		return domain.ErrDB
	}
	return nil
}

func (r *PRRepository) GetPRsForReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	var prs []domain.PullRequest
	if err := r.db.WithContext(ctx).Joins("JOIN pr_reviewers ON pr_reviewers.pull_request_id = pull_requests.pull_request_id").Where("pr_reviewers.user_id = ?", userID).Order("created_at desc").Find(&prs).Error; err != nil {
		slog.Error("GetPRsForReviewer db error", "err", err)
		return nil, domain.ErrDB
	}
	return prs, nil
}
