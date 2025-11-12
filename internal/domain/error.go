package domain

import "errors"

var (
	ErrTeamExists   = errors.New("team already exists")
	ErrTeamNotFound = errors.New("team not found")
	ErrUserNotFound = errors.New("user not found")

	ErrPRExists   = errors.New("pr exists")
	ErrPRNotFound = errors.New("pr not found")
	ErrPRMerged   = errors.New("pr merged")

	ErrNotAssigned    = errors.New("not assigned")
	ErrNoCandidate    = errors.New("no candidate")
	ErrAuthorNotFound = errors.New("author not found")
	ErrInvalidRequest = errors.New("invalid request")
	ErrDB             = errors.New("db error")
)
