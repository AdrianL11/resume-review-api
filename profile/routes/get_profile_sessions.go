package profile_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"resume-review-api/mongodb"
	session_db "resume-review-api/session/database"
)

func GetActiveSessions(c echo.Context) error {

	// Is Session Valid
	err := session_db.ValidateSession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Get User Id
	obj, err := mongodb.GetProfileBySession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Get Sessions
	sess, err := session_db.GetSessionsById(obj.ID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, sess)
}
