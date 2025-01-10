package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func GenerateSession(url string) (string, string, int) {
	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Extract cookies from the response headers
	cookies := resp.Header["Set-Cookie"]

	// Variables to store the desired cookies
	var aspNetCoreSession string
	var arrAffinitySameSite string

	// Iterate through cookies to find the ones we need
	for _, cookie := range cookies {
		if strings.Contains(cookie, ".AspNetCore.Session=") {
			aspNetCoreSession = ExtractCookieValue(cookie, ".AspNetCore.Session")
		} else if strings.Contains(cookie, "ARRAffinitySameSite=") {
			arrAffinitySameSite = ExtractCookieValue(cookie, "ARRAffinitySameSite")
		}
	}

	numberOfBrokers := 0
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return aspNetCoreSession, arrAffinitySameSite, numberOfBrokers
	}

	numberOfBrokers, _ = strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(FindElementByClass(doc, "span", "resultCount"), "\u00a0", ""), " ", ""))

	return aspNetCoreSession, arrAffinitySameSite, numberOfBrokers
}

func ExtractCookieValue(cookie string, key string) string {
	parts := strings.Split(cookie, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, key+"=") {
			return strings.TrimPrefix(part, key+"=")
		}
	}
	return ""
}
