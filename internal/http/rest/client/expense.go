package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
)

type postExpenseReqPayload struct {
	Label string       `json:"label"`
	Time  time.Time    `json:"time"`
	Value float32      `json:"value"`
	Unit  string       `json:"unit"`
	Tags  []domain.Tag `json:"tags"`
}

type postExpenseRespPayload struct {
	ID int `json:"id"`
}

// PostExpense sends a POST request with refresh token and
// gets a new Jwt Access Token
func PostExpense(exp domain.Expense, token string) (int, error) {
	// Marshall Expense to JSON
	payload := postExpenseReqPayload{
		Label: exp.Label,
		Value: exp.Value,
		Unit:  exp.Unit,
		Time:  exp.Time,
		Tags:  exp.Tags,
	}
	log.Println(payload)
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
