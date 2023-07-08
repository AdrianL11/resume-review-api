package admin_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	admin_db "resume-review-api/admin/database"
	"resume-review-api/mongodb"
	session_db "resume-review-api/session/database"
)

func GetUsers(c echo.Context) error {

	// Is Session Valid
	err := session_db.ValidateSession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Session is Valid, Get Current Profile
	profile, err := mongodb.GetProfileBySession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Check Role
	if profile.Role != mongodb.OwnerRole && profile.Role != mongodb.Administrator {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Acceptable Role, Grab Profiles
	profiles, err := admin_db.GetUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//Done
	return c.JSON(http.StatusOK, profiles)
}
