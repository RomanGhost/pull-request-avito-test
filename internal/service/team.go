package service

import (
	"context"
	"log/slog"

	"github.com/RomanGhost/pull-request-avito-test/internal/domain"
	"github.com/RomanGhost/pull-request-avito-test/internal/repository"
)

type TeamService struct{ repo *repository.TeamRepository }

func NewTeamService(r *repository.TeamRepository) *TeamService { return &TeamService{repo: r} }

func (s *TeamService) CreateTeam(ctx context.Context, team *domain.Team) error {
	exists, err := s.repo.TeamExists(ctx, team.TeamName)
	if err != nil {
		slog.Error("CreateTeam exists check error", "err", err)
		return err
	}
	if exists {
		return domain.ErrTeamExists
	}
	if err := s.repo.CreateTeam(ctx, team); err != nil {
		slog.Error("CreateTeam error", "err", err)
		return err
	}
	return nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	team, err := s.repo.GetTeam(ctx, teamName)
	if err != nil {
		slog.Error("GetTeam error", "err", err)
		return nil, err
	}
	if team == nil {
		return nil, domain.ErrTeamNotFound
	}
	return team, nil
}
