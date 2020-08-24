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
