package main

import "net/http"

func LeakyBucketMiddleware(lb *LeakyBucket) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if lb.TryTake() {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
			}
		})
	}
}
