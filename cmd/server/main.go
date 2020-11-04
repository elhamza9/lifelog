package main

import (
	"errors"
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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// getDbConn retrieves DB params from environment, constructs & returns the DB connection
func getDbConn() (string, error) {
	var (
		dbHost string = os.Getenv("LFLG_DB_HOST")
		dbPort string = os.Getenv("LFLG_DB_PORT")
		dbName string = os.Getenv("LFLG_DB_NAME")
		dbUser string = os.Getenv("LFLG_DB_USER")
		dbPass string = os.Getenv("LFLG_DB_PASS")
	)
	if dbHost == "" || dbPort == "" || dbName == "" || dbUser == "" || dbPass == "" {
		return "", errors.New("Db parameter missing in environment")
	}
	conn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", dbHost, dbPort, dbName, dbUser, dbPass)
	return conn, nil
}

// hash_var_name specifies the name of the environment variable where the password bcrypt hash is stored
const hash_var_name string = "LFLG_PASS_HASH"

func main() {
	dbConn, err := getDbConn()
	if err != nil {
		fmt.Printf("could not construct Db connection string: %s\n", err)
		os.Exit(1)
	}
	grmDb, err := gorm.Open(postgres.Open(dbConn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
		os.Exit(1)
	}
	if err := grmDb.AutoMigrate(&db.Tag{}, &db.Expense{}, &db.Activity{}); err != nil {
		fmt.Printf("Error Auto-Migrating Tables:\n\t%s\n", err)
		os.Exit(1)
	}

	repo := db.NewRepository(grmDb)

	lister := listing.NewService(&repo)
	adder := adding.NewService(&repo)
	editor := editing.NewService(&repo)
	deletor := deleting.NewService(&repo)
	authenticator := auth.NewService(hash_var_name)

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
