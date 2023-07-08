package authentication_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	authentication_db "resume-review-api/authentication/database"
)

type ForgotPasswordDetails struct {
	Email string `json:"email_address" validate:"required"`
}

func ForgotPassword(c echo.Context) error {

	// Create Forgot Password Details
	var forgotPasswordDetails ForgotPasswordDetails
	if err := c.Bind(&forgotPasswordDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(forgotPasswordDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Add Forgot Password to Database
	err := authentication_db.CreateForgotPassword(c, forgotPasswordDetails.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Return Success
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
