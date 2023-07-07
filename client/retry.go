package main

import (
	"net/http"
)

type RetryRoundTripper struct {
	RoundTripper http.RoundTripper
	RetryCount   int
}

func (rrt *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := rrt.RoundTripper
	if rt == nil {
		rt = http.DefaultTransport
	}

	var resp *http.Response
	var err error
	for i := 0; i <= rrt.RetryCount; i++ {
		resp, err = rt.RoundTrip(req)
		if err != nil || (resp != nil && resp.StatusCode >= http.StatusInternalServerError || resp.StatusCode == http.StatusTooManyRequests) {
			continue
		}
		return resp, nil
	}
	return resp, err
}
