# AREUY Times

# Instatlations
```sh
$ go get github.com/Fatiri/areuy/net
```

## List Of interfaces of HTTPClient

```sh
type IHTTPClient interface {
	Invoke(param *ParamaterHttpClient) (*http.Response, error)
	InvokeResty(param *ParamaterHttpClient) (*resty.Response, error)
	ReadHttpResponse(response *http.Response) ([]byte, error)
}
```

## How To test

   Run init func in specific test file for and mocking the HTTP client

```
func init() {
	client.Client = &mocks.MockClient{}
}
```

  and create mocking file for assign d=the HTTPCLient

```
```