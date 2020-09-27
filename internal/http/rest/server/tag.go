package server

import (
	"net/http"
	"strconv"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/labstack/echo/v4"
)

// GetAllTags handler returns a list of all tags
func (h *Handler) GetAllTags(c echo.Context) error {
	tags, err := h.lister.AllTags()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	respTags := make([]JSONRespListTag, len(tags))
	var respTag JSONRespListTag
	for i, t := range tags {
		respTag.From(t)
		respTags[i] = respTag
	}
	return c.JSON(http.StatusOK, tags)
}

// GetTagExpenses handler returns expenses of a given tag
func (h *Handler) GetTagExpenses(c echo.Context) error {
	// Get Tag ID from path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Get Expenses
	expenses, err := h.lister.ExpensesByTag(domain.TagID(id))
	if err != nil {
		return c.String(errToHTTPCode(err, "tags"), err.Error())
	}
	return c.JSON(http.StatusOK, expenses)
}

// GetTagActivities handler returns activities of a given tag
func (h *Handler) GetTagActivities(c echo.Context) error {
	// Get Tag ID from path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Get Activities
	activities, err := h.lister.ActivitiesByTag(domain.TagID(id))
	if err != nil {
		return c.String(errToHTTPCode(err, "tags"), err.Error())
	}
	return c.JSON(http.StatusOK, activities)
}

// AddTag handler calls adding service to create a tag
// with given name and returns the created tag
func (h *Handler) AddTag(c echo.Context) error {
	// Json unmarshall
	var jsTag JSONReqTag
	if err := c.Bind(&jsTag); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Create Tag
	tag := jsTag.ToDomain()
	id, err := h.adder.NewTag(tag)
	if err != nil {
		return c.String(errToHTTPCode(err, "tags"), err.Error())
	}
	// Get created Tag
	created, err := h.lister.GetTagByID(id)
	return c.JSON(http.StatusCreated, created)
}

// EditTag handler calls editing service to edit a tag
// with given name and returns the edited tag
func (h *Handler) EditTag(c echo.Context) error {
	// Get Tag ID from path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Json unmarshall
	var jsTag JSONReqTag
	if err := c.Bind(&jsTag); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Edit Tag
	tag := jsTag.ToDomain()
	tag.ID = domain.TagID(id)
	if err := h.editor.EditTag(tag); err != nil {
		return c.String(errToHTTPCode(err, "tags"), err.Error())
	}
	// Retrieve edited Tag
	edited, err := h.lister.GetTagByID(tag.ID)
	if err != nil {
		return c.String(errToHTTPCode(err, "tags"), err.Error())
	}
	return c.JSON(http.StatusOK, edited)
}

// DeleteTag handler calls deleting service to delete a tag
func (h *Handler) DeleteTag(c echo.Context) error {
	// Get Tag ID from path
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Delete Tag
	err = h.deleter.Tag(domain.TagID(id))
	if err != nil {
		return c.String(errToHTTPCode(err, "tags"), err.Error())
	}
	return c.String(http.StatusNoContent, "Tag Deleted Successfully")
}
