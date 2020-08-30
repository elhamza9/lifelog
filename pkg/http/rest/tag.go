package rest

import (
	"net/http"

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
