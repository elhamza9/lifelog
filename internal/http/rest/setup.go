package rest

import (
	"errors"
	"net/http"

	"github.com/elhamza90/lifelog/internal/usecase/adding"
	"github.com/elhamza90/lifelog/internal/usecase/auth"
	"github.com/elhamza90/lifelog/internal/usecase/deleting"
	"github.com/elhamza90/lifelog/internal/usecase/editing"
	"github.com/elhamza90/lifelog/internal/usecase/listing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

// Handler contains services required by it's methods
// (which are http handlers) to perform their jobs.
type Handler struct {
	lister        listing.Service
	adder         adding.Service
	editor        editing.Service
	deleter       deleting.Service
	authenticator auth.Service
}

// NewHandler construct & returns a new handler with provided services.
func NewHandler(lister *listing.Service, adder *adding.Service, editor *editing.Service, deleter *deleting.Service, authenticator *auth.Service) *Handler {
	return &Handler{
		lister:        *lister,
		adder:         *adder,
		editor:        *editor,
		deleter:       *deleter,
		authenticator: *authenticator,
	}
}

// RegisterRoutes registers routes with handlers.
func RegisterRoutes(r *echo.Echo, hnd *Handler) error {
	secret := jwtAccessSecret()
	if len(secret) == 0 {
		msg := "No JWT Secret was found in system"
		log.Fatal(msg)
		return errors.New(msg)
	}
	r.GET("/health-check", HealthCheck)
	// Group Auth
	auth := r.Group("/auth")
	auth.POST("/login", hnd.Login)
	auth.POST("/refresh", hnd.RefreshToken)
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
	return nil
}

// HealthCheck handler informs that api is up and running.
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Up!")
}
