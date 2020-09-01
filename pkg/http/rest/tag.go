package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store"
	"github.com/elhamza90/lifelog/pkg/usecase/deleting"
	"github.com/labstack/echo/v4"
)

// GetAllTags handler returns a list of all tags
func (h *Handler) GetAllTags(c echo.Context) error {
	tags, err := h.lister.AllTags()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, tags)
}

// AddTag handler calls adding service to create a tag
// with given name and returns the created tag
func (h *Handler) AddTag(c echo.Context) error {
	// Json unmarshall
	t := new(struct {
		Name string `json:"name"`
	})
	if err := c.Bind(t); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Create Tag
	id, err := h.adder.NewTag((*t).Name)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
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
	t := new(struct {
		Name string `json:"name"`
	})
	if err := c.Bind(t); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Edit Tag
	var tag domain.Tag = domain.Tag{
		ID:   domain.TagID(id),
		Name: t.Name,
	}
	if err := h.editor.EditTag(tag); err != nil {
		// return 404 if not found
		if errors.Is(err, store.ErrTagNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Get edited Tag
	edited, err := h.lister.GetTagByID(tag.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
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

	err = h.deleter.Tag(domain.TagID(id))
	if err != nil {
		if errors.Is(err, store.ErrTagNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		} else if errors.Is(err, deleting.ErrTagHasExpenses) || errors.Is(err, deleting.ErrTagHasActivities) {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		} else {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.String(http.StatusNoContent, "Tag Deleted Successfully")
}
