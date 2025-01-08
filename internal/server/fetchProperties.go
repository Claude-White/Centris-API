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
	"regexp"
	"sync"
)

func init() {
	fmt.Println("Test")
	log.Default()

	domain := "https://www.centris.ca"

	links := getAllProperties(domain)
	saveProperties(links)
}

func getAllProperties(domain string) []OutputData {
	var propertyLinks []OutputData
	ch := make(chan OutputData)
	var wg sync.WaitGroup
	chunksData := generateChunks()
	if chunksData == nil {
		log.Fatal("Error genarating chunks")
	}

	client := &http.Client{}

	for _, chunk := range chunksData {
		wg.Add(1)
		go getLinksFromChunk(domain, chunk, client, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for link := range ch {
		propertyLinks = append(propertyLinks, link)
	}

	return propertyLinks
}

func getLinksFromChunk(domain string, chunk InputData, client *http.Client, ch chan OutputData, wg *sync.WaitGroup) {
	defer wg.Done()
	markerJSON, _ := json.Marshal(chunk)
	markerReq, _ := http.NewRequest("POST", domain+"/api/property/map/GetMarkers",
		bytes.NewBuffer(markerJSON))
	markerReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(markerReq)
	if err != nil {
		fmt.Printf("Error fetching markers: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status: %d %s\n", resp.StatusCode, resp.Status)
		resp.Body.Close()
		return
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var markerResp MarkerResponse
	if err := json.Unmarshal(body, &markerResp); err != nil {
		log.Fatal("Error parsing marker response")
	}

	for _, marker := range markerResp.D.Result.Markers {
		body := getMarkerInfo(domain, marker, client)

		var infoResp MarkerInfoResponse
		if err := json.Unmarshal(body, &infoResp); err != nil {
			log.Fatal("Error parsing info response")
		}

		path := parseLinks(infoResp)

		ch <- OutputData{Link: domain + path}
	}
}

func parseLinks(infoResp MarkerInfoResponse) string {
	decodedHTML := html.UnescapeString(infoResp.D.Result.Html)
	regex := regexp.MustCompile(`href="([^"]+)"`)
	match := regex.FindStringSubmatch(decodedHTML)
	return match[1]
}

func generateChunks() []InputData {
	bigSquare := MapBounds{
		SouthWest: Coordinate{Lat: 44.980811, Lng: -79.670110},
		NorthEast: Coordinate{Lat: 51.998494, Lng: -57.107759},
	}
	chunkSize := Coordinate{
		Lat: 0.6141262171,
		Lng: 0.835647583,
	}

	var chunks []InputData
	currentLat := bigSquare.SouthWest.Lat
	southWestLatLimit := bigSquare.NorthEast.Lat
	southWestLngLimit := bigSquare.NorthEast.Lng

	for currentLat < southWestLatLimit {
		currentLng := bigSquare.SouthWest.Lng
		for currentLng < southWestLngLimit {
			chunk := InputData{
				ZoomLevel: 11,
				MapBounds: MapBounds{
					SouthWest: Coordinate{Lat: currentLat, Lng: currentLng},
					NorthEast: Coordinate{
						Lat: math.Min(currentLat+chunkSize.Lat, southWestLatLimit),
						Lng: math.Min(currentLng+chunkSize.Lng, southWestLngLimit),
					},
				},
			}
			chunks = append(chunks, chunk)
			currentLng += chunkSize.Lng
		}
		currentLat += chunkSize.Lat
	}

	return chunks
}

func getMarkerInfo(domain string, marker MarkerPosition, client *http.Client) []byte {
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

func saveProperties(arr []OutputData) {
	properties, err := json.Marshal(arr)
	if err != nil {
		log.Fatal("Failed to json parse outputData")
	}

	if err := os.WriteFile("data.json", properties, 0644); err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		return
	}
}
