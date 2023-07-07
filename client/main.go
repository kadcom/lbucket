package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	count   int
	method  string
	retries int
)

func init() {
	flag.IntVar(&count, "c", 1, "Number of requests to make")
	flag.StringVar(&method, "m", "none", "Retry method to use: none, retry, backoff, breaker")
	flag.IntVar(&retries, "r", 3, "Maximum number of retries for failed requests")
}

func main() {
	flag.Parse()

	urls := flag.Args()
	if len(urls) == 0 {
		fmt.Println("You must provide at least one URL")
		os.Exit(1)
	}

	client := &http.Client{
		Transport: http.DefaultTransport,
	}

	switch method {
	case "retry":
		client.Transport = &RetryRoundTripper{
			RoundTripper: http.DefaultTransport,
			RetryCount:   retries,
		}
	case "backoff":
		client.Transport = &ExponentialBackoffRoundTripper{
			RoundTripper: http.DefaultTransport,
			MaxRetries:   retries,
			BackoffTime:  1 * time.Second,
		}
	case "breaker":
		client.Transport = &CircuitBreakerRoundTripper{
			RoundTripper: http.DefaultTransport,
			FailureCount: 0,
			Open:         false,
			ResetTimeout: 1 * time.Second,
		}
	}

	for i := 0; i < count; i++ {
		for _, url := range urls {
			resp, err := client.Get(url)
			if err != nil {
				fmt.Println("Request error:", err)
				os.Exit(1)
			}
			defer resp.Body.Close()

			fmt.Println("Response status:", resp.Status)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading body:", err)
				os.Exit(1)
			}

			fmt.Println("Response body:", string(body))
		}
	}
}
