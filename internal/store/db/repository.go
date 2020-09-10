package db

import (
	"errors"

	"gorm.io/gorm"
)

var errNotImplemented error = errors.New("Repo Method not yet implemented")

// Repository manages data in a Database using GORM orm
type Repository struct {
	db *gorm.DB
}

// NewRepository returns new repository.
func NewRepository(db *gorm.DB) Repository {
	return Repository{db}
}
