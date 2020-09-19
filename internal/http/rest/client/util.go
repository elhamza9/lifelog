package client

import (
	"io"
	"io/ioutil"
)

func readResponseBody(respBody io.ReadCloser) ([]byte, error) {
	responseBody, err := ioutil.ReadAll(respBody)
	if err != nil {
		return []byte{}, err
	}
	respBody.Close()
	return responseBody, nil
}
