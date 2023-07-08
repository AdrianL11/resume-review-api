package admin_routes

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"resume-review-api/mongodb"
	session_db "resume-review-api/session/database"
)

type GetProfileDetails struct {
	Id string `json:"user_id" validate:"required"`
}

func GetProfileById(c echo.Context) error {

	// Create Get Session Details
	var getProfileDetails GetProfileDetails
	if err := c.Bind(&getProfileDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(getProfileDetails); err != nil {
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
	if profile.Role != mongodb.OwnerRole && profile.Role != mongodb.Administrator {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Acceptable Role, Grab User
	obj, err := primitive.ObjectIDFromHex(getProfileDetails.Id)
	profile, err = mongodb.GetProfilebyUserId(obj)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, profile)
}
