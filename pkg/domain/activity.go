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
)

// ************* Methods *************

// String returns a one line string representation of an activity
func (act Activity) String() string {
	return fmt.Sprintf("[%d | %s | %s ]", act.ID, act.Label, act.Time.Format("2006-01-02 15:04"))
}

// Valid checks primitive, non-db-related fields for validity
func (act Activity) Valid() error {
	now := time.Now()
	// Check Label Length
	if len(act.Label) < ActivityLabelMinLen || len(act.Label) > ActivityLabelMaxLen {
		return ErrActivityLabelLength
	}
	// Check Place Length
	if len(act.Place) > ActivityPlaceMaxLen {
		return ErrActivityPlaceLength
	}
	// Transform Place to Lowercase
	// Check Desc Length
	if len(act.Desc) > ActivityDescMaxLen {
		return ErrActivityDescLength
	}
	// Check TimeEnd not future
	if timeEnd := act.Time.Add(act.Duration); timeEnd.After(now) {
		return ErrActivityTimeFuture
	}
	// Everything is good
	return nil
}
