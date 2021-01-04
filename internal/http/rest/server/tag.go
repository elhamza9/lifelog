package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// GetAllTags handler returns a list of all tags.
func (h *Handler) GetAllTags(c echo.Context) error {
	tags, err := h.lister.AllTags()
	if err != nil {
		msg := "Internal Server Error while fetching tags"
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(http.StatusInternalServerError, msg)
	}
	logrus.Info("All tags fetched successfully")
	respTags := make([]JSONRespListTag, len(tags))
	var respTag JSONRespListTag
	for i, t := range tags {
		respTag.From(t)
		respTags[i] = respTag
	}
	return c.JSON(http.StatusOK, tags)
}

// GetTagExpenses handler returns expenses of a given tag.
func (h *Handler) GetTagExpenses(c echo.Context) error {
	// Get Tag ID from path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := fmt.Sprintf("Error while converting path param Tag ID with value %s to int", idStr)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(http.StatusBadRequest, msg)
	}
	tagID := domain.TagID(id)
	logrus.Debugf("Extracted tag id from path param: %s", tagID)
	// Get Expenses
	expenses, err := h.lister.ExpensesByTag(tagID)
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching expenses of tag %s", tagID)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(errToHTTPCode(err, "tags"), msg)
	}
	logrus.Infof("Expenses of tag with ID %s fetched successfully", tagID)
	// Construct response expenses from fetched expenses
	respExpenses := make([]JSONRespListExpense, len(expenses))
	var respExp JSONRespListExpense
	for i, exp := range expenses {
		respExp.From(exp)
		respExpenses[i] = respExp
	}
	return c.JSON(http.StatusOK, respExpenses)
}

// GetTagActivities handler returns activities of a given tag.
func (h *Handler) GetTagActivities(c echo.Context) error {
	// Get Tag ID from path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := fmt.Sprintf("Error while converting path param Tag ID with value %s to int", idStr)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(http.StatusBadRequest, msg)
	}
	tagID := domain.TagID(id)
	logrus.Debugf("Extracted tag id from path param: %s", tagID)
	// Get Activities
	activities, err := h.lister.ActivitiesByTag(domain.TagID(id))
	if err != nil {
		msg := fmt.Sprintf("Internal Server Error while fetching activities of tag %s", tagID)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(errToHTTPCode(err, "tags"), msg)
	}
	logrus.Infof("Activities of tag with ID %s fetched successfully", tagID)
	// Construct respActivities from fetched activities
	respActivities := make([]JSONRespListActivity, len(activities))
	var respAct JSONRespListActivity
	for i, act := range activities {
		respAct.From(act)
		respActivities[i] = respAct
	}
	return c.JSON(http.StatusOK, respActivities)
}

// AddTag handler adds a given tag and returns it.
func (h *Handler) AddTag(c echo.Context) error {
	// Json unmarshall
	var jsTag JSONReqTag
	if err := c.Bind(&jsTag); err != nil {
		var (
			msg     string = errInvalidJSON.Error()
			details string = httpErrorMsg(err)
			code    int    = errToHTTPCode(errInvalidJSON, "tags")
		)
		logrus.Error(msg + " | " + details)
		return c.String(code, msg)
	}
	// Create Tag
	tag := jsTag.ToDomain()
	id, err := h.adder.NewTag(tag)
	if err != nil {
		msg := err.Error()
		logrus.Error(msg)
		return c.String(errToHTTPCode(err, "tags"), msg)
	}
	logrus.Infof("Created Tag %s successfully", id)
	// Get created Tag
	created, err := h.lister.GetTagByID(id)
	logrus.Infof("Retrieved Tag %s successfully", created.ID)
	return c.JSON(http.StatusCreated, created)
}

// EditTag handler edits tag with given ID and returns it.
func (h *Handler) EditTag(c echo.Context) error {
	// Get Tag ID from path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := fmt.Sprintf("Error while converting path param Tag ID with value %s to int", idStr)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(http.StatusBadRequest, msg)
	}
	tagID := domain.TagID(id)
	logrus.Debugf("Extracted tag id from path param: %s", tagID)
	// Json unmarshall
	var jsTag JSONReqTag
	if err := c.Bind(&jsTag); err != nil {
		var (
			msg     string = errInvalidJSON.Error()
			details string = httpErrorMsg(err)
			code    int    = errToHTTPCode(errInvalidJSON, "tags")
		)
		logrus.Error(msg + " | " + details)
		return c.String(code, msg)
	}
	logrus.Debug("Unmarshalled JSON successfully")
	// Edit Tag
	tag := jsTag.ToDomain()
	tag.ID = tagID
	if err := h.editor.EditTag(tag); err != nil {
		msg := fmt.Sprintf("error while updating tag %s", tagID)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(errToHTTPCode(err, "tags"), msg)
	}
	logrus.Infof("Updated Tag %s successfully", tagID)
	// Retrieve edited Tag
	edited, err := h.lister.GetTagByID(tag.ID)
	if err != nil {
		msg := fmt.Sprintf("error while retrieving updated tag %s", tagID)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(errToHTTPCode(err, "tags"), msg)
	}
	logrus.Infof("Retrieved Tag %s successfully", tagID)
	return c.JSON(http.StatusOK, edited)
}

// DeleteTag handler deletes a tag with given ID.
func (h *Handler) DeleteTag(c echo.Context) error {
	// Get Tag ID from path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := fmt.Sprintf("Error while converting path param Tag ID with value %s to int", idStr)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(http.StatusBadRequest, msg)
	}
	tagID := domain.TagID(id)
	logrus.Debugf("Extracted tag id from path param: %s", tagID)
	// Delete Tag
	err = h.deleter.Tag(domain.TagID(id))
	if err != nil {
		msg := fmt.Sprintf("error while deleting tag with ID: %s", tagID)
		details := err.Error()
		logrus.Error(msg + " | " + details)
		return c.String(errToHTTPCode(err, "tags"), msg)
	}
	logrus.Infof("Deleted tag %s successfully", tagID)
	return c.String(http.StatusNoContent, "Tag Deleted Successfully")
}
