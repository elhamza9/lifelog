package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/labstack/echo/v4"
)

// jsonActivity is used to unmarshal a json activity
type jsonActivity struct {
	Label    string         `json:"label"`
	Desc     string         `json:"desc"`
	Place    string         `json:"place"`
	Time     time.Time      `json:"time"`
	Duration time.Duration  `json:"duration"`
	TagIds   []domain.TagID `json:"tagIds"`
}

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
	var date time.Time
	if len(dateStr) == 0 {
		date = defaultActivitiesDateFilter()
	} else {
		var err error
		date, err = time.Parse(dateFilterFormat, dateStr)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	}
	activities, err := h.lister.ActivitiesByTime(date)
	if err != nil {
		return c.String(errToHTTPCode(err, "activities"), err.Error())
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
		return c.String(errToHTTPCode(err, "activities"), err.Error())
	}
	return c.JSON(http.StatusOK, act)
}

// AddActivity handler adds an activity
func (h *Handler) AddActivity(c echo.Context) error {
	// Json unmarshall
	var jsAct jsonActivity
	if err := c.Bind(&jsAct); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Construct Tags slice from ids ( don't fetch anything )
	tags := []domain.Tag{}
	for _, id := range jsAct.TagIds {
		tags = append(tags, domain.Tag{ID: id})
	}
	// Call adding service
	act := domain.Activity{
		Label:    jsAct.Label,
		Desc:     jsAct.Desc,
		Place:    jsAct.Place,
		Time:     jsAct.Time,
		Duration: jsAct.Duration,
		Tags:     tags,
	}
	id, err := h.adder.NewActivity(act)
	if err != nil {
		return c.String(errToHTTPCode(err, "activities"), err.Error())
	}
	// Retrieve created activity
	created, err := h.lister.Activity(id)
	if err != nil {
		return c.String(errToHTTPCode(err, "activities"), err.Error())
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
		return c.String(errToHTTPCode(err, "activities"), err.Error())
	}
	// Json unmarshall
	var jsAct jsonActivity
	if err := c.Bind(&jsAct); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Construct Tags slice from ids ( don't fetch anything )
	tags := []domain.Tag{}
	for _, id := range jsAct.TagIds {
		tags = append(tags, domain.Tag{ID: id})
	}
	// Edit Activity
	act.Label = jsAct.Label
	act.Desc = jsAct.Desc
	act.Place = jsAct.Place
	act.Time = jsAct.Time
	act.Duration = jsAct.Duration
	act.Tags = tags
	err = h.editor.EditActivity(act)
	if err != nil {
		return c.String(errToHTTPCode(err, "activities"), err.Error())
	}
	// Retrieve edited activity
	edited, err := h.lister.Activity(act.ID)
	if err != nil {
		return c.String(errToHTTPCode(err, "activities"), err.Error())
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
		return c.String(errToHTTPCode(err, "activities"), err.Error())
	}
	return c.JSON(http.StatusNoContent, "Activity Deleted Successfully")
}
