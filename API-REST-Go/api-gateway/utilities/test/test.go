package test

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/goccy/go-json"
)

type RequestParams struct {
	Method      string
	Path        string
	ContentType string
	Body        any
	Headers     map[string]string
	Cookies     []*http.Cookie
}

func NewRequest(r *RequestParams) *http.Request {
	if r.ContentType == "" {
		r.ContentType = "application/json"
	}

	var body bytes.Buffer
	json.NewEncoder(&body).Encode(r.Body)

	req := httptest.NewRequest(r.Method, r.Path, &body)
	req.Header.Add("Content-Type", r.ContentType)
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	for _, cookie := range r.Cookies {
		req.AddCookie(cookie)
	}
	return req
}

func BodyToString(body io.ReadCloser) string {
	// read response body
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}
	// close response body
	body.Close()

	return string(bytes)
}
