package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func TestActivitiesByDate(t *testing.T) {
	const path string = "/activities%s"
	const param string = "from"
	const frmt string = "01-02-2006" // time format in query param
	now := time.Now()

	// Sub-tests definition
	tests := map[string]struct {
		filter       string
		expectedCode int
	}{
		"No Date Filter": {
			filter:       "",
			expectedCode: http.StatusOK,
		},
		"Last 2 days": {
			filter:       fmt.Sprintf("?%s=%s", param, now.AddDate(0, 0, -2).Format(frmt)),
			expectedCode: http.StatusOK,
		},
		"Date Future": {
			filter:       fmt.Sprintf("?%s=%s", param, now.AddDate(0, 0, 1).Format(frmt)),
			expectedCode: http.StatusBadRequest,
		},
		"Wrong Date Format": {
			filter:       fmt.Sprintf("?%s=2020-01-31", param),
			expectedCode: http.StatusBadRequest,
		},
	}

	// Sub-tests Execution
	var (
		req *http.Request
		rec *httptest.ResponseRecorder
		ctx echo.Context
	)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req = httptest.NewRequest(http.MethodGet, fmt.Sprintf(path, test.filter), nil)
			rec = httptest.NewRecorder()
			ctx = router.NewContext(req, rec)
			ctx.SetPath(path)
			hnd.ActivitiesByDate(ctx)
			if rec.Code != test.expectedCode {
				body := rec.Body.String()
				t.Fatalf("\nExpected Code: %d\nReturned Code: %v\nReturned Body: %s", test.expectedCode, rec.Code, body)
			}
		})
	}

}
