package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
)

type expenseReqPayload struct {
	Label      string            `json:"label"`
	Time       time.Time         `json:"time"`
	Value      float32           `json:"value"`
	Unit       string            `json:"unit"`
	ActivityID domain.ActivityID `json:"activityId"`
	Tags       []domain.Tag      `json:"tags"`
}

type postExpenseRespPayload struct {
	ID int `json:"id"`
}

// PostExpense sends a POST request with refresh token and
// gets a new Jwt Access Token
func PostExpense(exp domain.Expense, token string) (int, error) {
	// Marshall Expense to JSON
	payload := expenseReqPayload{
		Label:      exp.Label,
		Value:      exp.Value,
		Unit:       exp.Unit,
		Time:       exp.Time,
		ActivityID: exp.ActivityID,
		Tags:       exp.Tags,
	}
	//log.Println(payload)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	// Send HTTP Request
	const path string = url + "/expenses"
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
	// Extrexp ID created Expense
	respObj := postExpenseRespPayload{}
	if err := json.Unmarshal(responseBody, &respObj); err != nil {
		return 0, err
	}
	return respObj.ID, nil
}

// UpdateExpense sends a PUT request to update given expense
func UpdateExpense(exp domain.Expense, token string) error {
	// Marshall Expense to JSON
	payload := expenseReqPayload{
		Label:      exp.Label,
		Value:      exp.Value,
		Unit:       exp.Unit,
		Time:       exp.Time,
		ActivityID: exp.ActivityID,
		Tags:       exp.Tags,
	}
	//log.Println(payload)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	// Send HTTP Request
	path := url + "/expenses/" + strconv.Itoa(int(exp.ID))
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
func FetchExpenses(token string, minTime time.Time) ([]domain.Expense, error) {
	// Send HTTP Request
	path := url + "/expenses?from=" + minTime.Format("01-02-2006")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return []domain.Expense{}, err
	}
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []domain.Expense{}, err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return []domain.Expense{}, err
	}
	// Check Response Code
	if responseCode != http.StatusOK {
		return []domain.Expense{}, fmt.Errorf("error fetching expenses:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	// Extract Expenses
	var expenses []domain.Expense
	if err := json.Unmarshal(responseBody, &expenses); err != nil {
		return []domain.Expense{}, err
	}
	return expenses, nil
}

// FetchExpenseDetails sends a GET request to fetch expense with given id
func FetchExpenseDetails(id domain.ExpenseID, token string) (domain.Expense, error) {
	// Send HTTP Request
	path := url + "/expenses/" + strconv.Itoa(int(id))
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return domain.Expense{}, err
	}
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return domain.Expense{}, err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return domain.Expense{}, err
	}
	// Check Response Code
	if responseCode != http.StatusOK {
		return domain.Expense{}, fmt.Errorf("error fetching expenses:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	// Extract Expenses
	var expense domain.Expense
	if err := json.Unmarshal(responseBody, &expense); err != nil {
		return domain.Expense{}, err
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
