package client_test

import (
	"testing"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/http/rest/client"
	"github.com/elhamza90/lifelog/internal/http/rest/server"
)

func TestPostActivity(t *testing.T) {
	payload := server.JSONReqActivity{
		Label:    "Do smth",
		Place:    "Somewhere",
		Desc:     "With Details",
		Time:     time.Now().Add(time.Duration(-1 * time.Hour)),
		Duration: time.Duration(time.Minute * 20),
		TagIds:   []domain.TagID{},
	}
	_, err := client.PostActivity(payload, token)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchActivities(t *testing.T) {
	_, err := client.FetchActivities(token, time.Now().AddDate(0, -3, 0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteActivity(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"
	// Create Test Activity
	payload := server.JSONReqActivity{
		Label:    "Do smth",
		Place:    "Somewhere",
		Desc:     "With Details",
		Time:     time.Now().Add(time.Duration(-1 * time.Hour)),
		Duration: time.Duration(time.Minute * 20),
	}
	id, err := client.PostActivity(payload, token)
	if err != nil {
		t.Fatal(err)
	}
	// Test Delete
	if err = client.DeleteActivity(domain.ActivityID(id), token); err != nil {
		t.Fatal(err)
	}
}

func TestFetchActivityDetails(t *testing.T) {
	// Create Test Activity
	payload := server.JSONReqActivity{
		Label:    "Do smth",
		Place:    "Somewhere",
		Desc:     "With Details",
		Time:     time.Now().Add(time.Duration(-1 * time.Hour)),
		Duration: time.Duration(time.Minute * 20),
		TagIds:   []domain.TagID{},
	}
	id, err := client.PostActivity(payload, token)
	if err != nil {
		t.Fatal(err)
	}
	// Test Delete
	if _, err = client.FetchActivityDetails(domain.ActivityID(id), token); err != nil {
		t.Fatal(err)
	}
}
