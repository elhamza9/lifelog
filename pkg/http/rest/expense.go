package rest

import (
	"net/http"
	"strconv"
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

// ExpenseDetails handler returns details of expense with given ID
// It required a path parameter :id
func (h *Handler) ExpenseDetails(c echo.Context) error {
	// Get ID from Path param
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Get Expense
	act, err := h.lister.Expense(domain.ExpenseID(id))
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	return c.JSON(http.StatusOK, act)
}

// AddExpense handler adds an expense
func (h *Handler) AddExpense(c echo.Context) error {
	// Json unmarshall
	a := new(jsonExpense)
	if err := c.Bind(a); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Construct Tags slice from ids ( don't fetch anything )
	tags := []domain.Tag{}
	for _, id := range (*a).TagIds {
		tags = append(tags, domain.Tag{ID: id})
	}
	// Call adding service
	exp := domain.Expense{
		Label:      (*a).Label,
		Value:      (*a).Value,
		Unit:       (*a).Unit,
		Time:       (*a).Time,
		ActivityID: (*a).ActivityID,
		Tags:       tags,
	}
	id, err := h.adder.NewExpense(exp)
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	// Retrieve created expense
	created, err := h.lister.Expense(id)
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	return c.JSON(http.StatusCreated, created)
}

// EditExpense handler edits an activity with given ID
// It required a path parameter :id
func (h *Handler) EditExpense(c echo.Context) error {
	// Get ID from Path param
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Get Expense
	exp, err := h.lister.Expense(domain.ExpenseID(id))
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	// Json unmarshall
	e := new(jsonExpense)
	if err := c.Bind(e); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Construct Tags slice from ids ( don't fetch anything )
	tags := []domain.Tag{}
	for _, id := range (*e).TagIds {
		tags = append(tags, domain.Tag{ID: id})
	}
	// Edit Expense
	exp.Label = (*e).Label
	exp.Value = (*e).Value
	exp.Unit = (*e).Unit
	exp.Time = (*e).Time
	exp.ActivityID = (*e).ActivityID
	exp.Tags = tags
	err = h.editor.EditExpense(exp)
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	// Retrieve edited activity
	edited, err := h.lister.Expense(exp.ID)
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	return c.JSON(http.StatusOK, edited)
}
