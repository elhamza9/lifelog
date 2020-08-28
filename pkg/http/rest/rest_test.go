package rest_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elhamza90/lifelog/pkg/http/rest"
	"github.com/labstack/echo/v4"
)

var router *echo.Echo

func TestMain(m *testing.M) {
	log.Println("Setting Up Main")
	router = echo.New()
	rest.RegisterRoutes(router)
}

func TestHealthCheck(t *testing.T) {
	path := "/health-check"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	c := router.NewContext(req, rec)
	c.SetPath(path)
	if err := rest.HealthCheck(c); err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}

}
