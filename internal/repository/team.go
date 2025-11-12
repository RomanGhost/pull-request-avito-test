package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/RomanGhost/pull-request-avito-test/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team *domain.Team) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&domain.Team{TeamName: team.TeamName}).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return domain.ErrTeamExists // теперь используем custom error
			}
			slog.Error("CreateTeam db error", "err", err)
			return domain.ErrDB
		}
		for i := range team.Members {
			u := &team.Members[i]
			u.TeamName = team.TeamName
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(u).Error; err != nil {
				slog.Error("Upsert member error", "err", err)
				return domain.ErrDB
			}
		}
		return nil
	})
}

func (r *TeamRepository) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	var t domain.Team
	if err := r.db.WithContext(ctx).Preload("Members").First(&t, "team_name = ?", teamName).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("GetTeam db error", "err", err)
		return nil, domain.ErrDB
	}
	return &t, nil
}

func (r *TeamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Team{}).Where("team_name = ?", teamName).Count(&count).Error; err != nil {
		slog.Error("TeamExists db error", "err", err)
		return false, domain.ErrDB
	}
	return count > 0, nil
}
