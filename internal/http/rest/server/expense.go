package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
		msg := fmt.Sprintf("Internal Server Error while fetching expenses since %s", date.Format("2006-01-02"))
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "expenses"), msg)
	}
	logrus.Infof("Fetched expenses since %s successfully", date.Format("2006-01-02"))
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
		msg := fmt.Sprintf("Error while converting path param Expense ID with value %s to int", idStr)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(http.StatusBadRequest, msg)
	}
	expID := domain.ExpenseID(id)
	// Get Expense
	exp, err := h.lister.Expense(domain.ExpenseID(id))
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching expense %s", expID)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "expenses"), msg)
	}
	logrus.Infof("Fetched expense %s successfully", expID)
	// Get Expense Activity if exists
	act := domain.Activity{}
	if exp.ActivityID > 0 {
		act, err = h.lister.Activity(exp.ActivityID)
		if err != nil {
			msg := fmt.Sprintf("Internal Server Error while fetching expense %s's activity %s", expID, exp.ActivityID)
			logrus.Error(msg + " : " + err.Error())
			return c.String(errToHTTPCode(err, "expenses"), msg)
		}
		logrus.Infof("Fetched expense %s's activity %s successfully", expID, exp.ActivityID)
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
		var (
			msg     string = errInvalidJSON.Error()
			details string = httpErrorMsg(err)
			code    int    = errToHTTPCode(errInvalidJSON, "expenses")
		)
		logrus.Error(msg + " : " + details)
		return c.String(code, msg)
	}
	// Call adding service
	exp := jsExp.ToDomain()
	id, err := h.adder.NewExpense(exp)
	if err != nil {
		msg := "Internal Server Error while adding expense"
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "expenses"), msg)
	}
	logrus.Infof("Created expense %s successfully", id)
	// Retrieve created expense
	created, err := h.lister.Expense(id)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching expense %s", id)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "expenses"), msg)
	}
	logrus.Infof("Fetched expense %s successfully", id)
	// Get Expense Activity if exists
	act := domain.Activity{Label: "No Activity"}
	if exp.ActivityID > 0 {
		act, err = h.lister.Activity(exp.ActivityID)
		if err != nil {
			msg := fmt.Sprintf("Internal Server Error while fetching expense %s's activity %s", id, exp.ActivityID)
			logrus.Error(msg + " : " + err.Error())
			return c.String(errToHTTPCode(err, "expenses"), msg)
		}
		logrus.Infof("Fetched expense %s's activity %s successfully", id, exp.ActivityID)
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
		msg := fmt.Sprintf("Error while converting path param Expense ID with value %s to int", idStr)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(http.StatusBadRequest, msg)
	}
	expID := domain.ExpenseID(id)
	// Check Expense with given ID exists
	_, err = h.lister.Expense(expID)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching expense %s", expID)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "expenses"), msg)
	}
	logrus.Infof("Fetched expense %s successfully", expID)
	// Json unmarshall
	var jsExp JSONReqExpense
	if err := c.Bind(&jsExp); err != nil {
		var (
			msg     string = errInvalidJSON.Error()
			details string = httpErrorMsg(err)
			code    int    = errToHTTPCode(errInvalidJSON, "activities")
		)
		logrus.Error(msg + " | " + details)
		return c.String(code, msg)
	}
	// Update
	exp := jsExp.ToDomain()
	err = h.editor.EditExpense(exp)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while updating expense %s", expID)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "expenses"), msg)
	}
	logrus.Infof("Updated expense %s successfully", exp.ID)
	// Retrieve edited expense
	edited, err := h.lister.Expense(exp.ID)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching updated expense %s", exp.ID)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "expenses"), msg)
	}
	logrus.Infof("Fetched expense %s successfully", exp.ID)
	// Get Expense Activity if exists
	act := domain.Activity{Label: "No Activity"}
	if exp.ActivityID > 0 {
		act, err = h.lister.Activity(exp.ActivityID)
		if err != nil {
			msg := fmt.Sprintf("Internal Server Error while fetching expense %s's activity %s", exp.ID, exp.ActivityID)
			logrus.Error(msg + " : " + err.Error())
			return c.String(errToHTTPCode(err, "expenses"), msg)
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
		msg := fmt.Sprintf("Error while converting path param Expense ID with value %s to int", idStr)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(http.StatusBadRequest, msg)
	}
	expID := domain.ExpenseID(id)
	// Delete Expense
	err = h.deleter.Expense(expID)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while deleting expense %s", expID)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "expenses"), msg)
	}
	logrus.Infof("Deleted expense %s successfully", expID)
	return c.JSON(http.StatusNoContent, "Expense Deleted Successfully")
}
