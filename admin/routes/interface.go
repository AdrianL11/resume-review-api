package admin_routes

import (
	"github.com/labstack/echo/v4"
	admin_db "resume-review-api/admin/database"
	authentication_db "resume-review-api/authentication/database"
	"resume-review-api/mongodb"
	profile_db "resume-review-api/profile/database"
	"resume-review-api/resume_ai_middleware"
	"resume-review-api/util"
	"resume-review-api/util/resume_ai_env"
)

type AdminRouteHandler struct {
	serverSettings   resume_ai_env.ServerSettings
	adminDBService   *admin_db.ResumeAIAdminDBService
	profileDBService *profile_db.ProfileDBService
	authDBService    *authentication_db.ResumeAIAuthDBService
}

func NewAdminRouteHandler(
	serverSettings resume_ai_env.ServerSettings,
	adminDBService *admin_db.ResumeAIAdminDBService,
	profileDBService *profile_db.ProfileDBService,
	authDBService *authentication_db.ResumeAIAuthDBService,
) *AdminRouteHandler {
	return &AdminRouteHandler{
		serverSettings:   serverSettings,
		adminDBService:   adminDBService,
		profileDBService: profileDBService,
		authDBService:    authDBService,
	}
}

var _ util.RouteHandler = &AdminRouteHandler{}

func (h *AdminRouteHandler) RegisterRoutes(e *echo.Echo, requireAuthedSessionMiddleware echo.MiddlewareFunc) {
	requireOwnerOrAdminRole := resume_ai_middleware.RequireAuthedUserAnyRoles(h.serverSettings, mongodb.OwnerRole, mongodb.Administrator)

	authedRoutes := e.Group("", requireAuthedSessionMiddleware, requireOwnerOrAdminRole)
	authedRoutes.GET("/admin/get_profiles", h.GetUsers)
	authedRoutes.POST("/admin/get_user_sessions", h.GetUserSessions)
	authedRoutes.POST("/admin/invalidate_user_sessions", h.InvalidateUserSessions)
	authedRoutes.POST("/admin/get_profile", h.GetProfileById)
	authedRoutes.POST("/admin/update_profile", h.UpdateProfile)
	authedRoutes.POST("/admin/deactivate_user", h.DeactivateUser)
	authedRoutes.POST("/admin/add_user", h.AddUser)
}
