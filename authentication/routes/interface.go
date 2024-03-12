package authentication_routes

import (
	"github.com/labstack/echo/v4"
	admin_db "resume-review-api/admin/database"
	authentication_db "resume-review-api/authentication/database"
	"resume-review-api/util"
	"resume-review-api/util/resume_ai_env"
)

type AuthRouteHandler struct {
	serverSettings resume_ai_env.ServerSettings
	authService    *authentication_db.ResumeAIAuthDBService
	dbService      *admin_db.ResumeAIAdminDBService
}

func NewAuthRouteHandler(
	serverSettings resume_ai_env.ServerSettings,
	authService *authentication_db.ResumeAIAuthDBService,
	dbService *admin_db.ResumeAIAdminDBService,
) *AuthRouteHandler {
	return &AuthRouteHandler{
		serverSettings: serverSettings,
		authService:    authService,
		dbService:      dbService,
	}
}

var _ util.RouteHandler = &AuthRouteHandler{}

func (h *AuthRouteHandler) RegisterRoutes(e *echo.Echo, requireAuthedSessionMiddleware echo.MiddlewareFunc) {
	e.POST("/login", h.Login)
	e.POST("/forgot_password", h.ForgotPassword)
	e.POST("/forgot_password_validate", h.ForgotPasswordValidate)
	e.POST("/reset_password", h.ResetPassword)
	e.GET("/logged_in", h.LoggedIn)

	// Logout requires auth
	e.GET("/logout", h.Logout, requireAuthedSessionMiddleware)
}
