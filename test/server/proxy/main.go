package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// Define the target URL where you want to proxy requests
	targetUrl, err := url.Parse("http://localhost:5000")
	if err != nil {
		log.Fatal("Error parsing target URL:", err)
	}

	// Create a reverse proxy
	reverseProxy := httputil.NewSingleHostReverseProxy(targetUrl)

	// Handle requests by proxying them to the target URL
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Update the request host to match the target URL's host
		r.Host = targetUrl.Host
		// Proxy the request to the target URL
		reverseProxy.ServeHTTP(w, r)
	})

	// Start the server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
