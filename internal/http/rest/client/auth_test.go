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
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.G2f51qduWeW9XHILKdP-ZROENssiPEUG_I_pOB-3298"
	_, err := client.RefreshToken(token)
	if err != nil {
		t.Fatal(err)
	}
}
