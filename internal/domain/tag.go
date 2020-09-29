package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TagID is a value-object representing Id of a Tag.
type TagID uint

// String returns a string representation of the id
func (id TagID) String() string {
	return strconv.Itoa(int(id))
}

// Tag Entity
type Tag struct {
	ID   TagID
	Name string
}

// Constants for tag name conditions
const (
	TagNameMinLength  int    = 3
	TagNameMaxLength  int    = 20
	TagNameValidChars string = `^[\w-]*$` // Only Alphanumeric characters and dashes
)

// Errors
var (
	ErrTagNameLen               = fmt.Errorf("Tag name must be %d ~ %d characters long", TagNameMinLength, TagNameMaxLength)
	ErrTagNameInvalidCharacters = errors.New("Tag name can only contain alphanumeric characters and dashes")
	ErrTagNameDuplicate         = errors.New("Tag name duplicate")
)

// ************* Methods *************

// String returns a one-line representation of a tag
func (t Tag) String() string {
	return fmt.Sprintf("[%d | %s ]", t.ID, t.Name)
}

// Validate checks primitive, non-db-related fields for validity
func (t *Tag) Validate() error {
	// Check tag name length
	nameTooShort := len(t.Name) < TagNameMinLength
	nameTooLong := len(t.Name) > TagNameMaxLength
	if nameTooShort || nameTooLong {
		return ErrTagNameLen
	}
	t.Name = strings.ToLower(t.Name)
	// Check Tag name characters
	if match, _ := regexp.Match(TagNameValidChars, []byte(t.Name)); !match {
		return ErrTagNameInvalidCharacters
	}
	// Everything is good
	return nil
}
