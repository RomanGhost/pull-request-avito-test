package handler

import (
	"github.com/RomanGhost/pull-request-avito-test/internal/repository"
	"github.com/RomanGhost/pull-request-avito-test/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handlers struct {
	userHandler *UserHandler
	teamHandler *TeamHandler
	prHandler   *PRHandler
	pingDB      *gorm.DB
}

func RegisterHandlers(db *gorm.DB) *Handlers {
	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	prRepo := repository.NewPRRepository(db)

	teamSvc := service.NewTeamService(teamRepo)
	userSvc := service.NewUserService(userRepo)
	prSvc := service.NewPRService(prRepo, userRepo)

	teamHandler := TeamHandler{
		teamService: teamSvc,
	}
	prHandler := PRHandler{
		prService: prSvc,
	}
	userHandler := UserHandler{
		userService: userSvc,
		prService:   prSvc,
	}

	return &Handlers{
		userHandler: &userHandler,
		teamHandler: &teamHandler,
		prHandler:   &prHandler,
		pingDB:      db,
	}
}

func (h *Handlers) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.HealthCheck)

	team := r.Group("/team")
	{
		team.POST("/add", h.teamHandler.AddTeam)
		team.GET("/get", h.teamHandler.GetTeam)
	}

	user := r.Group("/users")
	{
		user.POST("/setIsActive", h.userHandler.SetUserActive)
		user.GET("/getReview", h.userHandler.GetUserReviews)
	}

	pr := r.Group("/pullRequest")
	{
		pr.POST("/create", h.prHandler.ReassignPR)
		pr.POST("/merge", h.prHandler.MergePR)
		pr.POST("/reassign", h.prHandler.ReassignPR)
	}
}
