package adding_test

import (
	"log"
	"os"
	"testing"

	"github.com/elhamza90/lifelog/pkg/store/memory"
	"github.com/elhamza90/lifelog/pkg/usecase/adding"
)

var adder adding.Service   // Instance of service we will be testing
var repo memory.Repository // Repository used by service

func TestMain(m *testing.M) {
	log.Println("Setting up tests")
	repo = memory.NewRepository()    // Work with In-Memory DB
	adder = adding.NewService(&repo) // Passing by reference to change db when testing
	os.Exit(m.Run())
}
