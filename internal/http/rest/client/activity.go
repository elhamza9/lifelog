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

type postActivityRespPayload struct {
	ID int `json:"id"`
}

// PostActivity sends a POST request with refresh token and
// gets a new Jwt Access Token
func PostActivity(payload server.JSONReqActivity, token string) (int, error) {
	// Marshall Activity to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	// Send HTTP Request
	const path string = url + "/activities"
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
		return 0, fmt.Errorf("error posting new activity:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	// Extract ID created Activity
	respObj := postActivityRespPayload{}
	if err := json.Unmarshal(responseBody, &respObj); err != nil {
		return 0, err
	}
	return respObj.ID, nil
}

// UpdateActivity sends a PUT request to update given activity
func UpdateActivity(payload server.JSONReqActivity, token string) error {
	// Marshall Activity to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	// Send HTTP Request
	path := url + "/activities/" + strconv.Itoa(int(payload.ID))
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
		return fmt.Errorf("error updating activity:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	return nil
}

// FetchActivities sends a GET request to fetch all activities
func FetchActivities(token string, minTime time.Time) ([]server.JSONRespListActivity, error) {
	// Send HTTP Request
	path := url + "/activities?from=" + minTime.Format("01-02-2006")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return []server.JSONRespListActivity{}, err
	}
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []server.JSONRespListActivity{}, err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return []server.JSONRespListActivity{}, err
	}
	// Check Response Code
	if responseCode != http.StatusOK {
		return []server.JSONRespListActivity{}, fmt.Errorf("error fetching activities:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	// Extract Activities
	var activities []server.JSONRespListActivity
	if err := json.Unmarshal(responseBody, &activities); err != nil {
		return []server.JSONRespListActivity{}, err
	}
	return activities, nil
}

// FetchActivityDetails sends a GET request to fetch activity with given id
func FetchActivityDetails(id domain.ActivityID, token string) (server.JSONRespDetailActivity, error) {
	// Send HTTP Request
	path := url + "/activities/" + strconv.Itoa(int(id))
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return server.JSONRespDetailActivity{}, err
	}
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return server.JSONRespDetailActivity{}, err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return server.JSONRespDetailActivity{}, err
	}
	// Check Response Code
	if responseCode != http.StatusOK {
		return server.JSONRespDetailActivity{}, fmt.Errorf("error fetching activity details:\n\t- code: %d\n\t- body: %s", responseCode, responseBody)
	}
	// Extract Activity
	var activity server.JSONRespDetailActivity
	if err := json.Unmarshal(responseBody, &activity); err != nil {
		return server.JSONRespDetailActivity{}, err
	}
	return activity, nil
}

// DeleteActivity sends a POST request with refresh token and
// gets a new Jwt Access Token
func DeleteActivity(id domain.ActivityID, token string) error {
	// Send HTTP Request
	path := url + "/activities/" + strconv.Itoa(int(id))
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
