package admin_routes

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"resume-review-api/mongodb"
	session_db "resume-review-api/session/database"
)

type GetSessionsDetails struct {
	Id string `json:"user_id" validate:"required"`
}

func GetUserSessions(c echo.Context) error {

	// Create Get Session Details
	var getSessionsDetails GetSessionsDetails
	if err := c.Bind(&getSessionsDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(getSessionsDetails); err != nil {
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

	// Acceptable Role, Grab Sessions
	obj, _ := primitive.ObjectIDFromHex(getSessionsDetails.Id)
	sessList, err := session_db.GetSessionsById(obj)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//Return List
	return c.JSON(http.StatusOK, sessList)
}
