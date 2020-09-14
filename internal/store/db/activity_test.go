package db_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
	"github.com/elhamza90/lifelog/internal/store/db"
)

func TestFindActivityByID(t *testing.T) {
	defer clearDB()
	// Create test Activity
	act := db.Activity{
		ID:       546,
		Label:    "test activity",
		Place:    "Somewhere",
		Desc:     "Details",
		Time:     time.Now(),
		Duration: time.Duration(time.Hour),
		Tags:     []db.Tag{},
	}
	if err := grmDb.Create(&act).Error; err != nil {
		t.Fatalf("\nError while creating test activity:\n  %v", err)
	}
	// Tests
	tests := map[string]struct {
		id          domain.ActivityID
		expectedErr error
	}{
		"Existing Activity":     {act.ID, nil},
		"Non Existing Activity": {23423423, store.ErrActivityNotFound},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := repo.FindActivityByID(test.id); err != test.expectedErr {
				t.Fatalf("\nExpected Error: %v\nReturned Error: %v", test.expectedErr, err)
			}
		})
	}
}

func TestSaveActivity(t *testing.T) {
	// Create test Activity
	defer clearDB()
	tags := []db.Tag{{ID: 1, Name: "test-tag-1"}, {ID: 2, Name: "test-tag-2"}}
	if err := grmDb.Create(&tags).Error; err != nil {
		t.Fatalf("\nError while creating Test Tags:\n  %v", err)
	}
	act := domain.Activity{
		ID:       546,
		Label:    "test activity",
		Place:    "Somewhere",
		Desc:     "Details",
		Time:     time.Now(),
		Duration: time.Duration(time.Hour),
		Tags:     []domain.Tag{{ID: tags[0].ID}},
	}
	// Test Save
	id, err := repo.SaveActivity(act)
	if err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
	var created db.Activity
	if err := grmDb.Preload("Tags").First(&created, id).Error; err != nil {
		t.Fatalf("\nUnexpected Error while retrieving saved activity:\n  %v", err)
	}
	if len(created.Tags) != len(act.Tags) {
		t.Fatalf("\nExpected %d Tags\nReturned %d Tags", len(act.Tags), len(created.Tags))
	}
}

func TestFindActivitiesByTime(t *testing.T) {
	// Create test 100 activities:
	// one for each day starting from today going backward
	defer clearDB()
	const nbrActivities int = 100
	activities := make([]db.Activity, nbrActivities)
	now := time.Now()
	for i := 0; i < nbrActivities; i++ {
		activities[i] = db.Activity{
			Label:    fmt.Sprintf("Test Activity %d", i),
			Place:    "Somewhere",
			Desc:     "Details",
			Time:     now.AddDate(0, 0, -i),
			Duration: time.Duration(time.Hour),
			Tags:     []db.Tag{},
		}
	}
	// Shuffle activities before saving them to DB to avoid getting them by insertion order
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(activities), func(i, j int) { activities[i], activities[j] = activities[j], activities[i] })
	if err := grmDb.Create(&activities).Error; err != nil {
		t.Fatalf("\nError while creating test activities:\n  %v", err)
	}
	// Test Get Activities of last 5 days (Should be 6 activities)
	minTime := now.AddDate(0, 0, -5)
	nbrExpectedActivities := 6
	res, err := repo.FindActivitiesByTime(minTime)
	if err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
	if len(res) != nbrExpectedActivities {
		t.Log(res)
		t.Fatalf("\nExpecting %d Activities\nReturned %d activities", nbrExpectedActivities, len(res))
	}
	// Test Order by time
	for i, exp := range res {
		if i < len(res)-1 {
			if exp.Time.Before(res[i+1].Time) {
				t.Fatal("Activities not ordered by time")
			}
		}
	}
}

func TestFindActivitiesByTag(t *testing.T) {
	// Create test activities & tags:
	defer clearDB()
	var (
		tag1 db.Tag = db.Tag{ID: 11, Name: "test-tag-1"}
		tag2 db.Tag = db.Tag{ID: 12, Name: "test-tag-2"}
		tag3 db.Tag = db.Tag{ID: 13, Name: "test-tag-3"}
	)
	if err := grmDb.Create(&[]db.Tag{tag1, tag2, tag3}).Error; err != nil {
		t.Fatalf("\nError while creating test tags:\n  %v", err)
	}
	now := time.Now()
	activities := []db.Activity{
		{
			Label:    "Test Activity 1 ( Tag1, Tag3 )",
			Place:    "Somewhere",
			Desc:     "Details",
			Time:     now.AddDate(0, 0, -20),
			Duration: time.Duration(time.Hour),
			Tags:     []db.Tag{tag1, tag3},
		},
		{
			Label:    "Test Activity 2 ( Tag2, Tag3 )",
			Place:    "Somewhere",
			Desc:     "Details",
			Time:     now.AddDate(0, 0, -3),
			Duration: time.Duration(time.Hour),
			Tags:     []db.Tag{tag2, tag3},
		},
		{
			Label:    "Test Activity 3 ( Tag1, Tag2 )",
			Place:    "Somewhere",
			Desc:     "Details",
			Time:     now.AddDate(0, 0, -15),
			Duration: time.Duration(time.Hour),
			Tags:     []db.Tag{tag1, tag2},
		},
	}
	if err := grmDb.Create(&activities).Error; err != nil {
		t.Fatalf("\nError while creating test activities:\n  %v", err)
	}
	// Test Get Activities of Tag 1
	res, err := repo.FindActivitiesByTag(tag1.ID)
	if err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
	expectedActivities := [2]db.Activity{activities[2], activities[0]}
	if len(res) != len(expectedActivities) {
		t.Fatalf("\nExpecting %d Activities\nReturned %d activities", len(expectedActivities), len(res))
	}
	for i, exp := range res {
		if exp.ID != expectedActivities[i].ID {
			t.Fatalf("\nExpecting activity ID %d in %d position, Got ID %d", expectedActivities[i].ID, i+1, exp.ID)
		}
	}
}

func TestDeleteActivity(t *testing.T) {
	t.Fatal("Test not yet implemented")
}

func TestEditActivity(t *testing.T) {
	t.Fatal("Test not yet implemented")
}
