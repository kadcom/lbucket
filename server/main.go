package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	lb := CreateLeakyBucket(10, time.Second)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(LeakyBucketMiddleware(&lb))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	http.ListenAndServe(":3000", r)
}
