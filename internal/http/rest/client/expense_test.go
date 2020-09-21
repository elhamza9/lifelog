package client_test

import (
	"testing"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/http/rest/client"
)

func TestPostExpense(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.gkSvHuAki4boAJkAbMgSsFRCqA80GjEHfg9rvjvICpY"
	testExp := domain.Expense{
		Label: "Do smth",
		Value: 100,
		Unit:  "eu",
		Time:  time.Now().Add(time.Duration(-1 * time.Hour)),
		Tags:  []domain.Tag{},
	}
	_, err := client.PostExpense(testExp, token)
	if err != nil {
		t.Fatal(err)
	}
}
