package rest

import (
	"net/http"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store"
	"github.com/elhamza90/lifelog/pkg/usecase/deleting"
)

// errToHTTPCode returns the http code that should be sent for an error
// The grp parameter specifies which handler group called the function
// because some errors will get treated differently depending on the handler
// in which they were received.
// Ex: a Tag not found will raise a StatusNotFound in a tag handler
//     but will raise s StatusUnprocessableEntity in an expense handler
func errToHTTPCode(err error, grp string) int {
	switch err {
	// domain errors
	case domain.ErrTagNameDuplicate:
		fallthrough
	case domain.ErrTagNameLen:
		fallthrough
	case domain.ErrTagNameInvalidCharacters:
		return http.StatusBadRequest
	case domain.ErrActivityTimeFuture:
		fallthrough
	case domain.ErrActivityLabelLength:
		fallthrough
	case domain.ErrActivityPlaceLength:
		fallthrough
	case domain.ErrActivityDescLength:
		return http.StatusBadRequest
	// usecase errors
	case deleting.ErrTagHasExpenses:
		fallthrough
	case deleting.ErrTagHasActivities:
		fallthrough
	case deleting.ErrActivityHasExpenses:
		return http.StatusUnprocessableEntity
	// domain errors
	case store.ErrTagNotFound:
		if grp == "tags" {
			return http.StatusNotFound
		}
		return http.StatusUnprocessableEntity
	case store.ErrActivityNotFound:
		if grp == "activities" {
			return http.StatusNotFound
		}
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
