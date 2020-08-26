package adding_test

import (
	"strings"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store"
)

func TestNewActivity(t *testing.T) {
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
		place       string
		desc        string
		time        time.Time
		dur         time.Duration
		tags        []domain.Tag
		expectedErr error
	}{
		"Correct": {
			label:       "New Activity",
			place:       "Beach",
			desc:        "Play Soccer",
			time:        time.Now().AddDate(0, 0, -1),
			dur:         time.Duration(time.Hour),
			tags:        []domain.Tag{{ID: 100002}, {ID: 100005}},
			expectedErr: nil,
		},
		"Short Label": {
			label:       "",
			place:       "Beach",
			desc:        "Play Soccer",
			time:        time.Now().AddDate(0, 0, -1),
			dur:         time.Duration(time.Hour),
			tags:        []domain.Tag{{ID: 100002}, {ID: 100005}},
			expectedErr: domain.ErrActivityLabelLength,
		},
		"Long Label": {
			label:       "My very very very very very very very very very very very very very very very Long Label",
			place:       "Beach",
			desc:        "Play Soccer",
			time:        time.Now().AddDate(0, 0, -1),
			dur:         time.Duration(time.Hour),
			tags:        []domain.Tag{{ID: 100002}, {ID: 100005}},
			expectedErr: domain.ErrActivityLabelLength,
		},
		"Long Place": {
			label:       "New Activity",
			place:       "My  very very very very very very very very very very very Long Place",
			desc:        "Play Soccer",
			time:        time.Now().AddDate(0, 0, -1),
			dur:         time.Duration(time.Hour),
			tags:        []domain.Tag{{ID: 100002}, {ID: 100005}},
			expectedErr: domain.ErrActivityPlaceLength,
		},
		"Time + Dur Future": {
			label:       "New Activity",
			place:       "Beach",
			desc:        "Play Soccer",
			time:        time.Now(),
			dur:         time.Duration(time.Hour),
			tags:        []domain.Tag{{ID: 100002}, {ID: 100005}},
			expectedErr: domain.ErrActivityTimeFuture,
		},
		"Non Existing Tag": {
			label:       "New Activity",
			place:       "Beach",
			desc:        "Play Soccer",
			time:        time.Now().AddDate(0, 0, -1),
			dur:         time.Duration(time.Hour),
			tags:        []domain.Tag{{ID: 100002}, {ID: 100010}},
			expectedErr: store.ErrTagNotFound,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resAct, err := adder.NewActivity(test.label, test.place, test.desc, test.time, test.dur, &test.tags)
			testFailed := err != test.expectedErr
			if testFailed {
				t.Fatalf("\nExpected Err: %v\nReturned Err: %v", test.expectedErr, err)
			}
			if err == nil {
				// Tests after creation successful
				expectedPlace := strings.ToLower(test.place)
				if resAct.Place != expectedPlace {
					t.Fatalf("Expected Place: %s\nReturned Place: %s", expectedPlace, resAct.Place)
				}
			}
		})
	}

}
