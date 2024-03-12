package resume_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"resume-review-api/mongodb"
	"resume-review-api/resume"
	"resume-review-api/resume_ai_middleware"
)

func (h *ResumeRouteHandler) GetResumeCountInfo(c echo.Context) error {
	viewerProfile, ok := c.Get(resume_ai_middleware.UserSessionProfile).(mongodb.Profile)
	if !ok {
		return echo.ErrInternalServerError
	}

	count, average, err := h.resumeDBService.GetResumeInfo(viewerProfile.ID)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, resume.ReturnResponseResumeInfo{
		Count:        count,
		ResponseTime: average,
	})
}
