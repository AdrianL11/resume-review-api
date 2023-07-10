package authentication_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	session_db "resume-review-api/session/database"
)

func LoggedIn(c echo.Context) error {

	if session_db.ValidateSession(c) == nil {
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusUnauthorized)
	}
}
