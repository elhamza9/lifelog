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
	if err := server.RegisterRoutes(router, hnd); err != nil {
		os.Exit(1)
	}

	port := ":8080"
	router.Start(port)
}
