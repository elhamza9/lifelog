package client_test

import (
	"testing"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/http/rest/client"
)

func TestPostTag(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"
	_, err := client.PostTag(domain.Tag{Name: "my-test-tag-1"}, token)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchTags(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"
	_, err := client.FetchTags(token)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTag(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"
	// Create Test Tag
	id, err := client.PostTag(domain.Tag{Name: "my-test-tag-3123"}, token)
	if err != nil {
		t.Fatal(err)
	}
	// Test Delete
	if err = client.DeleteTag(domain.TagID(id), token); err != nil {
		t.Fatal(err)
	}
}
