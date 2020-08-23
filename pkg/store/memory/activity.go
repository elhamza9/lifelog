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
