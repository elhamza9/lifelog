package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	const path string = "/auth/login"
	// Sub-test No original password is set
	t.Run("Password not found in system", func(t *testing.T) {
		os.Setenv(hashEnvVarName, "")
		json := `{"password":"somepassword"}`
		expectedCode := http.StatusInternalServerError
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(json))
		req.Header.Set("Content-type", "application/json")
		rec := httptest.NewRecorder()
		ctx := router.NewContext(req, rec)
		ctx.SetPath(path)
		hnd.Login(ctx)
		if rec.Code != expectedCode {
			body := rec.Body.String()
			t.Fatalf("\nExpected Code: %d\nReturned Code: %d\nReturned Body: %s", expectedCode, rec.Code, body)
		}
	})
	// Generate Password hash and save it in env variable
	const testPass string = "test_pass"
	hash, err := bcrypt.GenerateFromPassword([]byte(testPass), 10)
	if err != nil {
		t.Fatalf("Error generating bcrypt hash: %s", err)
	}
	os.Setenv(hashEnvVarName, string(hash))
	defer os.Setenv(hashEnvVarName, "")
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
		"Long Password": {
			json:         fmt.Sprintf("{\"password\": \"%s\"}", strings.Repeat("abc", 100)),
			expectedCode: http.StatusBadRequest,
		},
		"Invalid JSON": {
			json:         `{"password":}`,
			expectedCode: http.StatusBadRequest,
		},
	}
	// Subtests Execution
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
				pat := `^{"at":".*","rt":".*"}`
				if match, err := regexp.Match(pat, []byte(body)); !match {
					t.Fatal(body, err)
				}
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	correctSecret := os.Getenv("LFLG_JWT_REFRESH_SECRET")
	// Create a new token
	genToken := func(secret string, exp time.Time) string {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["exp"] = exp.Unix()
		signed, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Log(err)
			return ""
		}
		return signed
	}
	const path string = "/auth/refresh"
	// Subtests definition
	now := time.Now()
	tests := map[string]struct {
		json         string
		expectedCode int
	}{
		"Correct": {
			json:         fmt.Sprintf("{\"refresh\":\"%s\"}", genToken(correctSecret, now.Add(time.Duration(time.Hour)))),
			expectedCode: http.StatusOK,
		},
		"Expired Token": {
			json:         fmt.Sprintf("{\"refresh\":\"%s\"}", genToken(correctSecret, now.Add(time.Duration(-1*time.Hour)))),
			expectedCode: http.StatusUnprocessableEntity,
		},
		"Token Signed with wrong secret": {
			json:         fmt.Sprintf("{\"refresh\":\"%s\"}", genToken(os.Getenv("LFLG_JWT_ACCESS_SECRET"), now.Add(time.Duration(time.Hour)))),
			expectedCode: http.StatusUnprocessableEntity,
		},
		"No Token in JSON": {
			json:         `{"what":"is this?"}`,
			expectedCode: http.StatusBadRequest,
		},
	}
	// Subtests Execution
	var (
		req *http.Request
		rec *httptest.ResponseRecorder
		ctx echo.Context
	)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Log(test.json)
			req = httptest.NewRequest(http.MethodPost, path, strings.NewReader(test.json))
			req.Header.Set("Content-type", "application/json")
			rec = httptest.NewRecorder()
			ctx = router.NewContext(req, rec)
			ctx.SetPath(path)
			hnd.RefreshToken(ctx)
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
