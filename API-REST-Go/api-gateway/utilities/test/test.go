package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type RequestParams struct {
	Method  string
	Path    string
	Body    any
	Headers map[string]string
	Cookies []*http.Cookie
}

func NewRequest(r *RequestParams) *http.Request {
	var body bytes.Buffer
	_ = json.NewEncoder(&body).Encode(r.Body)
	req := httptest.NewRequest(r.Method, r.Path, &body)
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	for _, cookie := range r.Cookies {
		req.AddCookie(cookie)
	}
	return req
}
