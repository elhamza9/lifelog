package adding_test

import (
	"log"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

func TestNewActivity(t *testing.T) {
	t.Fatal("Not yet implemented")
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
		/*
			"Correct": {
					label: "New Activity",
					place: "Beach",
					desc:  "Play Soccer",
					time:  time.Now().AddDate(0, 0, -1),
					dur:   time.Duration(time.Hour),
					tags:  &[]domain.Tag{{ID: 100002}, {ID: 100005}},
				},
				"Short Label": {
					label: "",
					place: "Beach",
					desc:  "Play Soccer",
					time:  time.Now().AddDate(0, 0, -1),
					dur:   time.Duration(time.Hour),
					tags:  &[]domain.Tag{{ID: 100002}, {ID: 100005}},
				},
		*/
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resAct, err := adder.NewActivity(test.label, test.place, test.desc, test.time, test.dur, &test.tags)
			testFailed := err != test.expectedErr
			if testFailed {
				t.Fatalf("\nExpecting: %v\nBut Got: %v", test.expectedErr, err)
			}
			if err == nil {
				// Tests after creation successful
				log.Println(resAct)
			}

		})
	}

}
