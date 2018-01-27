package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type fixture struct {
	header     string
	hostURL    string
	method     string
	reqURL     string
	statusCode int
}

func TestHandler(t *testing.T) {
	testTable := make(map[string]fixture)

	testTable["TestHTTPSRedirect"] = fixture{
		statusCode: http.StatusMovedPermanently,
		reqURL:     "http://src.danbond.io/vanity",
	}
	testTable["TestProtoHeaderHTTPSRedirect"] = fixture{
		header:     "http",
		reqURL:     "http://src.danbond.io/vanity",
		statusCode: http.StatusMovedPermanently,
	}
	testTable["TestMethodNotAllowed"] = fixture{
		method:     http.MethodPost,
		reqURL:     "https://src.danbond.io/vanity",
		statusCode: http.StatusMethodNotAllowed,
	}
	testTable["TestRedirectURLError"] = fixture{
		hostURL:    "invalid host",
		method:     http.MethodGet,
		reqURL:     "https://src.danbond.io/vanity",
		statusCode: http.StatusInternalServerError,
	}
	testTable["TestNoGoGet"] = fixture{
		method:     http.MethodGet,
		reqURL:     "https://src.danbond.io/vanity",
		statusCode: http.StatusTemporaryRedirect,
	}
	testTable["TestNoPath"] = fixture{
		method:     http.MethodGet,
		reqURL:     "https://src.danbond.io?go-get=1",
		statusCode: http.StatusTemporaryRedirect,
	}
	testTable["TestSuccess"] = fixture{
		method:     http.MethodGet,
		reqURL:     "https://src.danbond.io/vanity?go-get=1",
		statusCode: http.StatusOK,
	}

	for name, test := range testTable {
		t.Run(name, func(t *testing.T) {
			u, err := url.Parse(test.reqURL)
			if err != nil {
				t.Fatal(err)
			}
			r := &http.Request{
				Header: make(http.Header),
				Method: test.method,
				URL:    u,
			}
			if test.header != "" {
				r.Header.Add(protoHeader, test.header)
			}
			u, err = url.Parse("https://github.com/danbondd")
			if err != nil {
				t.Fatal(err)
			}
			if test.hostURL != "" {
				u.Host = test.hostURL
			}

			h := handler("git", u)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, r)

			if rec.Code != test.statusCode {
				t.Errorf("expected status code: %d, got: %d", test.statusCode, rec.Code)
			}
		})
	}
}
