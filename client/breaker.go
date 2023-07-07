package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CircuitBreakerRoundTripper struct {
	sync.Mutex
	RoundTripper http.RoundTripper
	FailureCount int
	Open         bool
	ResetTimeout time.Duration
}

func (cb *CircuitBreakerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	cb.Lock()
	defer cb.Unlock()

	rt := cb.RoundTripper
	if rt == nil {
		rt = http.DefaultTransport
	}

	if cb.Open {
		return nil, fmt.Errorf("circuit breaker is open")
	}

	resp, err := rt.RoundTrip(req)
	if err != nil || (resp != nil && (resp.StatusCode >= http.StatusInternalServerError || resp.StatusCode == http.StatusTooManyRequests)) {
		cb.FailureCount++
		if cb.FailureCount >= 3 && !cb.Open {
			cb.Open = true
			time.AfterFunc(cb.ResetTimeout, func() {
				cb.Lock()
				defer cb.Unlock()
				cb.Open = false
				cb.FailureCount = 0
			})
		}
		return resp, err
	}

	cb.FailureCount = 0
	return resp, nil
}
