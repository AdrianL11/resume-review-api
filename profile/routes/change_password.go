package profile_routes

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ChangePasswordDetails struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

func (h *ProfileRouteHandler) ChangePassword(c echo.Context) error {

	// Create New User Check Details
	var changePasswordDetails ChangePasswordDetails
	if err := c.Bind(&changePasswordDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(changePasswordDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Change Password
	sess, _ := session.Get(h.serverSettings.SessionCookieName, c)
	obj, err := h.profileService.GetUserIdByEmail(sess.Values["email_address"].(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.profileService.ChangePassword(obj, changePasswordDetails.OldPassword, changePasswordDetails.NewPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
