package main

import (
	"fmt"
	"os"

	"github.com/elhamza90/lifelog/internal/http/rest/server"
	"github.com/elhamza90/lifelog/internal/store/db"
	"github.com/elhamza90/lifelog/internal/usecase/adding"
	"github.com/elhamza90/lifelog/internal/usecase/auth"
	"github.com/elhamza90/lifelog/internal/usecase/deleting"
	"github.com/elhamza90/lifelog/internal/usecase/editing"
	"github.com/elhamza90/lifelog/internal/usecase/listing"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	grmDb, err := gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
		os.Exit(1)
	}
	grmDb.AutoMigrate(&db.Tag{}, &db.Expense{}, &db.Activity{})

	repo := db.NewRepository(grmDb)

	lister := listing.NewService(&repo)
	adder := adding.NewService(&repo)
	editor := editing.NewService(&repo)
	deletor := deleting.NewService(&repo)
	authenticator := auth.NewService("LFLG_PASS_HASH")

	hnd := server.NewHandler(&lister, &adder, &editor, &deletor, &authenticator)

	router := echo.New()

	// Setup Routes
	if err := server.RegisterRoutes(router, hnd); err != nil {
		os.Exit(1)
	}

	// Setup Logger
	logrus.SetFormatter(&logrus.TextFormatter{
		//DisableColors: true,
		FullTimestamp: true,
	})
	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := logrus.WithFields(logrus.Fields{
				"remote_ip": c.RealIP(),
				"path":      c.Request().URL.Path,
				"method":    c.Request().Method,
			})
			logger.Info("New Request")
			return next(c)
		}
	})

	port := ":8080"
	router.Start(port)
}
