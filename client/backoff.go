package main

import (
	"net/http"
	"time"
)

type ExponentialBackoffRoundTripper struct {
	RoundTripper http.RoundTripper
	MaxRetries   int
	BackoffTime  time.Duration
}

func (eb *ExponentialBackoffRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := eb.RoundTripper
	if rt == nil {
		rt = http.DefaultTransport
	}

	var resp *http.Response
	var err error
	backoff := eb.BackoffTime
	for i := 0; i <= eb.MaxRetries; i++ {
		resp, err = rt.RoundTrip(req)
		if err != nil || (resp != nil && (resp.StatusCode >= http.StatusInternalServerError || resp.StatusCode == http.StatusTooManyRequests)) {
			time.Sleep(backoff)
			backoff *= 2
			continue
		}
		return resp, nil
	}
	return resp, err
}
