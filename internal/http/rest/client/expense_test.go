package client_test

import (
	"testing"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/http/rest/client"
	"github.com/elhamza90/lifelog/internal/http/rest/server"
)

func TestPostExpense(t *testing.T) {
	payload := server.JSONReqExpense{
		Label:  "Do smth",
		Value:  100,
		Unit:   "eu",
		Time:   time.Now().Add(time.Duration(-1 * time.Hour)),
		TagIds: []domain.TagID{},
	}
	_, err := client.PostExpense(payload, token)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchExpenses(t *testing.T) {
	_, err := client.FetchExpenses(token, time.Now().AddDate(0, -2, 0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateExpense(t *testing.T) {
	// Create Test Expense
	payload := server.JSONReqExpense{
		Label:  "Do smth",
		Value:  100,
		Unit:   "eu",
		Time:   time.Now().Add(time.Duration(-1 * time.Hour)),
		TagIds: []domain.TagID{},
	}
	id, err := client.PostExpense(payload, token)
	if err != nil {
		t.Fatal(err)
	}
	// Test Update
	payload.ID = id
	payload.Label = "Updated Label"
	payload.Value = payload.Value * 2
	payload.Unit = "dollar"
	payload.Time = time.Now().Add(time.Duration(-2 * time.Hour))
	if err = client.UpdateExpense(payload, token); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteExpense(t *testing.T) {
	// Create Test Expense
	payload := server.JSONReqExpense{
		Label:  "Do smth",
		Value:  100,
		Unit:   "eu",
		Time:   time.Now().Add(time.Duration(-1 * time.Hour)),
		TagIds: []domain.TagID{},
	}
	id, err := client.PostExpense(payload, token)
	if err != nil {
		t.Fatal(err)
	}
	// Test Delete
	if err = client.DeleteExpense(domain.ExpenseID(id), token); err != nil {
		t.Fatal(err)
	}
}

func TestFetchExpenseDetails(t *testing.T) {
	// Create Test Expense
	payload := server.JSONReqExpense{
		Label:  "Do smth",
		Value:  100,
		Unit:   "eu",
		Time:   time.Now().Add(time.Duration(-1 * time.Hour)),
		TagIds: []domain.TagID{},
	}
	id, err := client.PostExpense(payload, token)
	if err != nil {
		t.Fatal(err)
	}
	// Test FetchExpenseDetails
	if _, err = client.FetchExpenseDetails(domain.ExpenseID(id), token); err != nil {
		t.Fatal(err)
	}
}
