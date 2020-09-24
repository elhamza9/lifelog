package client

import (
	"io"
	"io/ioutil"

	"github.com/elhamza90/lifelog/internal/domain"
)

func getIdsFromTags(tags []domain.Tag) []domain.TagID {
	res := []domain.TagID{}
	for _, t := range tags {
		res = append(res, t.ID)
	}
	return res
}

func readResponseBody(respBody io.ReadCloser) ([]byte, error) {
	responseBody, err := ioutil.ReadAll(respBody)
	if err != nil {
		return []byte{}, err
	}
	respBody.Close()
	return responseBody, nil
}
