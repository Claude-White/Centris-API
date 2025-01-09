package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func GetAllProperties() {
	client := &http.Client{
	}

	resp, err := http.Get("https://www.centris.ca")
	if err != nil {
		log.Fatalf("Failed to make GET request: %v", err)
	}

	cookies := resp.Header["Set-Cookie"]

	var aspNetCoreSession string
	var arrAffinitySameSite string

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
	ch := make(chan string, 1500)
	var wg sync.WaitGroup

	
	for _, pin := range pins {
		for i := 0; i < pin.PointsCount; i++ {
			wg.Add(1)
			go GetAllHouseInPin(client, aspNetCoreSession, arrAffinitySameSite, pin, i, ch, &wg)
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for houseLink := range ch {
		houseLinks = append(houseLinks, houseLink)
	}
	return houseLinks
}

func GetAllHouseInPin(client *http.Client, aspNetCoreSession string, arrAffinitySameSite string, pin Marker, idx int, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
		bodyData := MarkerInfoInputData{
			PageIndex: idx,
			ZoomLevel: 11,
			Latitude: pin.Position.Lat,
			Longitude: pin.Position.Lng,
			GeoHash: pin.GeoHash,
		}
		pinJSON, _ := json.Marshal(bodyData)
		pinReq, _ := http.NewRequest("POST", "https://www.centris.ca/property/GetMarkerInfo",
			bytes.NewBuffer(pinJSON))
			pinReq.Header.Set("Content-Type", "application/json")
			pinReq.Header.Set("Cookie", ".AspNetCore.Session="+aspNetCoreSession+"; ARRAffinitySameSite="+arrAffinitySameSite+";")

		resp, err := client.Do(pinReq)
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

		var response MarkerInfoResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Error unmarshaling response: %v", err)
		}
		houseHTML, _ := html.Parse(strings.NewReader(response.D.Result.Html))
		houseLink := findElementAttribute(houseHTML, "a", "class", "property-thumbnail-summary-link", "href")
		fmt.Println(houseLink)
		ch <- houseLink
}