package rest

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store"
	"github.com/elhamza90/lifelog/pkg/usecase/deleting"
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

// ActivityDetails handler returns details of activity with given ID
// It required a path parameter :id
func (h *Handler) ActivityDetails(c echo.Context) error {
	// Get ID from Path param
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Get Activity
	act, err := h.lister.Activity(domain.ActivityID(id))
	if err != nil {
		if errors.Is(err, store.ErrActivityNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, act)
}

// AddActivity handler adds an activity
func (h *Handler) AddActivity(c echo.Context) error {
	// Json unmarshall
	a := new(struct {
		Label    string         `json:"label"`
		Desc     string         `json:"desc"`
		Place    string         `json:"place"`
		Time     time.Time      `json:"time"`
		Duration time.Duration  `json:"duration"`
		TagIds   []domain.TagID `json:"tagIds"`
	})
	if err := c.Bind(a); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Construct Tags slice from ids ( don't fetch anything )
	tags := []domain.Tag{}
	for _, id := range (*a).TagIds {
		tags = append(tags, domain.Tag{ID: id})
	}

	// Call adding service
	act := domain.Activity{
		Label:    (*a).Label,
		Desc:     (*a).Desc,
		Place:    (*a).Place,
		Time:     (*a).Time,
		Duration: (*a).Duration,
		Tags:     tags,
	}
	id, err := h.adder.NewActivity(act.Label, act.Place, act.Desc, act.Time, act.Duration, tags)
	if err != nil {
		if errors.Is(err, store.ErrTagNotFound) {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
		if errors.Is(err, domain.ErrActivityTimeFuture) {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Retrieve created activity
	created, err := h.lister.Activity(id)
	if err != nil {
		if errors.Is(err, store.ErrActivityNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
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
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Get Activity
	act, err := h.lister.Activity(domain.ActivityID(id))
	if err != nil {
		if errors.Is(err, store.ErrActivityNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Json unmarshall
	a := new(struct {
		Label    string         `json:"label"`
		Desc     string         `json:"desc"`
		Place    string         `json:"place"`
		Time     time.Time      `json:"time"`
		Duration time.Duration  `json:"duration"`
		TagIds   []domain.TagID `json:"tagIds"`
	})
	if err := c.Bind(a); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Construct Tags slice from ids ( don't fetch anything )
	tags := []domain.Tag{}
	for _, id := range (*a).TagIds {
		tags = append(tags, domain.Tag{ID: id})
	}

	// Edit Activity
	act.Label = (*a).Label
	act.Desc = (*a).Desc
	act.Place = (*a).Place
	act.Time = (*a).Time
	act.Duration = (*a).Duration
	act.Tags = tags
	err = h.editor.EditActivity(act)
	if err != nil {
		if errors.Is(err, store.ErrTagNotFound) {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Retrieve edited activity
	edited, err := h.lister.Activity(act.ID)
	if err != nil {
		if errors.Is(err, store.ErrActivityNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, edited)
}

// DeleteActivity handler deletes an activity with given ID
// It required a path parameter :id
func (h *Handler) DeleteActivity(c echo.Context) error {
	// Get ID from Path param
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Delete Activity
	err = h.deleter.Activity(domain.ActivityID(id))
	if err != nil {
		if errors.Is(err, store.ErrActivityNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, deleting.ErrActivityHasExpenses) {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusNoContent, "Activity Deleted Successfully")
}
