package net

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

type IHTTPClient interface {
	Invoke(param *ParamaterHttpClient) (*http.Response, error)
	InvokeResty(param *ParamaterHttpClient) (*resty.Response, error)
	ReadHttpResponse(response *http.Response) ([]byte, error)
}

type HTTPClientCtx struct {}

func ProvideIHTTPClient() IHTTPClient {
	return &HTTPClientCtx{}
}
