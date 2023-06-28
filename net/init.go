package net

import (
	"net"
	"net/http"
	"time"
)

var (
	HTTPClient IInitialHTTPClient
)

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

func init() {
	HTTPClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

type IInitialHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
