package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers routes with handlers
func RegisterRoutes(r *echo.Echo) {
	r.GET("/health-check", HealthCheck)
}

// HealthCheck handler informs that api is up and running
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Up!")
}
