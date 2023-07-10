package authentication_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	authentication_db "resume-review-api/authentication/database"
)

func LoggedIn(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"role": authentication_db.LoggedIn(c),
	})
}
