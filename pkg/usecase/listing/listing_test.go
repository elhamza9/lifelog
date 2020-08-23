package listing_test

import (
	"os"
	"testing"

	"github.com/elhamza90/lifelog/pkg/store/memory"
	"github.com/elhamza90/lifelog/pkg/usecase/listing"
)

var service listing.Service
var repo memory.Repository

func TestMain(m *testing.M) {
	repo = memory.NewRepository()       // Work with In-Memory DB
	service = listing.NewService(&repo) // Passing by reference to change db when testing
	os.Exit(m.Run())
}
