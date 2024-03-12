package admin_routes

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type GetProfileDetails struct {
	Id string `json:"user_id" validate:"required"`
}

func (h *AdminRouteHandler) GetProfileById(c echo.Context) error {

	// Create Get Session Details
	var getProfileDetails GetProfileDetails
	if err := c.Bind(&getProfileDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(getProfileDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Acceptable Role, Grab User
	obj, err := primitive.ObjectIDFromHex(getProfileDetails.Id)
	profile, err := h.profileDBService.GetProfileByUserId(obj)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, profile)
}
