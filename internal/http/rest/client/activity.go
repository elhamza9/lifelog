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

type postActivityReqPayload struct {
	Label    string        `json:"label"`
	Place    string        `json:"place"`
	Desc     string        `json:"desc"`
	Time     time.Time     `json:"time"`
	Duration time.Duration `json:"duration"`
	Tags     []domain.Tag  `json:"tags"`
}

type postActivityRespPayload struct {
	ID int `json:"id"`
}

// PostActivity sends a POST request with refresh token and
// gets a new Jwt Access Token
func PostActivity(act domain.Activity, token string) (int, error) {
	// Marshall Activity to JSON
	payload := postActivityReqPayload{
		Label:    act.Label,
		Place:    act.Place,
		Desc:     act.Desc,
		Time:     act.Time,
		Duration: act.Duration,
		Tags:     act.Tags,
	}
	log.Println(payload)
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
		return 0, fmt.Errorf("error posting new activity:\n\t- code: %d\n\t- body: %s\n", responseCode, responseBody)
	}
	// Extract ID created Activity
	respObj := postActivityRespPayload{}
	if err := json.Unmarshal(responseBody, &respObj); err != nil {
		return 0, err
	}
	return respObj.ID, nil
}

// FetchActivities sends a GET request to fetch all activities
func FetchActivities(token string, minTime time.Time) ([]domain.Activity, error) {
	// Send HTTP Request
	path := url + "/activities?from=" + minTime.Format("01-02-2006")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return []domain.Activity{}, err
	}
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []domain.Activity{}, err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return []domain.Activity{}, err
	}
	// Check Response Code
	if responseCode != http.StatusOK {
		return []domain.Activity{}, fmt.Errorf("error fetching activities:\n\t- code: %d\n\t- body: %s\n", responseCode, responseBody)
	}
	// Extract Activities
	var activities []domain.Activity
	if err := json.Unmarshal(responseBody, &activities); err != nil {
		return []domain.Activity{}, err
	}
	return activities, nil
}
