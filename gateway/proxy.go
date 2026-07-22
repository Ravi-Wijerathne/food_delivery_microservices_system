package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sony/gobreaker/v2"
)

var cb *gobreaker.CircuitBreaker[[]byte]

func init() {
	var st gobreaker.Settings
	st.Name = "Gateway_Circuit_Breaker"
	st.MaxRequests = 5
	st.Interval = time.Duration(10) * time.Second
	st.Timeout = time.Duration(5) * time.Second
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}
	cb = gobreaker.NewCircuitBreaker[[]byte](st)
}

func proxyRequest(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		targetURL := target

		if strings.HasPrefix(path, "/api") {
			path = strings.TrimPrefix(path, "/api")
		}

		if strings.HasPrefix(path, "/auth") {
			path = strings.TrimPrefix(path, "/auth")
		}

		if strings.Contains(targetURL, "{id}") {
			parts := strings.Split(path, "/")
			for _, part := range parts {
				if part != "" && part != "orders" && part != "menu" && part != "users" && part != "profile" {
					targetURL = strings.Replace(targetURL, "{id}", part, 1)
				}
			}
		}

		fullURL := targetURL + path
		if r.URL.RawQuery != "" {
			fullURL += "?" + r.URL.RawQuery
		}

		log.Printf("Proxying %s %s -> %s", r.Method, r.URL.Path, fullURL)

		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(r.Body)
		}

		// Use Circuit Breaker
		responseBytes, err := cb.Execute(func() ([]byte, error) {
			req, err := http.NewRequest(r.Method, fullURL, bytes.NewReader(body))
			if err != nil {
				return nil, err
			}

			for key, values := range r.Header {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}

			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.WriteHeader(resp.StatusCode)
			
			respBody, _ := io.ReadAll(resp.Body)
			return respBody, nil
		})

		if err != nil {
			log.Printf("Proxy error (CircuitBreaker): %v", err)
			http.Error(w, "Service Unavailable (Circuit Open)", http.StatusServiceUnavailable)
			return
		}

		w.Write(responseBytes)
	}
}
