package listing_test

import (
	"os"
	"testing"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store/memory"
	"github.com/elhamza90/lifelog/pkg/usecase/listing"
)

var service listing.Service
var repo memory.Repository

func TestMain(m *testing.M) {
	repo = memory.NewRepository()       // Work with In-Memory DB
	service = listing.NewService(&repo) // Passing by reference to change db when testing
	os.Exit(m.Run())
}

func TestAllTags(t *testing.T) {
	// Init Repo with some tags
	repo.Tags = &map[domain.TagID]domain.Tag{
		100000: {ID: 100000, Name: "tag-100000"},
		100001: {ID: 100001, Name: "tag-100001"},
		100002: {ID: 100002, Name: "tag-100002"},
		100003: {ID: 100003, Name: "tag-100003"},
		100004: {ID: 100004, Name: "tag-100004"},
		100005: {ID: 100005, Name: "tag-100005"},
	}
	resTags, _ := service.AllTags()
	if len(*resTags) != len(*repo.Tags) {
		t.Fatalf("\nExpecting tags: %v\nBut Got: %v", repo.Tags, resTags)
	}
}
