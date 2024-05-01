package model

import (
	"fmt"
	"net/url"
	"strings"
)

func matchURL(inputURL, domain string) bool {
	// Parse the input URL
	u, err := url.Parse(inputURL)
	if err != nil {
		// Handle invalid URLs
		return false
	}

	// Extract the host from the parsed URL
	host := u.Hostname()
	fmt.Println("host: ", host, "domain: ", domain, "inputURL: ", inputURL)

	// Check if the host contains the desired domain
	return strings.Contains(domain, host)
}
