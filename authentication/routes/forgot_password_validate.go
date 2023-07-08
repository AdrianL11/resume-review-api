package authentication_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	authentication_db "resume-review-api/authentication/database"
)

type ForgotPasswordValidateDetails struct {
	Token string `json:"token" validate:"required"`
}

func ForgotPasswordValidate(c echo.Context) error {

	// Create Forgot Password Details
	var forgotPasswordValidateDetails ForgotPasswordValidateDetails
	if err := c.Bind(&forgotPasswordValidateDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(forgotPasswordValidateDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validated, Check
	ret := authentication_db.ForgotPasswordValidateToken(forgotPasswordValidateDetails.Token)

	return c.JSON(http.StatusOK, map[string]bool{
		"valid_token": ret,
	})
}
