package rest

import (
	"net/http"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/labstack/echo/v4"
)

// jsonExpense is used to unmarshal a json expense
type jsonExpense struct {
	Label      string            `json:"label"`
	Time       time.Time         `json:"time"`
	Value      float32           `json:"value"`
	Unit       string            `json:"unit"`
	ActivityID domain.ActivityID `json:"activityId"`
	TagIds     []domain.TagID    `json:"tagIds"`
}

// ExpensesByDate handler returns a list of all expenses done
// from a specific date up to now.
// It requires a query parameter "from" specifying the date as mm-dd-yyyy
// If no parameter is found default is to return expenses of last 3 months
func (h *Handler) ExpensesByDate(c echo.Context) error {
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
	expenses, err := h.lister.ExpensesByTime(date)
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	return c.JSON(http.StatusOK, expenses)
}
