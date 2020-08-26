package usecase

import "errors"

// ErrTagNameDuplicate is returned when trying to add/edit a tag with a name that already exists in store
var ErrTagNameDuplicate error = errors.New("Tag name duplicate")
