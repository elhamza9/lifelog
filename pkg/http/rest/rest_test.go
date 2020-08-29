package rest_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/elhamza90/lifelog/pkg/http/rest"
	"github.com/elhamza90/lifelog/pkg/store/memory"
	"github.com/elhamza90/lifelog/pkg/usecase/adding"
	"github.com/elhamza90/lifelog/pkg/usecase/deleting"
	"github.com/elhamza90/lifelog/pkg/usecase/editing"
	"github.com/elhamza90/lifelog/pkg/usecase/listing"
	"github.com/labstack/echo/v4"
)

var router *echo.Echo
var hnd *rest.Handler
var repo memory.Repository

func TestMain(m *testing.M) {
	log.Println("Setting Up Main")
	router = echo.New()
	repo = memory.NewRepository()
	lister := listing.NewService(&repo)
	adder := adding.NewService(&repo)
	editor := editing.NewService(&repo)
	deletor := deleting.NewService(&repo)
	hnd = rest.NewHandler(&lister, &adder, &editor, &deletor)

	rest.RegisterRoutes(router, hnd)

	os.Exit(m.Run())
}

func TestHealthCheck(t *testing.T) {
	log.Print("Test Health Check")
	path := "/health-check"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	c := router.NewContext(req, rec)
	c.SetPath(path)
	if err := rest.HealthCheck(c); err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
}
