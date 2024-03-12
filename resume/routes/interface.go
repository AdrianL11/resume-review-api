package resume_routes

import (
	"github.com/labstack/echo/v4"
	resume_db "resume-review-api/resume/database"
	"resume-review-api/util"
	"resume-review-api/util/resume_ai_env"
)

type ResumeRouteHandler struct {
	serverSettings  resume_ai_env.ServerSettings
	resumeDBService *resume_db.ResumeDBService
}

func NewResumeRouteHandler(
	serverSettings resume_ai_env.ServerSettings,
	resumeDBService *resume_db.ResumeDBService,
) *ResumeRouteHandler {
	return &ResumeRouteHandler{
		serverSettings:  serverSettings,
		resumeDBService: resumeDBService,
	}
}

var _ util.RouteHandler = &ResumeRouteHandler{}

func (h *ResumeRouteHandler) RegisterRoutes(e *echo.Echo, requireAuthedSessionMiddleware echo.MiddlewareFunc) {
	authedRoutes := e.Group("", requireAuthedSessionMiddleware)
	authedRoutes.POST("/resume/review", h.ReviewResume)
	authedRoutes.GET("/resume/counts", h.GetResumeCountInfo)
}
