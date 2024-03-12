package admin_routes

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type GetSessionsDetails struct {
	Id string `json:"user_id" validate:"required"`
}

func (h *AdminRouteHandler) GetUserSessions(c echo.Context) error {

	// Create Get Session Details
	var getSessionsDetails GetSessionsDetails
	if err := c.Bind(&getSessionsDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(getSessionsDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Acceptable Role, Grab Sessions
	obj, _ := primitive.ObjectIDFromHex(getSessionsDetails.Id)
	sessList, err := h.authDBService.GetSessionsById(obj)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//Return List
	return c.JSON(http.StatusOK, sessList)
}
