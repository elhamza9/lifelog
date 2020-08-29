package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllTags(t *testing.T) {
	expectedCode := http.StatusOK
	path := "/tags"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	c := router.NewContext(req, rec)
	c.SetPath(path)
	var err error
	if err = hnd.GetAllTags(c); err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
	code := rec.Code
	if code != expectedCode {
		t.Fatalf("\nExpected Code: %d\nReturned Code: %d", expectedCode, code)
	}
}
