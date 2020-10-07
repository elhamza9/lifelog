package client

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/elhamza90/lifelog/internal/domain"
)

var (
	host string = os.Getenv("LFLG_HOST")
	port string = os.Getenv("LFLG_PORT")
	url  string = "http://" + host + ":" + port
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
