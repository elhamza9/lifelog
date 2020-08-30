package rest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/labstack/echo/v4"
)

func TestGetAllTags(t *testing.T) {
	expectedCode := http.StatusOK
	path := "/tags"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	ctx := router.NewContext(req, rec)
	ctx.SetPath(path)
	hnd.GetAllTags(ctx)
	if rec.Code != expectedCode {
		t.Fatalf("\nExpected Code: %d\nReturned Code: %d", expectedCode, rec.Code)
	}
}

func TestAddTag(t *testing.T) {
	// Init repo with a tag to test duplicate tag name return code.
	repo.Tags = map[domain.TagID]domain.Tag{
		8987: {ID: 8987, Name: "existing-tag"},
	}

	// Sub-tests definition
	tests := map[string]struct {
		json         string
		expectedCode int
	}{
		"Correct": {
			json:         `{"name":"new-tag"}`,
			expectedCode: http.StatusCreated,
		},
		"Duplicate": {
			json:         `{"name":"existing-tag"}`,
			expectedCode: http.StatusBadRequest,
		},
		"Invalid Chars": {
			json:         `{"name":"bad$tag"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	// Sub-tests execution
	const path string = "/tags"
	var (
		req *http.Request
		rec *httptest.ResponseRecorder
		ctx echo.Context
	)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req = httptest.NewRequest(http.MethodPost, path, strings.NewReader(test.json))
			req.Header.Set("Content-Type", "application/json")
			rec = httptest.NewRecorder()
			ctx = router.NewContext(req, rec)
			ctx.SetPath(path)
			hnd.AddTag(ctx)
			if rec.Code != test.expectedCode {
				t.Fatalf("\nExpected Code: %d\nReturned Code: %v", test.expectedCode, rec.Code)
			}
		})
	}
}
