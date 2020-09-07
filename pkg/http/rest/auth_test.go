package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestLogin(t *testing.T) {
	const testPass string = "test_pass"
	os.Setenv("LFLG_PASSWORD", testPass)
	// Subtests Definition
	tests := map[string]struct {
		json         string
		expectedCode int
	}{
		"Correct Credentials": {
			json:         fmt.Sprintf("{\"password\":\"%s\"}", testPass),
			expectedCode: http.StatusOK,
		},
		"Incorrect Credentials": {
			json:         `{"password": "mywrongtestpass"}`,
			expectedCode: http.StatusUnauthorized,
		},
		"Short Password": {
			json:         `{"password": "pswd"}`,
			expectedCode: http.StatusBadRequest,
		},
	}
	// Subtests Execution
	const path string = "/auth/login"
	var (
		req *http.Request
		rec *httptest.ResponseRecorder
		ctx echo.Context
	)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req = httptest.NewRequest(http.MethodPost, path, strings.NewReader(test.json))
			req.Header.Set("Content-type", "application/json")
			rec = httptest.NewRecorder()
			ctx = router.NewContext(req, rec)
			ctx.SetPath(path)
			hnd.Login(ctx)
			body := rec.Body.String()
			if rec.Code != test.expectedCode {
				t.Fatalf("\nExpected Code: %d\nReturned Code: %d\nReturned Body: %s", test.expectedCode, rec.Code, body)
			}
			if rec.Code == http.StatusOK {
				pat := `^{"at":".*"}`
				if match, err := regexp.Match(pat, []byte(body)); !match {
					t.Fatal(body, err)
				}
			}
		})
	}
}
