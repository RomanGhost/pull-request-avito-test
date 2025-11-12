package service

import (
	"context"
	"log/slog"

	"github.com/RomanGhost/pull-request-avito-test/internal/domain"
	"github.com/RomanGhost/pull-request-avito-test/internal/repository"
)

type UserService struct{ repo *repository.UserRepository }

func NewUserService(r *repository.UserRepository) *UserService { return &UserService{repo: r} }

func (s *UserService) UpsertUser(ctx context.Context, u *domain.User) error {
	if err := s.repo.UpsertUser(ctx, u); err != nil {
		slog.Error("UpsertUser error", "err", err)
		return err
	}
	return nil
}

func (s *UserService) SetUserActive(ctx context.Context, userID string, active bool) (*domain.User, error) {
	user, err := s.repo.SetUserActive(ctx, userID, active)
	if err != nil {
		slog.Error("SetUserActive error", "err", err)
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		slog.Error("GetUser error", "err", err)
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetActiveTeamMembers(ctx context.Context, teamName, exclude string) ([]domain.User, error) {
	members, err := s.repo.GetActiveTeamMembers(ctx, teamName, exclude)
	if err != nil {
		slog.Error("GetActiveTeamMembers error", "err", err)
		return nil, err
	}
	return members, nil
}
