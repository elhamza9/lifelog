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
	ctx := router.NewContext(req, rec)
	ctx.SetPath(path)
	hnd.GetAllTags(ctx)
	if rec.Code != expectedCode {
		t.Fatalf("\nExpected Code: %d\nReturned Code: %d", expectedCode, rec.Code)
	}
}
