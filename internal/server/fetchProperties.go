package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/time/rate"
)

// Global rate limiter
var limiter = rate.NewLimiter(rate.Every(5*time.Millisecond), 20) // 200 requests per second

const (
	maxConcurrentRequests = 100 // Increased from 5
	maxIdleConns          = 400 // Increased from 100
	requestTimeout        = 30 * time.Second
)

func GetAllProperties() {
	// Create a transport with connection pooling
	transport := &http.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConns,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Add timeout
	}

	resp, err := http.Get("https://www.centris.ca")
	if err != nil {
		log.Fatalf("Failed to make GET request: %v", err)
	}

	cookies := resp.Header["Set-Cookie"]
	var aspNetCoreSession, arrAffinitySameSite string

	for _, cookie := range cookies {
		if strings.Contains(cookie, ".AspNetCore.Session=") {
			aspNetCoreSession = extractCookieValue(cookie, ".AspNetCore.Session")
		} else if strings.Contains(cookie, "ARRAffinitySameSite=") {
			arrAffinitySameSite = extractCookieValue(cookie, "ARRAffinitySameSite")
		}
	}
	resp.Body.Close()

	pins := getAllPins(client, aspNetCoreSession, arrAffinitySameSite)
	housesHTML := getAllHouses(client, pins, aspNetCoreSession, arrAffinitySameSite)

	housesJSON, err := json.MarshalIndent(housesHTML, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling markers to JSON: %v", err)
	}

	fileName := "housesHTML.json"
	err = os.WriteFile(fileName, housesJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Printf("Markers JSON saved to file: %s\n", fileName)
}

func getAllPins(client *http.Client, aspNetCoreSession string, arrAffinitySameSite string) []Marker {
	bodyData := InputData{
		ZoomLevel: 11,
		MapBounds: MapBounds{
			NorthEast: Coordinate{
				Lat: 51.41553513240069,
				Lng: -57.19362267293036,
			},
			SouthWest: Coordinate{
				Lat: 44.99651987571269,
				Lng: -79.53598220832646,
			},
		},
	}

	markerJSON, _ := json.Marshal(bodyData)
	markerReq, _ := http.NewRequest("POST", "https://www.centris.ca/api/property/map/GetMarkers",
		bytes.NewBuffer(markerJSON))
	markerReq.Header.Set("Content-Type", "application/json")
	markerReq.Header.Set("Cookie", ".AspNetCore.Session="+aspNetCoreSession+"; ARRAffinitySameSite="+arrAffinitySameSite+";")

	resp, err := client.Do(markerReq)
	if err != nil {
		fmt.Printf("Error fetching markers: %v\n", err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status: %d %s\n", resp.StatusCode, resp.Status)
		resp.Body.Close()
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshaling response: %v", err)
	}

	markers := response.D.Result.Markers

	return markers
}

func getAllHouses(client *http.Client, pins []Marker, aspNetCoreSession string, arrAffinitySameSite string) []string {
	var houseLinks []string
	var mu sync.Mutex // Mutex for thread-safe append

	// Use semaphore to limit concurrent goroutines
	sem := make(chan struct{}, maxConcurrentRequests) // Limit to 5 concurrent requests

	var wg sync.WaitGroup

	for _, pin := range pins {
		for i := 0; i < pin.PointsCount; i++ {
			wg.Add(1)
			go func(pin Marker, idx int) {
				defer wg.Done()

				// Acquire semaphore
				sem <- struct{}{}
				defer func() { <-sem }()

				// Wait for rate limiter
				err := limiter.Wait(context.Background())
				if err != nil {
					log.Printf("Rate limiter error: %v\n", err)
					return
				}

				link := GetAllHouseInPin(client, aspNetCoreSession, arrAffinitySameSite, pin, idx)
				fmt.Println(link)
				if link != "" {
					mu.Lock()
					houseLinks = append(houseLinks, link)
					mu.Unlock()
				}
			}(pin, i)
		}
	}

	wg.Wait()
	return houseLinks
}

func GetAllHouseInPin(client *http.Client, aspNetCoreSession string, arrAffinitySameSite string, pin Marker, idx int) string {
	bodyData := MarkerInfoInputData{
		PageIndex: idx,
		ZoomLevel: 11,
		Latitude:  pin.Position.Lat,
		Longitude: pin.Position.Lng,
		GeoHash:   pin.GeoHash,
	}

	pinJSON, err := json.Marshal(bodyData)
	if err != nil {
		log.Printf("Error marshaling body data: %v\n", err)
		return ""
	}

	pinReq, err := http.NewRequest("POST", "https://www.centris.ca/property/GetMarkerInfo", bytes.NewBuffer(pinJSON))
	if err != nil {
		log.Printf("Error creating HTTP request: %v\n", err)
		return ""
	}

	pinReq.Header.Set("Content-Type", "application/json")
	pinReq.Header.Set("Cookie", ".AspNetCore.Session="+aspNetCoreSession+"; ARRAffinitySameSite="+arrAffinitySameSite+";")

	// Add retries with exponential backoff
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		resp, err := client.Do(pinReq)
		if err != nil {
			log.Printf("Error fetching markers (attempt %d/%d): %v\n", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(time.Duration(1<<uint(i)) * time.Second)
				continue
			}
			return ""
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error status (attempt %d/%d): %d %s\n", i+1, maxRetries, resp.StatusCode, resp.Status)
			if i < maxRetries-1 {
				time.Sleep(time.Duration(1<<uint(i)) * time.Second)
				continue
			}
			return ""
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v\n", err)
			return ""
		}

		var response MarkerInfoResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Printf("Error unmarshaling response: %v\n", err)
			return ""
		}

		houseHTML, err := html.Parse(strings.NewReader(response.D.Result.Html))
		if err != nil {
			log.Printf("Error parsing HTML: %v\n", err)
			return ""
		}

		return findElementAttribute(houseHTML, "a", "class", "property-thumbnail-summary-link", "href")
	}

	return ""
}
