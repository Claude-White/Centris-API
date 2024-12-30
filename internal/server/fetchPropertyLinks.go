package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"time"
)

func main() {
	// Set up channel for handling interrupt signal
	domain := "https://www.centris.ca"
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Create a mutex to protect allData
	var mu sync.Mutex
	var allData []OutputData

	// Create a channel to signal when we're done
	done := make(chan bool)

	go func() {
		fetchPropertyLinks(domain, done, &allData, &mu)
	}()

	// Wait for either interrupt or completion
	select {
	case <-interrupt:
		fmt.Println("Script is about to terminate. Saving data...")
		savePropertyLinks(allData)
	case <-done:
		fmt.Println("Processing complete. Saving data...")
		savePropertyLinks(allData)
	}
}

func savePropertyLinks(arr []OutputData) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(arr); err != nil {
		fmt.Printf("Error creating JSON: %v\n", err)
		return
	}

	if err := os.WriteFile("data.json", buffer.Bytes(), 0644); err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		return
	}
}

func fetchMarkers(domain string, item InputData, client *http.Client) []byte {
	markerPayload := map[string]interface{}{
		"zoomLevel": 11,
		"mapBounds": item.MapBounds,
	}

	markerJSON, _ := json.Marshal(markerPayload)
	markerReq, _ := http.NewRequest("POST", domain+"/api/property/map/GetMarkers",
		bytes.NewBuffer(markerJSON))
	markerReq.Header.Set("Content-Type", "application/json")

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
	return body
}

func fetchMarkerInfo(domain string, marker MarkerPosition, client *http.Client) []byte {
	infoPayload := map[string]interface{}{
		"pageIndex": 0,
		"zoomLevel": 11,
		"latitude":  marker.Position.Lat,
		"longitude": marker.Position.Lng,
		"geoHash":   "f25ds",
	}

	infoJSON, _ := json.Marshal(infoPayload)
	infoReq, _ := http.NewRequest("POST", domain+"/property/GetMarkerInfo",
		bytes.NewBuffer(infoJSON))
	infoReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(infoReq)
	if err != nil {
		fmt.Printf("Error fetching marker info: %v\n", err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status: %d %s\n", resp.StatusCode, resp.Status)
		resp.Body.Close()
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return body
}

func generateMapChunks() []byte {
	bigSquare := map[string]map[string]float64{
		"SouthWest": {"Lat": 44.980811, "Lng": -79.670110},
		"NorthEast": {"Lat": 51.998494, "Lng": -57.107759},
	}
	chunkSize := map[string]float64{
		"Lat": 0.6141262171,
		"Lng": 0.835647583,
	}

	var chunks []map[string]interface{}
	currentLat := bigSquare["SouthWest"]["Lat"]
	southWestLatLimit := bigSquare["NorthEast"]["Lat"]
	southWestLngLimit := bigSquare["NorthEast"]["Lng"]

	for currentLat < southWestLatLimit {
		currentLng := bigSquare["SouthWest"]["Lng"]
		for currentLng < southWestLngLimit {
			chunk := map[string]interface{}{
				"zoomLevel": 11,
				"mapBounds": map[string]map[string]float64{
					"SouthWest": {"Lat": currentLat, "Lng": currentLng},
					"NorthEast": {
						"Lat": math.Min(currentLat+chunkSize["Lat"], southWestLatLimit),
						"Lng": math.Min(currentLng+chunkSize["Lng"], southWestLngLimit),
					},
				},
			}
			chunks = append(chunks, chunk)
			currentLng += chunkSize["Lng"]
		}
		currentLat += chunkSize["Lat"]
	}

	output, err := json.MarshalIndent(chunks, "", "  ")
	if err != nil {
		log.Printf("Error marshaling map chunks: %v\n", err)
		return nil
	}
	return output
}

func parsePropertyLinks(infoResp MarkerInfoResponse) string {
	decodedHTML := html.UnescapeString(infoResp.D.Result.Html)
	regex := regexp.MustCompile(`href="([^"]+)"`)
	match := regex.FindStringSubmatch(decodedHTML)
	return match[1]
}

func fetchPropertyLinks(domain string, done chan bool, allData *[]OutputData, mu *sync.Mutex) {
	// Generate map chunks instead of reading from file
	chunksData := generateMapChunks()
	if chunksData == nil {
		fmt.Println("Error generating map chunks")
		done <- true
		return
	}

	var inputData []InputData
	if err := json.Unmarshal(chunksData, &inputData); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		done <- true
		return
	}

	client := &http.Client{}

	for _, item := range inputData {
		// Fetch markers
		body := fetchMarkers(domain, item, client)

		var markerResp MarkerResponse
		if err := json.Unmarshal(body, &markerResp); err != nil {
			fmt.Printf("Error parsing marker response: %v\n", err)
			continue
		}

		for _, marker := range markerResp.D.Result.Markers {
			// Fetch marker info
			body := fetchMarkerInfo(domain, marker, client)

			var infoResp MarkerInfoResponse
			if err := json.Unmarshal(body, &infoResp); err != nil {
				fmt.Printf("Error parsing info response: %v\n", err)
				continue
			}

			path := parsePropertyLinks(infoResp)

			mu.Lock()
			*allData = append(*allData, OutputData{Link: domain + path})
			fmt.Println("Successfully Fetched -", time.Now().Format("2006-01-02 15:04:05"))
			mu.Unlock()
			// Delay for 0.5 second
			time.Sleep(time.Millisecond * 500)
		}

		// Delay for 1 second before next item
		time.Sleep(time.Millisecond * 500)
	}

	done <- true
}
