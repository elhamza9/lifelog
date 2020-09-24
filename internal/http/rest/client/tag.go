package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/elhamza90/lifelog/internal/domain"
)

type tagReqPayload struct {
	Name string `json:"name"`
}

type postTagRespPayload struct {
	ID int `json:"id"`
}

// PostTag sends a POST request with refresh token and
// gets a new Jwt Access Token
func PostTag(tag domain.Tag, token string) (int, error) {
	// Marshall Tag to JSON
	payload := tagReqPayload{Name: tag.Name}
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
		return 0, fmt.Errorf("error posting new tag:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	// Extract ID created Tag
	respObj := postTagRespPayload{}
	if err := json.Unmarshal(responseBody, &respObj); err != nil {
		return 0, err
	}
	return respObj.ID, nil
}

// UpdateTag sends a POST request with refresh token and
// gets a new Jwt Access Token
func UpdateTag(tag domain.Tag, token string) error {
	// Marshall Tag to JSON
	payload := tagReqPayload{Name: tag.Name}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	// Send HTTP Request
	path := url + "/tags/" + strconv.Itoa(int(tag.ID))
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
		return fmt.Errorf("error updating tag:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	return nil
}

// DeleteTag sends a POST request with refresh token and
// gets a new Jwt Access Token
func DeleteTag(id domain.TagID, token string) error {
	// Send HTTP Request
	path := url + "/tags/" + strconv.Itoa(int(id))
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
		return fmt.Errorf("error deleting tag:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	return nil
}

// FetchTags sends a GET request to fetch all tags
func FetchTags(token string) ([]domain.Tag, error) {
	// Send HTTP Request
	const path string = url + "/tags"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return []domain.Tag{}, err
	}
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []domain.Tag{}, err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return []domain.Tag{}, err
	}
	// Check Response Code
	if responseCode != http.StatusOK {
		return []domain.Tag{}, fmt.Errorf("error posting new tag:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	// Extract Tags
	var tags []domain.Tag
	if err := json.Unmarshal(responseBody, &tags); err != nil {
		return []domain.Tag{}, err
	}
	return tags, nil
}
