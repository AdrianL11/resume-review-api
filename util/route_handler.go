package util

import "github.com/labstack/echo/v4"

type RouteHandler interface {
	RegisterRoutes(e *echo.Echo, requireAuthedSessionMiddleware echo.MiddlewareFunc)
}
