package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elhamza90/lifelog/internal/domain"
)

type postTagReqPayload struct {
	Name string `json:"name"`
}

type postTagRespPayload struct {
	ID int `json:"id"`
}

// PostTag sends a POST request with refresh token and
// gets a new Jwt Access Token
func PostTag(tag domain.Tag, token string) (int, error) {
	// Marshall Tag to JSON
	payload := postTagReqPayload{Name: tag.Name}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	// Send HTTP Request
	const path string = url + "/tags"
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
		return 0, fmt.Errorf("error posting new tag:\n\t- code: %d\n\t- body: %s\n", responseCode, responseBody)
	}
	// Extract ID created Tag
	respObj := postTagRespPayload{}
	if err := json.Unmarshal(responseBody, &respObj); err != nil {
		return 0, err
	}
	return respObj.ID, nil
}
