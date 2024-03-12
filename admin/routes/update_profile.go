package admin_routes

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"resume-review-api/mongodb"
	"resume-review-api/resume_ai_middleware"
)

type UpdateProfileDetails struct {
	Id           string `json:"user_id" validate:"required"`
	Email        string `json:"email_address" validate:"omitempty"`
	FirstName    string `json:"first_name" validate:"omitempty"`
	LastName     string `json:"last_name" validate:"omitempty"`
	Country      string `json:"country" validate:"omitempty"`
	ProfileImage string `json:"profile_image" validate:"omitempty"`
	Password     string `json:"password" validate:"omitempty"`
	Role         string `json:"role" validate:"omitempty"`
}

func (h *AdminRouteHandler) UpdateProfile(c echo.Context) error {

	// Create Update Profile Details
	var updateProfileDetails UpdateProfileDetails
	if err := c.Bind(&updateProfileDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(updateProfileDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	viewerProfile, ok := c.Get(resume_ai_middleware.UserSessionProfile).(mongodb.Profile)
	if !ok {
		return echo.ErrInternalServerError
	}

	// Acceptable Role, Grab User
	obj, err := primitive.ObjectIDFromHex(updateProfileDetails.Id)
	profile, err := h.profileDBService.GetProfileByUserId(obj)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if Current User can Change based on Role
	if viewerProfile.Role != mongodb.OwnerRole {
		if viewerProfile.Role == mongodb.Administrator && profile.Role != mongodb.User {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("unauthorized authority"))
		}
	}

	// Can Change
	err = h.adminDBService.UpdateProfile(
		updateProfileDetails.Id,
		updateProfileDetails.Email,
		updateProfileDetails.FirstName,
		updateProfileDetails.LastName,
		updateProfileDetails.Country,
		updateProfileDetails.ProfileImage,
		updateProfileDetails.Password,
		updateProfileDetails.Role,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
