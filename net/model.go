package net

import "net/url"

type ParamaterHttpClient struct {
	URL               string           `json:"url"`
	Method            string           `json:"method"`
	ContentType       string           `json:"content_type"`
	KeyAuthorization  string           `json:"key_authorization"`
	Authorization     string           `json:"authorization"`
	WithAuthorization bool             `json:"with_authorization"`
	BodyRequest       []byte          `json:"body_request"`
	Headers           []RequestHttpClient `json:"headers"`
	Query             []RequestHttpClient `json:"query"`
	UrlValue          url.Values
}

type RequestHttpClient struct {
	Key   string
	Value string
}
