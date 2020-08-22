package domain

import (
	"errors"
	"fmt"
)

// TagID is a value object representing Id of a Tag.
type TagID uint

// Tag Entity
type Tag struct {
	ID   TagID
	Name string
}

// Constants
const (
	TagNameMinLength       int    = 3
	TagNameMaxLength       int    = 20
	TagNameValidCharacters string = `^[\w-]*$` // Only Alphanumeric characters and dashes
)

// Errors
var (
	ErrTagNameTooShort          = fmt.Errorf("Tag name must be at least %d long", TagNameMinLength)
	ErrTagNameTooLong           = fmt.Errorf("Tag name must be at most %d long", TagNameMinLength)
	ErrTagNameInvalidCharacters = errors.New("Tag name can only contain alphanumeric characters and dashes")
)
