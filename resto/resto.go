package resto

import (
	"bytes"
	"net/http"
	"time"

	json "github.com/goccy/go-json"
	"github.com/hashicorp/go-retryablehttp"
)

type Resto struct {
	retry *retryablehttp.Client
}

/*
Tạo Retry HTTP Client với tham số
tryTimes: số lần cố gắng thử lại
waitMax: thời gian đợi tối đa sau mỗi lần tính bằng millisecond
*/
func Retry(tryTimes int, waitMax time.Duration) Resto {
	retryHTTP := retryablehttp.NewClient()
	retryHTTP.RetryMax = tryTimes
	retryHTTP.RetryWaitMax = waitMax * time.Millisecond
	retryHTTP.Logger = nil
	return Resto{
		retry: retryHTTP,
	}
}

func (r Resto) Post(url string, data interface{}) (*http.Response, error) {
	jsonstr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	response, err := r.retry.Post(url, "application/json", bytes.NewBuffer(jsonstr))
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r Resto) Get(url string) (*http.Response, error) {
	response, err := r.retry.Get(url)
	if err != nil {
		return nil, err
	}
	return response, nil
}
