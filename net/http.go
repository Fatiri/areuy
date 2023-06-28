package net

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"net/http"
	"time"
)

func (hc HTTPClientCtx) Invoke(param *ParamaterHttpClient) (*http.Response, error) {
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

	response, err := HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (hc HTTPClientCtx) InvokeResty(param *ParamaterHttpClient) (*resty.Response, error) {
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

func (hc HTTPClientCtx) ReadHttpResponse(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
