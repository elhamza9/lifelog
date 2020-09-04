package rest

import (
	"net/http"

	"github.com/elhamza90/lifelog/pkg/usecase/adding"
	"github.com/elhamza90/lifelog/pkg/usecase/deleting"
	"github.com/elhamza90/lifelog/pkg/usecase/editing"
	"github.com/elhamza90/lifelog/pkg/usecase/listing"
	"github.com/labstack/echo/v4"
)

// Handler contains services required by it's methods
// (which are http handlers) to perform their jobs.
type Handler struct {
	lister  listing.Service
	adder   adding.Service
	editor  editing.Service
	deleter deleting.Service
}

// NewHandler construct & returns a new handler with provided services.
func NewHandler(lister *listing.Service, adder *adding.Service, editor *editing.Service, deleter *deleting.Service) *Handler {
	return &Handler{
		lister:  *lister,
		adder:   *adder,
		editor:  *editor,
		deleter: *deleter,
	}
}

// RegisterRoutes registers routes with handlers.
func RegisterRoutes(r *echo.Echo, hnd *Handler) {
	r.GET("/health-check", HealthCheck)
	// Group Tags
	tags := r.Group("/tags")
	tags.GET("", hnd.GetAllTags)
	tags.POST("", hnd.AddTag)
	tags.PUT("/:id", hnd.EditTag)
	tags.DELETE("/:id", hnd.DeleteTag)
	// Group Activities
	activities := r.Group("/activities")
	activities.GET("", hnd.ActivitiesByDate)
	activities.GET("/:id", hnd.ActivityDetails)
	activities.POST("", hnd.AddActivity)
	activities.PUT("/:id", hnd.EditActivity)
	activities.DELETE("/:id", hnd.DeleteActivity)
	// Group Expenses
	expenses := r.Group("/expenses")
	expenses.GET("", hnd.ExpensesByDate)
	expenses.POST("", hnd.AddExpense)
}

// HealthCheck handler informs that api is up and running.
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Up!")
}
