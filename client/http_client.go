package client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ParamaterHttpClient struct {
	URL               string           `json:"url"`
	Method            string           `json:"method"`
	ContentType       string           `json:"content_type"`
	KeyAuthorization  string           `json:"key_authorization"`
	Authorization     string           `json:"authorization"`
	WithAuthorization bool             `json:"with_authorization"`
	BodyRequest       []byte           `json:"body_request"`
	Headers           []RequestHttpClient `json:"headers"`
	Query             []RequestHttpClient `json:"query"`
	UrlValue          url.Values
}

type RequestHttpClient struct {
	Key   string
	Value string
}

type HTTPClients interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClients
)

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

func init() {
	Client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

func HttpClient(param *ParamaterHttpClient) (*http.Response, error) {
	// byteBody := bytes.NewReader(param.BodyRequest)

	request, err := http.NewRequest(param.Method, param.URL, nil)
	if err != nil {
		return nil, err
	}
	for _, header := range param.Headers {
		request.Header.Set(header.Key, header.Value)
	}

	if len(param.Query) != 0 {
		q := request.URL.Query()
		for _, header := range param.Query {
			q.Add(header.Key, header.Value)
		}
		request.URL.RawQuery = q.Encode()
	}

	if param.WithAuthorization {
		request.Header.Set(param.KeyAuthorization, param.Authorization)
	}

	response, err := Client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func HttpClientV2(param *ParamaterHttpClient, req string) (*http.Response, error) {
	request, err := http.NewRequest(param.Method, param.URL, strings.NewReader(req))
	if err != nil {
		return nil, err
	}

	for _, header := range param.Headers {
		request.Header.Set(header.Key, header.Value)
	}

	if param.WithAuthorization {
		request.Header.Set(param.KeyAuthorization, param.Authorization)
	}

	if len(param.Query) != 0 {
		q := request.URL.Query()
		for _, header := range param.Query {
			q.Add(header.Key, header.Value)
		}
		request.URL.RawQuery = q.Encode()
	}

	response, err := Client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func HTTPClientTokoCrypto(param *ParamaterHttpClient) (*resty.Response, error) {
	client := resty.New()
	response, errAPI := client.SetRetryCount(10).SetRetryWaitTime(5*time.Second).R().
		SetBasicAuth("X-MBX-APIKEY", param.Authorization).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParamsFromValues(param.UrlValue).
		Get(param.URL)

	if errAPI != nil {
		fmt.Println(errAPI)
	}

	return response, nil
}

func ReadHttpResponse(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
