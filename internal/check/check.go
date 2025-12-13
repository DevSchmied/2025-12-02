package check

import (
	"net/http"
	"strings"
	"time"
)

func CheckLink(url string) bool {
	// Check whether the URL has an http/https prefix. This is required for further processing.
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Create an HTTP client with a 2-second timeout.
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	// Check whether the site responds to the request.
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	// Close the response body to avoid resource leaks.
	defer resp.Body.Close()

	return resp.StatusCode < 400
}
