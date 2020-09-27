package client_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/http/rest/client"
	"github.com/elhamza90/lifelog/internal/http/rest/server"
)

func TestPostTag(t *testing.T) {
	_, err := client.PostTag(server.JSONReqTag{Name: fmt.Sprintf("my-test-tag-%d", rand.Intn(999))}, token)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchTags(t *testing.T) {
	_, err := client.FetchTags(token)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateTag(t *testing.T) {
	// Create Test Tag
	payload := server.JSONReqTag{
		Name: fmt.Sprintf("my-test-tag-%d", rand.Intn(999)),
	}
	id, err := client.PostTag(payload, token)
	if err != nil {
		t.Fatal(err)
	}
	// Test Delete
	payload = server.JSONReqTag{
		ID:   id,
		Name: fmt.Sprintf("my-edited-tag-%d", rand.Intn(999)),
	}
	if err = client.UpdateTag(payload, token); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTag(t *testing.T) {
	// Create Test Tag
	id, err := client.PostTag(server.JSONReqTag{Name: fmt.Sprintf("my-test-tag-%d", rand.Intn(999))}, token)
	if err != nil {
		t.Fatal(err)
	}
	// Test Delete
	if err = client.DeleteTag(domain.TagID(id), token); err != nil {
		t.Fatal(err)
	}
}
