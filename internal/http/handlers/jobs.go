package handlers

import (
	"encoding/json"
	"job-queue/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type JobHandler struct {
	service *service.JobService
}

func NewJobHandler(service *service.JobService) *JobHandler {
	return &JobHandler{service: service}
}

func (h *JobHandler) CreateJob(c echo.Context) error {
	var req struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	job, err := h.service.CreateJob(c.Request().Context(), req.Type, req.Payload)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create job")
	}

	return c.JSON(http.StatusAccepted, job)
}

func (h *JobHandler) GetJob(c echo.Context) error {
	jobID := c.Param("id")
	if jobID == "" {
		return c.String(http.StatusBadRequest, "Job ID is required")
	}

	job, err := h.service.GetJob(c.Request().Context(), jobID)
	if err != nil {
		return c.String(http.StatusNotFound, "Job not found")
	}

	return c.JSON(http.StatusOK, job)
}
