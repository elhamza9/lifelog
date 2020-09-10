package db

import (
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
)

// Tag Model
type Tag struct {
	ID        domain.TagID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName specifies the name of the table for the tag model
func (t Tag) TableName() string { return "tags" }
