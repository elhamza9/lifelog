package domain

import (
	"errors"
	"fmt"
)

// TagID is a value-object representing Id of a Tag.
type TagID uint

// Tag Entity
type Tag struct {
	ID   TagID
	Name string
}

// Constants for tag name conditions
const (
	TagNameMinLength       int    = 3
	TagNameMaxLength       int    = 20
	TagNameValidCharacters string = `^[\w-]*$` // Only Alphanumeric characters and dashes
)

// Errors
var (
	ErrTagNameLen               = fmt.Errorf("Tag name must be %d ~ %d characters long", TagNameMinLength, TagNameMaxLength)
	ErrTagNameInvalidCharacters = errors.New("Tag name can only contain alphanumeric characters and dashes")
	ErrTagNameDuplicate         = errors.New("Tag name duplicate")

	ErrTagNotFound = errors.New("Tag with provided ID not found")
)
