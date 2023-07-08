package profile_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	profile_db "resume-review-api/profile/database"
)

type NewUserCheckDetails struct {
	Id string `json:"token" validate:"required"`
}

func NewUserValidate(c echo.Context) error {

	// Create New User Check Details
	var newUserCheckDetails NewUserCheckDetails
	if err := c.Bind(&newUserCheckDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(newUserCheckDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validated, Check if First time Logging In
	ret := profile_db.NewUserValidateToken(newUserCheckDetails.Id)

	return c.JSON(http.StatusOK, map[string]bool{
		"valid_token": ret,
	})
}
