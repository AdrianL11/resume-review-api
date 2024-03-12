package profile_routes

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"resume-review-api/resume_ai_middleware"
)

func (h *ProfileRouteHandler) GetActiveSessions(c echo.Context) error {
	sessionID, ok := c.Get(resume_ai_middleware.UserSessionIDContextKey).(primitive.ObjectID)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Get Sessions
	sess, err := h.authService.GetSessionsById(sessionID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, sess)
}
