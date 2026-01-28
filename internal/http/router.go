package http

import (
	"job-queue/internal/http/handlers"

	"github.com/labstack/echo/v4"
)

func NewRouter(jobHandler *handlers.JobHandler) *echo.Echo {
	e := echo.New()

	e.POST("/jobs", jobHandler.CreateJob)
	e.GET("/jobs/:id", jobHandler.GetJob)

	return e
}
