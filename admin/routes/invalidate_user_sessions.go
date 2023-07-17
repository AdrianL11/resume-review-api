package admin_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"resume-review-api/mongodb"
	session_db "resume-review-api/session/database"
)

type InvalidateSessionUserDetails struct {
	Id string `json:"user_id" validate:"required"`
}

func InvalidateUserSessions(c echo.Context) error {

	// Create Get Session Details
	var getSessionsDetails InvalidateSessionUserDetails
	if err := c.Bind(&getSessionsDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(getSessionsDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

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

	// Acceptable Role, Grab Sessions
	err = session_db.InvalidateAllSessions(getSessionsDetails.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
