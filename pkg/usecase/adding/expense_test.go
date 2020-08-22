package adding_test

import (
	"testing"

	"github.com/elhamza90/lifelog/pkg/domain"
)

func TestNewExpense(t *testing.T) {
	// Init Repo with some tags to test checking if tags exist
	repo.Tags = &map[domain.TagID]domain.Tag{
		100000: {ID: 100000, Name: "tag-100000"},
		100001: {ID: 100001, Name: "tag-100001"},
		100002: {ID: 100002, Name: "tag-100002"},
		100003: {ID: 100003, Name: "tag-100003"},
		100004: {ID: 100004, Name: "tag-100004"},
		100005: {ID: 100005, Name: "tag-100005"},
	}

	// Sub-tests
	tests := map[string]struct {
		label       string
		val         float32
		unit        string
		tags        *[]domain.Tag
		expectedErr error
	}{}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
		})
	}
}
