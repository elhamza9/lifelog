package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const host string = "localhost"
const port string = "8080"
const url string = "http://" + host + ":" + port

// loginReqPayload is used to marshall credentials to JSON when authenticating
type loginReqPayload struct {
	Password string `json:"password"`
}

// loginRespPayload is used to unmarshall JSON Response
// when authentication is successful
type loginRespPayload struct {
	Access  string `json:"at"`
	Refresh string `json:"rt"`
}

// Login sends a POST request with password and
// gets a Jwt Token Pair if password is correct
func Login(pass string) (string, string, error) {
	// Marshall Tag to JSON
	payload := loginReqPayload{Password: pass}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}
	// Send HTTP Request
	const path string = url + "/auth/login"
	requestBody := bytes.NewBuffer(jsonPayload)
	resp, err := http.Post(path, "application/json", requestBody)
	if err != nil {
		return "", "", err
	}
	// Read Response
	responseCode := resp.StatusCode
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return "", "", err
	}
	// Check Response Code
	if responseCode != http.StatusOK {
		return "", "", fmt.Errorf("error authenticating:\n\t- code: %d\n\t- body: %s\n", responseCode, responseBody)
	}
	// Extract Tokens
	respObj := loginRespPayload{}
	if err := json.Unmarshal(responseBody, &respObj); err != nil {
		return "", "", err
	}
	return respObj.Access, respObj.Refresh, nil
}
