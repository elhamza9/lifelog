package rest

import (
	"net/http"
	"os"

	"github.com/elhamza90/lifelog/pkg/usecase/adding"
	"github.com/elhamza90/lifelog/pkg/usecase/deleting"
	"github.com/elhamza90/lifelog/pkg/usecase/editing"
	"github.com/elhamza90/lifelog/pkg/usecase/listing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

// JwtSecret returns a byte array containing Jwt Signing Key
func JwtSecret() []byte {
	secret := os.Getenv("LFLG_JWT_SECRET")
	return []byte(secret)
}

// RegisterRoutes registers routes with handlers.
func RegisterRoutes(r *echo.Echo, hnd *Handler) {
	secret := JwtSecret()
	r.GET("/health-check", HealthCheck)
	// Group Auth
	auth := r.Group("/auth")
	auth.POST("/login", hnd.Login)
	// Group Tags
	tags := r.Group("/tags", middleware.JWT(secret))
	tags.GET("", hnd.GetAllTags)
	tags.GET("/:id/expenses", hnd.GetTagExpenses)
	tags.GET("/:id/activities", hnd.GetTagActivities)
	tags.POST("", hnd.AddTag)
	tags.PUT("/:id", hnd.EditTag)
	tags.DELETE("/:id", hnd.DeleteTag)
	// Group Activities
	activities := r.Group("/activities", middleware.JWT(secret))
	activities.GET("", hnd.ActivitiesByDate)
	activities.GET("/:id", hnd.ActivityDetails)
	activities.POST("", hnd.AddActivity)
	activities.PUT("/:id", hnd.EditActivity)
	activities.DELETE("/:id", hnd.DeleteActivity)
	// Group Expenses
	expenses := r.Group("/expenses", middleware.JWT(secret))
	expenses.GET("", hnd.ExpensesByDate)
	expenses.GET("/:id", hnd.ExpenseDetails)
	expenses.POST("", hnd.AddExpense)
	expenses.PUT("/:id", hnd.EditExpense)
	expenses.DELETE("/:id", hnd.DeleteExpense)
}

// HealthCheck handler informs that api is up and running.
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Up!")
}
