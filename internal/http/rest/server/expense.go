package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/labstack/echo/v4"
)

// defaultExpensesMinDate specifies default date filter when listing expenses
// and no filter was provided
func defaultExpensesDateFilter() time.Time {
	return time.Now().AddDate(0, -3, 0)
}

// ExpensesByDate handler returns a list of all expenses done
// from a specific date up to now.
// It requires a query parameter "from" specifying the date as mm-dd-yyyy
// If no parameter is found default is to return expenses of last 3 months
func (h *Handler) ExpensesByDate(c echo.Context) error {
	dateStr := c.QueryParam("from")
	var date time.Time
	if len(dateStr) == 0 {
		date = defaultExpensesDateFilter()
	} else {
		var err error
		if date, err = time.Parse(dateFilterFormat, dateStr); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	}
	expenses, err := h.lister.ExpensesByTime(date)
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	// Construct response expenses from fetched expenses
	respExpenses := make([]JSONRespListExpense, len(expenses))
	var respExp JSONRespListExpense
	for i, exp := range expenses {
		respExp.From(exp)
		respExpenses[i] = respExp
	}
	return c.JSON(http.StatusOK, respExpenses)
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
	exp, err := h.lister.Expense(domain.ExpenseID(id))
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	// Get Expense Activity if exists
	act := domain.Activity{}
	if exp.ActivityID > 0 {
		act, err = h.lister.Activity(exp.ActivityID)
		if err != nil {
			return c.String(errToHTTPCode(err, "expenses"), err.Error())
		}
	}
	var respExp JSONRespDetailExpense
	respExp.From(exp, act)
	return c.JSON(http.StatusOK, respExp)
}

// AddExpense handler adds an expense
func (h *Handler) AddExpense(c echo.Context) error {
	// Json unmarshall
	var jsExp JSONReqExpense
	if err := c.Bind(&jsExp); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Call adding service
	exp := jsExp.ToDomain()
	id, err := h.adder.NewExpense(exp)
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	// Retrieve created expense
	created, err := h.lister.Expense(id)
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	// Get Expense Activity if exists
	act := domain.Activity{Label: "No Activity"}
	if exp.ActivityID > 0 {
		act, err = h.lister.Activity(exp.ActivityID)
		if err != nil {
			return c.String(errToHTTPCode(err, "expenses"), err.Error())
		}
	}
	var respExp JSONRespDetailExpense
	respExp.From(created, act)
	return c.JSON(http.StatusCreated, respExp)
}

// EditExpense handler edits an expense with given ID
// It required a path parameter :id
func (h *Handler) EditExpense(c echo.Context) error {
	// Get ID from Path param
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Check Expense with given ID exists
	_, err = h.lister.Expense(domain.ExpenseID(id))
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	// Json unmarshall
	var jsExp JSONReqExpense
	if err := c.Bind(&jsExp); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Update
	exp := jsExp.ToDomain()
	err = h.editor.EditExpense(exp)
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	// Retrieve edited expense
	edited, err := h.lister.Expense(exp.ID)
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	// Get Expense Activity if exists
	act := domain.Activity{Label: "No Activity"}
	if exp.ActivityID > 0 {
		act, err = h.lister.Activity(exp.ActivityID)
		if err != nil {
			return c.String(errToHTTPCode(err, "expenses"), err.Error())
		}
	}
	var respExp JSONRespDetailExpense
	respExp.From(edited, act)
	return c.JSON(http.StatusOK, respExp)
}

// DeleteExpense handler deletes an expense with given ID
// It required a path parameter :id
func (h *Handler) DeleteExpense(c echo.Context) error {
	// Get ID from Path param
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Delete Expense
	err = h.deleter.Expense(domain.ExpenseID(id))
	if err != nil {
		return c.String(errToHTTPCode(err, "expenses"), err.Error())
	}
	return c.JSON(http.StatusNoContent, "Expense Deleted Successfully")
}
