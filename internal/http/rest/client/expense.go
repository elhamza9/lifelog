package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
)

type postExpenseReqPayload struct {
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
	payload := postExpenseReqPayload{
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
		return 0, fmt.Errorf("error posting new expense:\n\t- code: %d\n\t- body: %s\n", responseCode, responseBody)
	}
	// Extrexp ID created Expense
	respObj := postExpenseRespPayload{}
	if err := json.Unmarshal(responseBody, &respObj); err != nil {
		return 0, err
	}
	return respObj.ID, nil
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
		return []domain.Expense{}, fmt.Errorf("error fetching expenses:\n\t- code: %d\n\t- body: %s\n", responseCode, responseBody)
	}
	// Extract Expenses
	var expenses []domain.Expense
	if err := json.Unmarshal(responseBody, &expenses); err != nil {
		return []domain.Expense{}, err
	}
	return expenses, nil
}
