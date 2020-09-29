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

// defaultActivitiesMinDate specifies default date filter when listing activities
// and no filter was provided
func defaultActivitiesDateFilter() time.Time {
	return time.Now().AddDate(0, -3, 0)
}

// ActivitiesByDate handler returns a list of all activities done
// from a specific date up to now.
// It requires a query parameter "from" specifying the date as mm-dd-yyyy
// If no parameter is found default is to return activities of last 3 months
func (h *Handler) ActivitiesByDate(c echo.Context) error {
	dateStr := c.QueryParam("from")
	logrus.Debugf("Extracted query param from: %s", dateStr)
	var date time.Time
	if len(dateStr) == 0 {
		date = defaultActivitiesDateFilter()
	} else {
		var err error
		date, err = time.Parse(dateFilterFormat, dateStr)
		if err != nil {
			msg := "Error parsing from param to valid date"
			details := err.Error()
			logrus.Error(msg + ": " + details)
			return c.String(http.StatusBadRequest, msg)
		}
	}
	// Fetch Activities
	activities, err := h.lister.ActivitiesByTime(date)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching activities since %s", date.Format("2006-01-02"))
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "activities"), msg)
	}
	logrus.Infof("Fetched activities since %s successfully", date.Format("2006-01-02"))
	// Construct respActivities from fetched activities
	respActivities := make([]JSONRespListActivity, len(activities))
	var respAct JSONRespListActivity
	for i, act := range activities {
		respAct.From(act)
		respActivities[i] = respAct
	}
	return c.JSON(http.StatusOK, respActivities)
}

// ActivityDetails handler returns details of activity with given ID
// It required a path parameter :id
func (h *Handler) ActivityDetails(c echo.Context) error {
	// Get ID from Path param
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := fmt.Sprintf("Error while converting path param Activity ID with value %s to int", idStr)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(http.StatusBadRequest, msg)
	}
	actID := domain.ActivityID(id)
	// Get Activity
	act, err := h.lister.Activity(actID)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching activity %s", actID)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "activities"), msg)
	}
	logrus.Infof("Fetched activity %s successfully", actID)
	// Fetch Expenses
	expenses, err := h.lister.ExpensesByActivity(act.ID)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching expenses of activity %s", actID)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "activities"), err.Error())
	}
	logrus.Infof("Fetched expenses of activity %s successfully", actID)
	var actResp JSONRespDetailActivity
	actResp.From(act, expenses)
	return c.JSON(http.StatusOK, actResp)
}

// AddActivity handler adds an activity
func (h *Handler) AddActivity(c echo.Context) error {
	// Json unmarshall
	var jsAct JSONReqActivity
	if err := c.Bind(&jsAct); err != nil {
		var (
			msg     string = errInvalidJSON.Error()
			details string = httpErrorMsg(err)
			code    int    = errToHTTPCode(errInvalidJSON, "activities")
		)
		logrus.Error(msg + " : " + details)
		return c.String(code, msg)
	}
	act := jsAct.ToDomain()
	id, err := h.adder.NewActivity(act)
	if err != nil {
		msg := "Internal Server Error while adding activity"
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "activities"), msg)
	}
	logrus.Infof("Created activity %s successfully", id)
	// Retrieve created activity
	created, err := h.lister.Activity(id)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching activity %s", id)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "activities"), msg)
	}
	return c.JSON(http.StatusCreated, created)
}

// EditActivity handler edits an activity with given ID
// It required a path parameter :id
func (h *Handler) EditActivity(c echo.Context) error {
	// Get ID from Path param
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := fmt.Sprintf("Error while converting path param Activity ID with value %s to int", idStr)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(http.StatusBadRequest, msg)
	}
	actID := domain.ActivityID(id)
	// Get Activity
	act, err := h.lister.Activity(actID)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching activity %s", actID)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "activities"), msg)
	}
	logrus.Infof("Fetched activity %s successfully", actID)
	// Json unmarshall
	var jsAct JSONReqActivity
	if err := c.Bind(&jsAct); err != nil {
		var (
			msg     string = errInvalidJSON.Error()
			details string = httpErrorMsg(err)
			code    int    = errToHTTPCode(errInvalidJSON, "activities")
		)
		logrus.Error(msg + " | " + details)
		return c.String(code, msg)
	}
	updated := jsAct.ToDomain()
	err = h.editor.EditActivity(updated)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while updating activity %s", actID)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "activities"), msg)
	}
	logrus.Infof("Updated activity %s successfully", actID)
	// Retrieve edited activity
	edited, err := h.lister.Activity(act.ID)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching updated activity %s", act.ID)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "activities"), msg)
	}
	logrus.Infof("Fetched activity %s successfully", act.ID)
	return c.JSON(http.StatusOK, edited)
}

// DeleteActivity handler deletes an activity with given ID
// It required a path parameter :id
func (h *Handler) DeleteActivity(c echo.Context) error {
	// Get ID from Path param
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := fmt.Sprintf("Error while converting path param Activity ID with value %s to int", idStr)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(http.StatusBadRequest, msg)
	}
	actID := domain.ActivityID(id)
	// Delete Activity
	err = h.deleter.Activity(actID)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while deleting activity %s", actID)
		logrus.Error(msg + " : " + err.Error())
		return c.String(errToHTTPCode(err, "activities"), msg)
	}
	logrus.Infof("Deleted activity %s successfully", actID)
	return c.JSON(http.StatusNoContent, "Activity Deleted Successfully")
}
