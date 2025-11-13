package model

type TeamMember struct {
	UserID   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name" binding:"required"`
	Members  []TeamMember `json:"members" binding:"required"`
}

type TeamResponse struct {
	Team Team `json:"team"`
}

type GetTeamRequest struct {
	TeamName string `form:"team_name" binding:"required"`
}
