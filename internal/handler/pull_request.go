package handler

import (
	"log/slog"
	"net/http"

	"github.com/RomanGhost/pull-request-avito-test/internal/domain"
	"github.com/RomanGhost/pull-request-avito-test/internal/handler/model"
	"github.com/RomanGhost/pull-request-avito-test/internal/service"
	"github.com/gin-gonic/gin"
)

type PRHandler struct {
	prService *service.PRService
}

func (h *PRHandler) CreatePR(c *gin.Context) {
	ctx := c.Request.Context()
	var req model.CreatePullRequestRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("CreatePR bind error", "err", err)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.CodeBadRequest,
			err.Error(),
		))
		return
	}

	domainPR := model.ToDomainPR(&req)
	createdPR, err := h.prService.CreatePR(ctx, domainPR)
	if err != nil {
		slog.Error("CreatePR service error", "err", err)
		switch err {
		case domain.ErrAuthorNotFound:
			c.JSON(http.StatusNotFound, model.NewErrorResponse(
				model.CodeNotFound,
				"author not found",
			))
		case domain.ErrPRExists:
			c.JSON(http.StatusConflict, model.NewErrorResponse(
				model.CodePRExists,
				"PR id already exists",
			))
		default:
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
				model.InternalError,
				"database error",
			))
		}
		return
	}

	pr := model.FromDomainPR(createdPR)
	c.JSON(http.StatusCreated, model.PullRequestResponse{PR: pr})
}

func (h *PRHandler) MergePR(c *gin.Context) {
	ctx := c.Request.Context()
	var req model.MergePullRequestRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("MergePR bind error", "err", err)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.CodeBadRequest,
			err.Error(),
		))
		return
	}

	domainPR, err := h.prService.MergePR(ctx, req.PullRequestID)
	if err != nil {
		slog.Error("MergePR service error", "err", err)
		if err == domain.ErrPRNotFound {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(
				model.CodeNotFound,
				"PR not found",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.InternalError,
			"database error",
		))
		return
	}

	pr := model.FromDomainPR(domainPR)
	c.JSON(http.StatusOK, model.PullRequestResponse{PR: pr})
}

func (h *PRHandler) ReassignPR(c *gin.Context) {
	ctx := c.Request.Context()
	var req model.ReassignReviewerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("ReassignPR bind error", "err", err)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.CodeBadRequest,
			err.Error(),
		))
		return
	}

	newReviewer, domainPR, err := h.prService.ReassignPR(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		slog.Error("ReassignPR service error", "err", err)
		switch err {
		case domain.ErrPRNotFound, domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, model.NewErrorResponse(
				model.CodeNotFound,
				"PR or user not found",
			))
		case domain.ErrPRMerged:
			c.JSON(http.StatusConflict, model.NewErrorResponse(
				model.CodePRMerged,
				"cannot reassign on merged PR",
			))
		case domain.ErrNotAssigned:
			c.JSON(http.StatusConflict, model.NewErrorResponse(
				model.CodeNotAssigned,
				"reviewer is not assigned to this PR",
			))
		case domain.ErrNoCandidate:
			c.JSON(http.StatusConflict, model.NewErrorResponse(
				model.CodeNoCandidate,
				"no active replacement candidate in team",
			))
		default:
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
				model.InternalError,
				"database error",
			))
		}
		return
	}

	pr := model.FromDomainPR(domainPR)
	c.JSON(http.StatusOK, model.ReassignReviewerResponse{
		PR:         pr,
		ReplacedBy: newReviewer,
	})
}
