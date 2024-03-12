package admin_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type InvalidateSessionUserDetails struct {
	Id string `json:"user_id" validate:"required"`
}

func (h *AdminRouteHandler) InvalidateUserSessions(c echo.Context) error {

	// Create Get Session Details
	var getSessionsDetails InvalidateSessionUserDetails
	if err := c.Bind(&getSessionsDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(getSessionsDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Acceptable Role, Grab Sessions
	err := h.authDBService.InvalidateAllSessions(getSessionsDetails.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
