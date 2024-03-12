package admin_routes

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"resume-review-api/mongodb"
	"resume-review-api/resume_ai_middleware"
)

type DeactivateUserDetails struct {
	Id string `json:"user_id" validate:"required"`
}

func (h *AdminRouteHandler) DeactivateUser(c echo.Context) error {

	// Create Deactivate User Details
	var deactivateUserDetails DeactivateUserDetails
	if err := c.Bind(&deactivateUserDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(deactivateUserDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Acceptable Role, Grab User
	obj, err := primitive.ObjectIDFromHex(deactivateUserDetails.Id)
	profile, err := h.profileDBService.GetProfileByUserId(obj)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	viewerProfile, ok := c.Get(resume_ai_middleware.UserSessionProfile).(mongodb.Profile)
	if !ok {
		return echo.ErrInternalServerError
	}

	// Check if Current User can Change based on Role
	if viewerProfile.Role != mongodb.OwnerRole {
		if viewerProfile.Role == mongodb.Administrator && profile.Role != mongodb.User {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("unauthorized authority"))
		}
	}

	// Can Change, Deactivate
	err = h.adminDBService.DeactivateUser(c, obj)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
