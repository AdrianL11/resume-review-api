package profile_routes

import (
	"github.com/labstack/echo/v4"
	authentication_db "resume-review-api/authentication/database"
	profile_db "resume-review-api/profile/database"
	"resume-review-api/util"
	"resume-review-api/util/resume_ai_env"
)

type ProfileRouteHandler struct {
	serverSettings resume_ai_env.ServerSettings
	profileService *profile_db.ProfileDBService
	authService    *authentication_db.ResumeAIAuthDBService
}

var _ util.RouteHandler = &ProfileRouteHandler{}

func NewProfileRouteHandler(
	serverSettings resume_ai_env.ServerSettings,
	profileService *profile_db.ProfileDBService,
	authService *authentication_db.ResumeAIAuthDBService,
) *ProfileRouteHandler {
	return &ProfileRouteHandler{
		serverSettings: serverSettings,
		profileService: profileService,
		authService:    authService,
	}
}

func (h *ProfileRouteHandler) RegisterRoutes(e *echo.Echo, requireAuthedSessionMiddleware echo.MiddlewareFunc) {
	authedRoutes := e.Group("", requireAuthedSessionMiddleware)
	authedRoutes.GET("/profile", h.GetProfile)
	authedRoutes.POST("/profile/set", h.SetProfile)
	authedRoutes.GET("/profile/sessions", h.GetActiveSessions)
	authedRoutes.POST("/profile/update", h.UpdateProfile)
	authedRoutes.POST("/profile/change_password", h.ChangePassword)
	authedRoutes.POST("/new_user_validate", h.NewUserValidate)
}
