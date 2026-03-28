package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
)

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

		if strings.Contains(target, "{id}") {
			parts := strings.Split(path, "/")
			for _, part := range parts {
				if part != "" && part != "orders" && part != "menu" {
					targetURL = strings.Replace(target, "{id}", part, 1)
				}
			}
		}

		fullURL := targetURL + path
		if r.URL.RawQuery != "" {
			fullURL += "?" + r.URL.RawQuery
		}

		log.Printf("Proxying %s %s -> %s", r.Method, r.URL.Path, fullURL)

		var body io.Reader
		if r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			body = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequest(r.Method, fullURL, body)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}

		for key, values := range r.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Proxy error: %v", err)
			http.Error(w, "Error forwarding request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
