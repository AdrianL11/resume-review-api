package authentication_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	authentication_db "resume-review-api/authentication/database"
)

type ResetPasswordDetails struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func ResetPassword(c echo.Context) error {

	// Create Reset Password Details
	var resetPasswordDetails ResetPasswordDetails
	if err := c.Bind(&resetPasswordDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(resetPasswordDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validated, Reset Password
	err := authentication_db.ResetPassword(resetPasswordDetails.Token, resetPasswordDetails.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Return Success
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
