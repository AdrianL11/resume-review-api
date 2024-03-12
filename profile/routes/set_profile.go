package profile_routes

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SetProfileDetails struct {
	Token        string `json:"token" validate:"required"`
	Password     string `json:"password" validate:"required"`
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Country      string `json:"country" validate:"required"`
	ProfileImage string `json:"profile_image"`
}

func (h *ProfileRouteHandler) SetProfile(c echo.Context) error {

	// Create New User Check Details
	var setProfileDetails SetProfileDetails
	if err := c.Bind(&setProfileDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(setProfileDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validated, Check if First time Logging In
	if h.profileService.NewUserValidateToken(setProfileDetails.Token) == false {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid token information"))
	}

	// Allowed, Insert into Database
	err := h.profileService.SetProfile(
		setProfileDetails.Token,
		setProfileDetails.Password,
		setProfileDetails.FirstName,
		setProfileDetails.LastName,
		setProfileDetails.Country,
		setProfileDetails.ProfileImage,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
