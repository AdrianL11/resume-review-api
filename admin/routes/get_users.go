package admin_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *AdminRouteHandler) GetUsers(c echo.Context) error {
	// Acceptable Role, Grab Profiles
	profiles, err := h.adminDBService.GetUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//Done
	return c.JSON(http.StatusOK, profiles)
}
