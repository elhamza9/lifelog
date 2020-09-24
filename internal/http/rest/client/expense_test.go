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

func TestFetchExpenses(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"
	_, err := client.FetchExpenses(token, time.Now().AddDate(0, -2, 0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteExpense(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"
	// Create Test Expense
	testExp := domain.Expense{
		Label: "Do smth",
		Value: 100,
		Unit:  "eu",
		Time:  time.Now().Add(time.Duration(-1 * time.Hour)),
		Tags:  []domain.Tag{},
	}
	id, err := client.PostExpense(testExp, token)
	if err != nil {
		t.Fatal(err)
	}
	// Test Delete
	if err = client.DeleteExpense(domain.ExpenseID(id), token); err != nil {
		t.Fatal(err)
	}
}

func TestFetchExpenseDetails(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"
	// Create Test Expense
	testExp := domain.Expense{
		Label: "Do smth",
		Value: 100,
		Unit:  "eu",
		Time:  time.Now().Add(time.Duration(-1 * time.Hour)),
		Tags:  []domain.Tag{},
	}
	id, err := client.PostExpense(testExp, token)
	if err != nil {
		t.Fatal(err)
	}
	// Test FetchExpenseDetails
	if _, err = client.FetchExpenseDetails(domain.ExpenseID(id), token); err != nil {
		t.Fatal(err)
	}
}
