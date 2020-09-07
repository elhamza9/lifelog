package main

import (
	"github.com/elhamza90/lifelog/pkg/http/rest"
	"github.com/elhamza90/lifelog/pkg/store/memory"
	"github.com/elhamza90/lifelog/pkg/usecase/adding"
	"github.com/elhamza90/lifelog/pkg/usecase/deleting"
	"github.com/elhamza90/lifelog/pkg/usecase/editing"
	"github.com/elhamza90/lifelog/pkg/usecase/listing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	router := echo.New()
	router.Use(middleware.Logger())

	repo := memory.NewRepository()

	lister := listing.NewService(&repo)
	adder := adding.NewService(&repo)
	editor := editing.NewService(&repo)
	deletor := deleting.NewService(&repo)

	hnd := rest.NewHandler(&lister, &adder, &editor, &deletor)

	rest.RegisterRoutes(router, hnd)

	router.Logger.Fatal(router.Start(":8080"))
}
