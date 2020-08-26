package domain

import (
	"errors"
	"fmt"
	"regexp"
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
)

// String returns a one-line representation of an tag
func (t Tag) String() string {
	return fmt.Sprintf("[%d | %s ]", t.ID, t.Name)
}

// Valid returns an error if the tag fields are invalid
func (t Tag) Valid() error {
	// Check tag name length
	nameTooShort := len(t.Name) < TagNameMinLength
	nameTooLong := len(t.Name) > TagNameMaxLength
	if nameTooShort || nameTooLong {
		return ErrTagNameLen
	}
	// Check Tag name characters
	if match, _ := regexp.Match(TagNameValidCharacters, []byte(t.Name)); !match {
		return ErrTagNameInvalidCharacters
	}
	// Everything is good
	return nil

}
