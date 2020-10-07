package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/http/rest/server"
)

// PostExpense sends a POST request with refresh token and
// gets a new Jwt Access Token
func PostExpense(payload server.JSONReqExpense, token string) (domain.ExpenseID, error) {
	// Marshall Expense to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	// Send HTTP Request
	path := url + "/expenses"
	requestBody := bytes.NewBuffer(jsonPayload)
	req, err := http.NewRequest("POST", path, requestBody)
	if err != nil {
		return 0, err
	}
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return 0, err
	}
	// Check Response Code
	if responseCode != http.StatusCreated {
		return 0, fmt.Errorf("error posting new expense:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	// Extrac created Expense expense
	respObj := struct{ ID domain.ExpenseID }{}
	if err := json.Unmarshal(responseBody, &respObj); err != nil {
		return 0, err
	}
	return respObj.ID, nil
}

// UpdateExpense sends a PUT request to update given expense
func UpdateExpense(payload server.JSONReqExpense, token string) error {
	// Marshall Expense to JSON
	//log.Println(payload)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	// Send HTTP Request
	path := url + "/expenses/" + strconv.Itoa(int(payload.ID))
	requestBody := bytes.NewBuffer(jsonPayload)
	req, err := http.NewRequest("PUT", path, requestBody)
	if err != nil {
		return err
	}
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return err
	}
	// Check Response Code
	if responseCode != http.StatusOK {
		return fmt.Errorf("error updating expense:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	return nil
}

// FetchExpenses sends a GET request to fetch all expenses
func FetchExpenses(token string, minTime time.Time) ([]server.JSONRespListExpense, error) {
	// Send HTTP Request
	path := url + "/expenses?from=" + minTime.Format("01-02-2006")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return []server.JSONRespListExpense{}, err
	}
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []server.JSONRespListExpense{}, err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return []server.JSONRespListExpense{}, err
	}
	// Check Response Code
	if responseCode != http.StatusOK {
		return []server.JSONRespListExpense{}, fmt.Errorf("error fetching expenses:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	// Extract Expenses
	var expenses []server.JSONRespListExpense
	if err := json.Unmarshal(responseBody, &expenses); err != nil {
		return []server.JSONRespListExpense{}, err
	}
	return expenses, nil
}

// FetchExpenseDetails sends a GET request to fetch expense with given id
func FetchExpenseDetails(id domain.ExpenseID, token string) (server.JSONRespDetailExpense, error) {
	// Send HTTP Request
	path := url + "/expenses/" + strconv.Itoa(int(id))
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return server.JSONRespDetailExpense{}, err
	}
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return server.JSONRespDetailExpense{}, err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return server.JSONRespDetailExpense{}, err
	}
	// Check Response Code
	if responseCode != http.StatusOK {
		return server.JSONRespDetailExpense{}, fmt.Errorf("error fetching expenses:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	// Unmarshall Expense
	var expense server.JSONRespDetailExpense
	if err := json.Unmarshal(responseBody, &expense); err != nil {
		return server.JSONRespDetailExpense{}, err
	}
	return expense, nil
}

// DeleteExpense sends a POST request with refresh token and
// gets a new Jwt Access Token
func DeleteExpense(id domain.ExpenseID, token string) error {
	// Send HTTP Request
	path := url + "/expenses/" + strconv.Itoa(int(id))
	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return err
	}
	// Check Response Code
	if responseCode != http.StatusNoContent {
		return fmt.Errorf("error deleting expense:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	return nil
}
