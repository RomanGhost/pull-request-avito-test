package handler

import (
	"log/slog"
	"net/http"

	"github.com/RomanGhost/pull-request-avito-test/internal/domain"
	"github.com/RomanGhost/pull-request-avito-test/internal/handler/model"
	"github.com/RomanGhost/pull-request-avito-test/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
	prService   *service.PRService
}

func (h *UserHandler) SetUserActive(c *gin.Context) {
	ctx := c.Request.Context()
	var req model.SetIsActiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("SetUserActive bind error", "err", err)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.CodeBadRequest,
			err.Error(),
		))
		return
	}

	domainUser, err := h.userService.SetUserActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		slog.Error("SetUserActive service error", "err", err)
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(
				model.CodeNotFound,
				"user not found",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.InternalError,
			"database error",
		))
		return
	}

	user := model.FromDomainUser(domainUser)
	c.JSON(http.StatusOK, model.UserResponse{User: user})
}

func (h *UserHandler) GetUserReviews(c *gin.Context) {
	ctx := c.Request.Context()

	var req model.GetUserReviewRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.CodeBadRequest,
			"user_id is required",
		))
		return
	}

	domainPRs, err := h.prService.GetPRsForReviewer(ctx, req.UserID)
	if err != nil {
		slog.Error("GetUserReviews service error", "err", err)
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.InternalError,
			"database error",
		))
		return
	}

	prs := model.FromDomainPRList(domainPRs)
	c.JSON(http.StatusOK, model.UserReviewsResponse{
		UserID:       req.UserID,
		PullRequests: prs,
	})
}
