package authentication_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ForgotPasswordValidateDetails struct {
	Token string `json:"token" validate:"required"`
}

func (h *AuthRouteHandler) ForgotPasswordValidate(c echo.Context) error {

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
	ret := h.authService.ForgotPasswordValidateToken(forgotPasswordValidateDetails.Token)

	return c.JSON(http.StatusOK, map[string]bool{
		"valid_token": ret,
	})
}
