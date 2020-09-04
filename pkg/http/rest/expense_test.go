package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/labstack/echo/v4"
)

func TestExpensesByDate(t *testing.T) {
	const path string = "/expenses%s"
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
			hnd.ExpensesByDate(ctx)
			if rec.Code != test.expectedCode {
				body := rec.Body.String()
				t.Fatalf("\nExpected Code: %d\nReturned Code: %v\nReturned Body: %s", test.expectedCode, rec.Code, body)
			}
		})
	}
}

func TestAddExpense(t *testing.T) {
	// Init Repo with some tags
	repo.Tags = map[domain.TagID]domain.Tag{
		1: {ID: 1, Name: "tag1"},
		2: {ID: 2, Name: "tag2"},
		3: {ID: 3, Name: "tag3"},
	}
	// Sub-tests definitions
	tests := map[string]struct {
		json         string
		expectedCode int
	}{
		"Correct": {
			json:         `{"label":"New Expense","value":9.5,"unit":"eu","time":"2020-04-01T18:00:00Z","tagIds":[1,3]}`,
			expectedCode: http.StatusCreated,
		},
		"Non-Existing Tag": {
			json:         `{"label":"New Expense","value":9.5,"unit":"eu","time":"2020-04-01T18:00:00Z","tagIds":[1,33]}`,
			expectedCode: http.StatusUnprocessableEntity,
		},
		"Time Future": {
			json:         `{"label":"New Expense","value":9.5,"unit":"eu",` + fmt.Sprintf("\"time\":\"%s\"", time.Now().AddDate(0, 0, 1).Format("2006-01-02")) + `,"tagIds":[1,3]}`,
			expectedCode: http.StatusBadRequest,
		},
		"Value zero": {
			json:         `{"label":"New Expense","value":0,"unit":"eu","time":"2020-04-01T18:00:00Z","tagIds":[1,3]}`,
			expectedCode: http.StatusBadRequest,
		},
	}
	// Sub-tests Execution
	const path string = "/expenses"
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
			hnd.AddExpense(ctx)
			if rec.Code != test.expectedCode {
				body := rec.Body.String()
				t.Fatalf("\nExpected Code: %d\nReturned Code: %d\nReturned Body: %s", test.expectedCode, rec.Code, body)
			}
		})
	}
}
