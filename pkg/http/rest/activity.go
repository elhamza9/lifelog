package rest

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store"
	"github.com/labstack/echo/v4"
)

// ActivitiesByDate handler returns a list of all activities done
// from a specific date up to now.
// It requires a query parameter "from" specifying the date as mm-dd-yyyy
// If no parameter is found default is to return activities of last 3 months
func (h *Handler) ActivitiesByDate(c echo.Context) error {
	dateStr := c.QueryParam("from")
	var date time.Time
	if len(dateStr) == 0 {
		date = time.Now().AddDate(0, -3, 0)
	} else {
		var err error
		date, err = time.Parse("01-02-2006", dateStr)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	}
	activities, err := h.lister.ActivitiesByTime(date)
	if err != nil {
		if errors.Is(err, domain.ErrActivityTimeFuture) {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, activities)
}

// ActivityDetails handler returns details of activity with given ID
// It required a path parameter :id
func (h *Handler) ActivityDetails(c echo.Context) error {
	// Get ID from Path param
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Get Activity
	act, err := h.lister.Activity(domain.ActivityID(id))
	if err != nil {
		if errors.Is(err, store.ErrActivityNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, act)
}
