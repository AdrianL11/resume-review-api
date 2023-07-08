package admin_routes

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	admin_db "resume-review-api/admin/database"
	"resume-review-api/mongodb"
	session_db "resume-review-api/session/database"
)

type DeactivateUserDetails struct {
	Id string `json:"user_id" validate:"required"`
}

func DeactivateUser(c echo.Context) error {

	// Create Deactivate User Details
	var deactivateUserDetails DeactivateUserDetails
	if err := c.Bind(&deactivateUserDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(deactivateUserDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Is Session Valid
	err := session_db.ValidateSession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Session is Valid, Get Current Profile
	profile, err := mongodb.GetProfileBySession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Check Role
	currentRole := profile.Role
	if profile.Role != mongodb.OwnerRole && profile.Role != mongodb.Administrator {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Acceptable Role, Grab User
	obj, err := primitive.ObjectIDFromHex(deactivateUserDetails.Id)
	profile, err = mongodb.GetProfilebyUserId(obj)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if Current User can Change based on Role
	if currentRole != mongodb.OwnerRole {
		if currentRole == mongodb.Administrator && profile.Role != mongodb.User {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("unauthorized authority"))
		}
	}

	// Can Change, Deactivate
	err = admin_db.DeactivateUser(c, obj)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
