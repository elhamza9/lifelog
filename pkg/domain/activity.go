package domain

import (
	"errors"
	"fmt"
	"time"
)

// ActivityID is an object-value representing ID of an activity
type ActivityID uint

// Activity Entity
type Activity struct {
	ID       ActivityID
	Label    string
	Place    string
	Desc     string
	Time     time.Time
	Duration time.Duration
	Tags     []Tag
}

// Constants
const (
	ActivityLabelMinLen int = 5
	ActivityLabelMaxLen int = 30
	ActivityPlaceMaxLen int = 30
	ActivityDescMaxLen  int = 255
)

// Errors
var (
	ErrActivityLabelLength error = fmt.Errorf("Activity Label must be %d ~ %d long", ActivityLabelMinLen, ActivityLabelMaxLen)
	ErrActivityPlaceLength error = fmt.Errorf("Activity Place must be maximum %d long", ActivityPlaceMaxLen)
	ErrActivityDescLength  error = fmt.Errorf("Activity Description must be maximum %d long", ActivityDescMaxLen)
	ErrActivityTimeFuture  error = errors.New("Activity Time + Duration can not result in future date")
	ErrActivityNotFound    error = errors.New("Activity Not Found")
)
