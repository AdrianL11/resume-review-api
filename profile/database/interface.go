package profile_db

import (
	"github.com/labstack/echo/v4"
	"resume-review-api/util"
	"resume-review-api/util/resume_ai_env"
)

type ProfileDBService struct {
	serverSettings resume_ai_env.ServerSettings
}

var _ util.RouteHandler = &ProfileDBService{}

func NewProfileDBService(serverSettings resume_ai_env.ServerSettings) *ProfileDBService {
	return &ProfileDBService{
		serverSettings: serverSettings,
	}
}

func (s *ProfileDBService) RegisterRoutes(c *echo.Echo, RequireAuthedSessionMiddleware echo.MiddlewareFunc) {

}
