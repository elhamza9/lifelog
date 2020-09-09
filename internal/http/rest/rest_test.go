package rest_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/elhamza90/lifelog/internal/http/rest"
	"github.com/elhamza90/lifelog/internal/store/memory"
	"github.com/elhamza90/lifelog/internal/usecase/adding"
	"github.com/elhamza90/lifelog/internal/usecase/auth"
	"github.com/elhamza90/lifelog/internal/usecase/deleting"
	"github.com/elhamza90/lifelog/internal/usecase/editing"
	"github.com/elhamza90/lifelog/internal/usecase/listing"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var (
	router *echo.Echo
	hnd    *rest.Handler
	repo   memory.Repository
)

// hashEnvVarName specifies the name of the environment variable
// where the testing password hash should be stored
const hashEnvVarName string = "LFLG_TEST_PASS_HASH"

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	log.Debug("Setting Up Test router")

	router = echo.New()

	repo = memory.NewRepository()
	lister := listing.NewService(&repo)
	adder := adding.NewService(&repo)
	editor := editing.NewService(&repo)
	deletor := deleting.NewService(&repo)
	authenticator := auth.NewService(hashEnvVarName)
	hnd = rest.NewHandler(&lister, &adder, &editor, &deletor, &authenticator)

	os.Setenv("LFLG_JWT_SECRET", "testsecret")
	if err := rest.RegisterRoutes(router, hnd); err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
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
