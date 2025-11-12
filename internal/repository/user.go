package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/RomanGhost/pull-request-avito-test/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) UpsertUser(ctx context.Context, u *domain.User) error {
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(u).Error; err != nil {
		slog.Error("UpsertUser db error", "err", err)
		return domain.ErrDB
	}
	return nil
}

func (r *UserRepository) SetUserActive(ctx context.Context, userID string, active bool) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("SetUserActive find error", "err", err)
		return nil, domain.ErrDB
	}
	user.IsActive = active
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		slog.Error("SetUserActive save error", "err", err)
		return nil, domain.ErrDB
	}
	return &user, nil
}

func (r *UserRepository) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("GetUser db error", "err", err)
		return nil, domain.ErrDB
	}
	return &user, nil
}

func (r *UserRepository) GetActiveTeamMembers(ctx context.Context, teamName, excludeUserID string) ([]domain.User, error) {
	var users []domain.User
	q := r.db.WithContext(ctx).Where("team_name = ? AND is_active = true", teamName)
	if excludeUserID != "" {
		q = q.Where("user_id != ?", excludeUserID)
	}
	if err := q.Order("user_id").Find(&users).Error; err != nil {
		slog.Error("GetActiveTeamMembers db error", "err", err)
		return nil, domain.ErrDB
	}
	return users, nil
}
