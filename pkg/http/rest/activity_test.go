package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
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

func TestActivityDetails(t *testing.T) {
	// Init Repo with one test activity
	act := domain.Activity{
		ID:       1,
		Label:    "Test Activity",
		Time:     time.Now().AddDate(0, 0, -2),
		Duration: time.Duration(time.Hour),
		Tags: []domain.Tag{
			{ID: 1, Name: "tag1"},
			{ID: 2, Name: "tag2"},
		},
	}
	repo.Activities = map[domain.ActivityID]domain.Activity{
		act.ID: act,
	}

	// Sub-tests definitions
	tests := map[string]struct {
		idStr        string
		expectedCode int
	}{
		"Existing Activity": {
			idStr:        strconv.Itoa(int(act.ID)),
			expectedCode: http.StatusOK,
		},
		"Non-Existing Activity": {
			idStr:        "234234", // Random non-existing ID
			expectedCode: http.StatusNotFound,
		},
		"Wrong Id format": {
			idStr:        "blabls",
			expectedCode: http.StatusBadRequest,
		},
	}

	// Sub-tests Execution
	const path string = "/activities/:id"
	const url string = "/activities/%s"
	var (
		req *http.Request
		rec *httptest.ResponseRecorder
		ctx echo.Context
	)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req = httptest.NewRequest(http.MethodGet, fmt.Sprintf(url, test.idStr), nil)
			rec = httptest.NewRecorder()
			ctx = router.NewContext(req, rec)
			ctx.SetPath(path)
			ctx.SetParamNames("id")
			ctx.SetParamValues(test.idStr)
			hnd.ActivityDetails(ctx)
			if rec.Code != test.expectedCode {
				body := rec.Body.String()
				t.Fatalf("\nExpected Code: %d\nReturned Code: %d\nReturned Body: %s", test.expectedCode, rec.Code, body)
			}
		})
	}
}

func TestDeleteActivity(t *testing.T) {

	// Init Repo with two test activities: one with and one without expense
	actWithoutExpense := domain.Activity{
		ID:       1,
		Label:    "Test Activity",
		Time:     time.Now().AddDate(0, 0, -2),
		Duration: time.Duration(time.Hour),
		Tags: []domain.Tag{
			{ID: 1, Name: "tag1"},
			{ID: 2, Name: "tag2"},
		},
	}
	actWithExpense := domain.Activity{
		ID:       2,
		Label:    "Test Activity",
		Time:     time.Now().AddDate(0, 0, -3),
		Duration: time.Duration(time.Hour),
		Tags: []domain.Tag{
			{ID: 2, Name: "tag2"},
		},
	}
	repo.Activities = map[domain.ActivityID]domain.Activity{
		actWithoutExpense.ID: actWithoutExpense,
		actWithExpense.ID:    actWithExpense,
	}
	repo.Expenses = map[domain.ExpenseID]domain.Expense{
		1: {
			ID:         1,
			Label:      "Expense for activity 2",
			Value:      14,
			Unit:       "Eu",
			Time:       time.Now().AddDate(0, 0, -3),
			ActivityID: actWithExpense.ID,
		},
	}

	// Sub-tests definitions
	tests := map[string]struct {
		idStr        string
		expectedCode int
	}{
		"Existing Activity Without Expense": {
			idStr:        strconv.Itoa(int(actWithoutExpense.ID)),
			expectedCode: http.StatusNoContent,
		},
		"Existing Activity With expense": {
			idStr:        strconv.Itoa(int(actWithExpense.ID)),
			expectedCode: http.StatusUnprocessableEntity,
		},
		"Non-Existing Activity": {
			idStr:        "234234", // Random non-existing ID
			expectedCode: http.StatusNotFound,
		},
		"Wrong Id format": {
			idStr:        "blabls",
			expectedCode: http.StatusBadRequest,
		},
	}

	// Sub-tests Execution
	const path string = "/activities/:id"
	const url string = "/activities/%s"
	var (
		req *http.Request
		rec *httptest.ResponseRecorder
		ctx echo.Context
	)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req = httptest.NewRequest(http.MethodGet, fmt.Sprintf(url, test.idStr), nil)
			rec = httptest.NewRecorder()
			ctx = router.NewContext(req, rec)
			ctx.SetPath(path)
			ctx.SetParamNames("id")
			ctx.SetParamValues(test.idStr)
			hnd.DeleteActivity(ctx)
			if rec.Code != test.expectedCode {
				body := rec.Body.String()
				t.Fatalf("\nExpected Code: %d\nReturned Code: %d\nReturned Body: %s", test.expectedCode, rec.Code, body)
			}
		})
	}
}
