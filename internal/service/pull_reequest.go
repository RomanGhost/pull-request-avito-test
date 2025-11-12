package service

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/RomanGhost/pull-request-avito-test/internal/domain"
	"github.com/RomanGhost/pull-request-avito-test/internal/repository"
)

type PRService struct {
	prRepo   *repository.PRRepository
	userRepo *repository.UserRepository
}

func NewPRService(pr *repository.PRRepository, ur *repository.UserRepository) *PRService {
	return &PRService{prRepo: pr, userRepo: ur}
}

func (s *PRService) CreatePR(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	existing, err := s.prRepo.GetPR(ctx, pr.PullRequestID)
	if err != nil {
		slog.Error("CreatePR get existing error", "err", err)
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrPRExists
	}
	author, err := s.userRepo.GetUser(ctx, pr.AuthorID)
	if err != nil {
		slog.Error("CreatePR get author error", "err", err)
		return nil, err
	}
	if author == nil {
		return nil, domain.ErrAuthorNotFound
	}
	candidates, err := s.userRepo.GetActiveTeamMembers(ctx, author.TeamName, pr.AuthorID)
	if err != nil {
		slog.Error("CreatePR get candidates error", "err", err)
		return nil, err
	}
	ids := make([]string, 0, len(candidates))
	for _, u := range candidates {
		ids = append(ids, u.UserID)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	pick := []string{}
	for i := 0; i < len(ids) && i < 2; i++ {
		pick = append(pick, ids[i])
	}
	pr.Status = "OPEN"
	pr.CreatedAt = time.Now()
	if err := s.prRepo.CreatePR(ctx, pr); err != nil {
		slog.Error("CreatePR create error", "err", err)
		return nil, err
	}
	if len(pick) > 0 {
		if err := s.prRepo.AssignReviewers(ctx, pr.PullRequestID, pick); err != nil {
			slog.Error("CreatePR assign reviewers error", "err", err)
			return nil, err
		}
	}
	created, err := s.prRepo.GetPR(ctx, pr.PullRequestID)
	if err != nil {
		slog.Error("CreatePR get final error", "err", err)
		return nil, err
	}
	return created, nil
}

func (s *PRService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := s.prRepo.MergePR(ctx, prID)
	if err != nil {
		slog.Error("MergePR error", "err", err)
		return nil, err
	}
	if pr == nil {
		return nil, domain.ErrPRNotFound
	}
	return pr, nil
}

func (s *PRService) ReassignPR(ctx context.Context, prID, oldUserID string) (string, *domain.PullRequest, error) {
	pr, err := s.prRepo.GetPR(ctx, prID)
	if err != nil {
		slog.Error("ReassignPR get pr error", "err", err)
		return "", nil, err
	}
	if pr == nil {
		return "", nil, domain.ErrPRNotFound
	}
	if pr.Status == "MERGED" {
		return "", nil, domain.ErrPRMerged
	}
	assigned, err := s.prRepo.IsReviewerAssigned(ctx, prID, oldUserID)
	if err != nil {
		slog.Error("ReassignPR is assigned error", "err", err)
		return "", nil, err
	}
	if !assigned {
		return "", nil, domain.ErrNotAssigned
	}
	oldUser, err := s.userRepo.GetUser(ctx, oldUserID)
	if err != nil {
		slog.Error("ReassignPR get old user error", "err", err)
		return "", nil, err
	}
	if oldUser == nil {
		return "", nil, domain.ErrUserNotFound
	}
	candidates, err := s.userRepo.GetActiveTeamMembers(ctx, oldUser.TeamName, oldUser.UserID)
	if err != nil {
		slog.Error("ReassignPR get candidates error", "err", err)
		return "", nil, err
	}
	existingReviewers := map[string]bool{}
	for _, r := range pr.Reviewers {
		existingReviewers[r.UserID] = true
	}
	avail := []string{}
	for _, c := range candidates {
		if !existingReviewers[c.UserID] {
			avail = append(avail, c.UserID)
		}
	}
	if len(avail) == 0 {
		return "", nil, domain.ErrNoCandidate
	}
	rand.Seed(time.Now().UnixNano())
	newReviewer := avail[rand.Intn(len(avail))]
	if err := s.prRepo.RemoveReviewer(ctx, prID, oldUserID); err != nil {
		slog.Error("ReassignPR remove error", "err", err)
		return "", nil, err
	}
	if err := s.prRepo.AddReviewer(ctx, prID, newReviewer); err != nil {
		slog.Error("ReassignPR add error", "err", err)
		return "", nil, err
	}
	pr, err = s.prRepo.GetPR(ctx, prID)
	if err != nil {
		slog.Error("ReassignPR get final pr error", "err", err)
		return "", nil, err
	}
	return newReviewer, pr, nil
}

func (s *PRService) GetPR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	return s.prRepo.GetPR(ctx, prID)
}

func (s *PRService) GetPRsForReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	prs, err := s.prRepo.GetPRsForReviewer(ctx, userID)
	if err != nil {
		slog.Error("GetPRsForReviewer error", "err", err)
		return nil, err
	}
	return prs, nil
}
