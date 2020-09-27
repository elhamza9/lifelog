package client_test

import (
	"testing"

	"github.com/elhamza90/lifelog/internal/http/rest/client"
)

func TestLogin(t *testing.T) {
	_, _, err := client.Login("pass_pass")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRefreshToken(t *testing.T) {
	refresh := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.CxJEg6krZCATx1I0ttD9LdWNZPJiBg_7vF9lAFsufl0"
	_, err := client.RefreshToken(refresh)
	if err != nil {
		t.Fatal(err)
	}
}
