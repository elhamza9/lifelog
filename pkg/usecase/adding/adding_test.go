package adding_test

import (
	"log"
	"os"
	"testing"

	"github.com/elhamza90/lifelog/pkg/store/memory"
	"github.com/elhamza90/lifelog/pkg/usecase/adding"
)

var adder adding.Service
var repo memory.Repository

func TestMain(m *testing.M) {
	log.Println("Setting up tests")
	repo = memory.NewRepository()      // Work with In-Memory DB
	adder = adding.NewService(&repo) // Passing by reference to change db when testing
	os.Exit(m.Run())
}
