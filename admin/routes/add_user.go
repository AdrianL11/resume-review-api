package admin_routes

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"resume-review-api/mongodb"
	"resume-review-api/resume_ai_middleware"
)

type AddUserDetails struct {
	EmailAddress string `json:"email_address" validate:"required"`
	Role         string `json:"role" validate:"required"`
}

func (h *AdminRouteHandler) AddUser(c echo.Context) error {

	// Create Add User Details
	var addUserDetails AddUserDetails
	if err := c.Bind(&addUserDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(addUserDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	viewerProfile, ok := c.Get(resume_ai_middleware.UserSessionProfile).(mongodb.Profile)
	if !ok {
		return echo.ErrInternalServerError
	}

	// Acceptable Role, Can User Add This Role?
	if viewerProfile.Role != mongodb.OwnerRole {
		if mongodb.Role(addUserDetails.Role) == mongodb.OwnerRole || mongodb.Role(addUserDetails.Role) == mongodb.Administrator {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("unable to add role"))
		}
	}

	// Can Add User, Add User
	err := h.adminDBService.AddUser(viewerProfile.ID, addUserDetails.EmailAddress, mongodb.Role(addUserDetails.Role))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
