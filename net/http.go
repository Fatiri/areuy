package net

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// GetHTTPRequestJSON ...
func GetHTTPRequestJSON(ctx context.Context, method string, url string, body io.Reader, headers ...map[string]string) (res []byte, statusCode int, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// iterate optional data of headers
	for _, header := range headers {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	timeout, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT"))
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	r, err := client.Do(req)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	resp := StreamToByte(r.Body)

	if os.Getenv("DEVELOPMENT") == "1" {
		// fmt.Println("code : " + r.Status)
		// fmt.Println("url : " + url)
		// fmt.Println("mtd : " + method)
		// fmt.Printf("%v\n", headers)
		// fmt.Println("resp :", string(resp))
	}

	defer func() {
		r.Body.Close()
	}()

	return resp, r.StatusCode, nil
}

// StreamToString func
func StreamToString(stream io.Reader) string {
	if stream == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}

// StreamToByte ...
func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
