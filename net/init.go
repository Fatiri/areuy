package net

import (
	"net/http"
	"time"
)

var (
	HTTPClient IInitialHTTPClient
)

func init() {
	HTTPClient = &http.Client{
		Timeout:   time.Second * 10,
	}
}

type IInitialHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
