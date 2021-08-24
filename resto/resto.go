package resto

import (
	"bytes"
	"net/http"
	"time"

	json "github.com/goccy/go-json"
	"github.com/hashicorp/go-retryablehttp"
)

func Post(url string, data interface{}) (*http.Response, error) {
	jsonstr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	restclient := retryablehttp.NewClient()
	restclient.RetryMax = 4
	restclient.RetryWaitMax = 100 * time.Millisecond
	restclient.Logger = nil

	response, err := restclient.Post(url, "application/json", bytes.NewBuffer(jsonstr))
	if err != nil {
		return nil, err
	}
	return response, nil
}

func Get(url string) (*http.Response, error) {
	restclient := retryablehttp.NewClient()
	restclient.RetryMax = 4
	restclient.RetryWaitMax = 100 * time.Millisecond
	restclient.Logger = nil

	response, err := restclient.Get(url)
	if err != nil {
		return nil, err
	}
	return response, nil
}
