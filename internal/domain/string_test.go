package domain

// Tests just to visually check the string format of the entities

import (
	"log"
	"testing"
	"time"
)

func TestActivityString(t *testing.T) {
	act := Activity{
		ID:       1,
		Label:    "Act 1",
		Place:    "Some Place",
		Desc:     "Some Details",
		Time:     time.Now().AddDate(0, 0, -1),
		Duration: time.Duration(time.Minute * 45),
	}
	log.Print(act)
}

func TestExpenseString(t *testing.T) {
	exp := Expense{
		ID:    1,
		Label: "Test Expense",
		Value: 10,
		Unit:  "Eu",
		Time:  time.Now().AddDate(0, 0, -1),
	}
	log.Print(exp)
}

func TestTagString(t *testing.T) {
	tag := Tag{ID: 1, Name: "tag-1"}
	log.Print(tag)
}
