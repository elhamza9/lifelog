package memory

import (
	"math/rand"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store"
)

func generateRandomActivityID() domain.ActivityID {
	rand.Seed(time.Now().UnixNano())
	res := rand.Intn(10000)
	return domain.ActivityID(res)
}

// FindActivityByID returns activity with given ID.
// If none is found, returns error
func (repo Repository) FindActivityByID(id domain.ActivityID) (domain.Activity, error) {
	for _, act := range *repo.Activities {
		if act.ID == id {
			return act, nil
		}
	}
	return domain.Activity{}, store.ErrActivityNotFound
}

// AddActivity stores the given activity in memory and returns created activity
func (repo Repository) AddActivity(act domain.Activity) (domain.ActivityID, error) {
	act.ID = generateRandomActivityID()
	(*repo.Activities)[act.ID] = act
	return act.ID, nil
}

// FindActivitiesByTime returns activities
// with Time field greater than or equal to the given time
func (repo Repository) FindActivitiesByTime(t time.Time) (*[]domain.Activity, error) {
	res := []domain.Activity{}
	for _, act := range *repo.Activities {
		if !act.Time.Before(t) {
			res = append(res, act)
		}
	}
	return &res, nil
}

// FindActivitiesByTag returns actenses that have the provided tag in their Tags field
func (repo Repository) FindActivitiesByTag(tid domain.TagID) (*[]domain.Activity, error) {
	res := []domain.Activity{}
	for _, act := range *repo.Activities {
		for _, tag := range act.Tags {
			if tag.ID == tid {
				res = append(res, act)
			}
		}
	}
	return &res, nil
}

// DeleteActivity removes activity with provided ID from memory
func (repo Repository) DeleteActivity(id domain.ActivityID) error {
	delete((*repo.Activities), id)
	return nil
}

// EditActivity edits given activity in memory
func (repo Repository) EditActivity(act domain.Activity) (domain.Activity, error) {

	(*repo.Activities)[act.ID] = act
	return (*repo.Activities)[act.ID], nil
}
