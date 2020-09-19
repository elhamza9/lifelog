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
