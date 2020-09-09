package rest

import (
	"errors"
	"net/http"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
	"github.com/elhamza90/lifelog/internal/usecase/auth"
	"github.com/elhamza90/lifelog/internal/usecase/deleting"
	"github.com/labstack/echo/v4"
)

// dateFilterFormat specifies the format of date in a query filter
const dateFilterFormat string = "01-02-2006"

// errInvalidJSON represents an error that occured while binding
// (unmarshaling) a json request to struct type
var errInvalidJSON error = errors.New("Invalid JSON")

// httpErrorMsg extracts error message from echo.HTTPError struct
func httpErrorMsg(err error) string {
	he, _ := err.(*echo.HTTPError)
	return he.Message.(string)
}

// errToHTTPCode returns the http code that should be sent for an error
// The grp parameter specifies which handler group called the function
// because some errors will get treated differently depending on the handler
// in which they were received.
// Ex: a Tag not found will raise a StatusNotFound in a tag handler
//     but will raise s StatusUnprocessableEntity in an expense handler
func errToHTTPCode(err error, grp string) int {
	switch err {
	// rest errors
	case errInvalidJSON:
		return http.StatusBadRequest
	case errSigningJwt:
		return http.StatusInternalServerError
	// auth errors
	case auth.ErrPasswordLength:
		return http.StatusBadRequest
	case auth.ErrIncorrectCredentials:
		return http.StatusUnauthorized
	case auth.ErrHashNotFound:
		return http.StatusInternalServerError
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
		fallthrough
	case domain.ErrExpenseLabelLength:
		fallthrough
	case domain.ErrExpenseValue:
		fallthrough
	case domain.ErrExpenseUnitLength:
		fallthrough
	case domain.ErrExpenseTimeFuture:
		return http.StatusBadRequest
	// usecase errors
	case deleting.ErrTagHasExpenses:
		fallthrough
	case deleting.ErrTagHasActivities:
		fallthrough
	case deleting.ErrActivityHasExpenses:
		return http.StatusUnprocessableEntity
	// store errors
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
	case store.ErrExpenseNotFound:
		if grp == "expenses" {
			return http.StatusNotFound
		}
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
