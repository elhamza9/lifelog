package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
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
