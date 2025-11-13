package handler

import (
	"log/slog"
	"net/http"

	"github.com/RomanGhost/pull-request-avito-test/internal/domain"
	"github.com/RomanGhost/pull-request-avito-test/internal/handler/model"
	"github.com/RomanGhost/pull-request-avito-test/internal/service"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func (h *TeamHandler) AddTeam(c *gin.Context) {
	ctx := c.Request.Context()
	var req model.Team

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("AddTeam bind error", "err", err)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.CodeBadRequest,
			err.Error(),
		))
		return
	}

	// Конвертируем в domain
	domainTeam := model.ToDomainTeam(&req)

	if err := h.teamService.CreateTeam(ctx, domainTeam); err != nil {
		slog.Error("AddTeam service error", "err", err)
		if err == domain.ErrTeamExists {
			c.JSON(http.StatusBadRequest, model.NewErrorResponse(
				model.CodeTeamExists,
				"team_name already exists",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.InternalError,
			"database error",
		))
		return
	}

	// Конвертируем обратно для ответа
	responseTeam := model.FromDomainTeam(domainTeam)
	c.JSON(http.StatusCreated, model.TeamResponse{Team: responseTeam})
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	ctx := c.Request.Context()

	var req model.GetTeamRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.CodeBadRequest,
			"team_name required",
		))
		return
	}

	domainTeam, err := h.teamService.GetTeam(ctx, req.TeamName)
	if err != nil {
		slog.Error("GetTeam service error", "err", err)
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.InternalError,
			"database error",
		))
		return
	}

	if domainTeam == nil {
		c.JSON(http.StatusNotFound, model.NewErrorResponse(
			model.CodeNotFound,
			"team not found",
		))
		return
	}

	responseTeam := model.FromDomainTeam(domainTeam)
	c.JSON(http.StatusOK, responseTeam)
}
