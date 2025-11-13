package model

import (
	"github.com/RomanGhost/pull-request-avito-test/internal/domain"
)

// ToDomainTeam конвертирует Team в domain.Team
func ToDomainTeam(t *Team) *domain.Team {
	members := make([]domain.User, len(t.Members))
	for i, m := range t.Members {
		members[i] = domain.User{
			UserID:   m.UserID,
			Username: m.Username,
			TeamName: t.TeamName,
			IsActive: m.IsActive,
		}
	}
	return &domain.Team{
		TeamName: t.TeamName,
		Members:  members,
	}
}

// FromDomainTeam конвертирует domain.Team в Team
func FromDomainTeam(dt *domain.Team) Team {
	members := make([]TeamMember, len(dt.Members))
	for i, m := range dt.Members {
		members[i] = TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}
	return Team{
		TeamName: dt.TeamName,
		Members:  members,
	}
}

// FromDomainUser конвертирует domain.User в User
func FromDomainUser(du *domain.User) User {
	return User{
		UserID:   du.UserID,
		Username: du.Username,
		TeamName: du.TeamName,
		IsActive: du.IsActive,
	}
}

// ToDomainPR конвертирует CreatePullRequestRequest в domain.PullRequest
func ToDomainPR(req *CreatePullRequestRequest) *domain.PullRequest {
	return &domain.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
	}
}

// FromDomainPR конвертирует domain.PullRequest в PullRequest
func FromDomainPR(dpr *domain.PullRequest) PullRequest {
	// Извлекаем UserID из Reviewers
	reviewers := make([]string, 0, len(dpr.Reviewers))
	for _, r := range dpr.Reviewers {
		reviewers = append(reviewers, r.UserID)
	}

	return PullRequest{
		PullRequestID:     dpr.PullRequestID,
		PullRequestName:   dpr.PullRequestName,
		AuthorID:          dpr.AuthorID,
		Status:            dpr.Status,
		AssignedReviewers: reviewers,
		CreatedAt:         &dpr.CreatedAt,
		MergedAt:          dpr.MergedAt,
	}
}

// FromDomainPRList конвертирует []domain.PullRequest в []PullRequestShort
func FromDomainPRList(dprs []domain.PullRequest) []PullRequestShort {
	result := make([]PullRequestShort, len(dprs))
	for i, dpr := range dprs {
		result[i] = PullRequestShort{
			PullRequestID:   dpr.PullRequestID,
			PullRequestName: dpr.PullRequestName,
			AuthorID:        dpr.AuthorID,
			Status:          dpr.Status,
		}
	}
	return result
}
