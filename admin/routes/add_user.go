package admin_routes

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	admin_db "resume-review-api/admin/database"
	"resume-review-api/mongodb"
	session_db "resume-review-api/session/database"
)

type AddUserDetails struct {
	EmailAddress string `json:"email_address" validate:"required"`
	Role         string `json:"role" validate:"required"`
}

func AddUser(c echo.Context) error {

	// Create Add User Details
	var addUserDetails AddUserDetails
	if err := c.Bind(&addUserDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(addUserDetails); err != nil {
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

	// Acceptable Role, Can User Add This Role?
	if profile.Role != mongodb.OwnerRole {
		if mongodb.Role(addUserDetails.Role) == mongodb.OwnerRole || mongodb.Role(addUserDetails.Role) == mongodb.Administrator {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("unable to add role"))
		}
	}

	// Can Add User, Add User
	err = admin_db.AddUser(profile.ID, addUserDetails.EmailAddress, mongodb.Role(addUserDetails.Role))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
