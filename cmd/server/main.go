package main

import (
	"os"

	"github.com/elhamza90/lifelog/pkg/http/rest"
	"github.com/elhamza90/lifelog/pkg/store/memory"
	"github.com/elhamza90/lifelog/pkg/usecase/adding"
	"github.com/elhamza90/lifelog/pkg/usecase/auth"
	"github.com/elhamza90/lifelog/pkg/usecase/deleting"
	"github.com/elhamza90/lifelog/pkg/usecase/editing"
	"github.com/elhamza90/lifelog/pkg/usecase/listing"
	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()

	repo := memory.NewRepository()

	lister := listing.NewService(&repo)
	adder := adding.NewService(&repo)
	editor := editing.NewService(&repo)
	deletor := deleting.NewService(&repo)
	authenticator := auth.NewService("LFLG_PASS_HASH")

	hnd := rest.NewHandler(&lister, &adder, &editor, &deletor, &authenticator)

	if err := rest.RegisterRoutes(router, hnd); err != nil {
		os.Exit(1)
	}

	port := ":8080"
	router.Start(port)
}
