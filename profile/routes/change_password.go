package profile_routes

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"resume-review-api/mongodb"
	profile_db "resume-review-api/profile/database"
	session_db "resume-review-api/session/database"
)

type ChangePasswordDetails struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

func ChangePassword(c echo.Context) error {

	// Create New User Check Details
	var changePasswordDetails ChangePasswordDetails
	if err := c.Bind(&changePasswordDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(changePasswordDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Is Session Valid
	err := session_db.ValidateSession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Change Password
	sess, _ := session.Get(os.Getenv("session_name"), c)
	obj, err := mongodb.GetUserIdByEmail(sess.Values["email_address"].(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = profile_db.ChangePassword(obj, changePasswordDetails.OldPassword, changePasswordDetails.NewPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
