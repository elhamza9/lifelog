package memory

import (
	"math/rand"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
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
	return domain.Activity{}, domain.ErrActivityNotFound
}

// AddActivity stores the given activity in memory and returns created activity
func (repo Repository) AddActivity(act domain.Activity) (domain.Activity, error) {
	act.ID = generateRandomActivityID()
	(*repo.Activities)[act.ID] = act
	return act, nil
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
