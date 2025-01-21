package server

import (
	"bytes"
	"centris-api/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/html"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Every(5*time.Millisecond), 20) // 200 requests per second

const (
	maxConcurrentRequests = 100 // Increased from 5
	maxIdleConns          = 400 // Increased from 100
	requestTimeout        = 120 * time.Second
)

func RunPropertyScraper() {
	properties, propertiesExpenses, propertiesFeatures, propertiesPhotos, brokersProperties := getProperties()
	conn, dbErr := pgx.Connect(context.Background(), os.Getenv("SUPABASE_DB"))
	if dbErr != nil {
		log.Fatalf("Failed to connect to the database: %v", dbErr)
	}
	defer conn.Close(context.Background())

	dbServer := CreateServer(conn)
	dbServer.uploadPropertiesToDB(properties, propertiesExpenses, propertiesFeatures, propertiesPhotos, brokersProperties)
}

func getProperties() ([]repository.Property, [][]repository.PropertyExpense, [][]repository.PropertyFeature, [][]repository.PropertyPhoto, [][]repository.BrokerProperty) {
	transport := &http.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConns,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50)

	links := GetAllProperties()

	var properties []repository.Property
	var propertiesExpenses [][]repository.PropertyExpense
	var propertiesFeatures [][]repository.PropertyFeature
	var propertiesPhotos [][]repository.PropertyPhoto
	var brokersProperties [][]repository.BrokerProperty

	for _, link := range links {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Wait for rate limiter
			err := limiter.Wait(context.Background())
			if err != nil {
				log.Printf("Rate limiter error for %s: %v", url, err)
				return
			}

			// Make the request
			resp, err := client.Get(url)
			if err != nil {
				log.Printf("Error making request to %s: %v", url, err)
				return
			}
			defer resp.Body.Close()

			if strings.Contains(resp.Request.URL.String(), "listingnotfound") {
				fmt.Printf("Listing not found for URL: %s", url)
				return
			}

			// Parse the HTML content
			doc, err := html.Parse(io.NopCloser(resp.Body))
			if err != nil {
				log.Printf("Error parsing HTML for %s: %v", url, err)
				return
			}

			if link != "https://www.centris.ca" {
				fmt.Println(link)
				property := getProperty(doc)
				propertyExpenses := getPropertyExpenses(doc, property.ID)
				propertyFeatures := getPropertyFeatures(doc, property.ID)
				propertyPhotos := getPropertyPhotos(property.ID)
				brokerProperties := getPropertyBroker(doc, property.ID)

				properties = append(properties, property)
				propertiesExpenses = append(propertiesExpenses, propertyExpenses)
				propertiesFeatures = append(propertiesFeatures, propertyFeatures)
				propertiesPhotos = append(propertiesPhotos, propertyPhotos)
				brokersProperties = append(brokersProperties, brokerProperties)
			}

		}(link)
	}
	wg.Wait()

	return properties, propertiesExpenses, propertiesFeatures, propertiesPhotos, brokersProperties
}

func GetAllProperties() []string {
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

	aspNetCoreSession, arrAffinitySameSite, _ := GenerateSession(baseUrl + PropertyMapUrl)

	pins := getAllPins(client, aspNetCoreSession, arrAffinitySameSite)
	housesHTML := getAllHouses(client, pins, aspNetCoreSession, arrAffinitySameSite)
	return housesHTML
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

				link := baseUrl + GetAllHouseInPin(client, aspNetCoreSession, arrAffinitySameSite, pin, idx)
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

		return FindElementAttribute(houseHTML, "a", "class", "property-thumbnail-summary-link", "href")
	}

	return ""
}

func getPropertyCategory(title string) string {
	return strings.Split(title, " ")[0]
}

func getPropertyId(doc *html.Node) int64 {
	element := FindElementById(doc, "ListingDisplayId")
	elementText := ExtractText(element)
	propertyId, _ := strconv.ParseInt(elementText, 10, 64)
	return propertyId
}

func getPropertyPrice(doc *html.Node) float32 {
	propertyPriceString := FindElementAttribute(doc, "meta", "itemprop", "price", "content")
	propertyPriceFloat64, _ := strconv.ParseFloat(propertyPriceString, 64)
	return float32(propertyPriceFloat64)
}

func getPropertyCoordinate(doc *html.Node, itempropValue string) float32 {
	propertyLatitudeFloat64, _ := strconv.ParseFloat(FindElementAttribute(doc, "meta", "itemprop", itempropValue, "content"), 64)
	return float32(propertyLatitudeFloat64)
}

func getPropretyDescription(doc *html.Node) *string {
	var propertyDescription *string
	element := FindElementByAttribute(doc, "itemprop", "description")
	if element == nil {
		return nil
	}
	description := ExtractText(element)
	if description != "" {
		propertyDescription = &description // Take the address of the variable
	} else {
		propertyDescription = nil
	}

	return propertyDescription
}

func getPropertyRoomNumber(doc *html.Node, className string) *int32 {
	var propertyRoomNumber *int32
	elementNode := FindElementByClassNode(doc, "div", "col-lg-3 col-sm-6 "+className)
	if elementNode == nil {
		return nil
	}

	roomNumber := ExtractText(elementNode)
	if roomNumber != "" {
		regex := regexp.MustCompile(`\b([1-9][0-9]?|0)\b`)
		matches := regex.FindAllString(roomNumber, -1)
		var sum int32 = 0 // Initialize a sum variable

		for _, match := range matches {
			if className == "cac" {
				parsedRoomNumber, err := strconv.ParseInt(matches[0], 10, 32)
				if err == nil {
					num := int32(parsedRoomNumber)
					return &num
				}
			}

			parsedBathroomNumber, err := strconv.ParseInt(match, 10, 32)
			if err == nil {
				sum += int32(parsedBathroomNumber) // Add to the sum
			}
		}
		propertyRoomNumber = &sum // Assign the pointer to the final sum
	}

	return propertyRoomNumber
}

func getPropertyTitle(doc *html.Node) string {
	element := FindElementByAttribute(doc, "data-id", "PageTitle")
	elementText := ExtractText(element)
	return elementText
}

func getPropertyAddress(doc *html.Node) string {
	element := FindElementByAttribute(doc, "itemprop", "address")
	elementText := ExtractText(element)
	return elementText
}

func getPropertyCityName(doc *html.Node) string {
	parentElement := FindElementById(doc, "divStatistique")
	element := FindElementByTagName(parentElement, "h2")
	elementChild := FindElementByTagName(element, "span")
	elementChildText := ExtractText(elementChild)
	return elementChildText
}

func getProperty(doc *html.Node) repository.Property {

	currentTime := time.Now()

	// propertyTransactionType := getPropertyTransactionType(pageTitle)

	return repository.Property{
		ID:             getPropertyId(doc),
		Title:          getPropertyTitle(doc),
		Category:       getPropertyCategory(getPropertyTitle(doc)),
		Address:        getPropertyAddress(doc),
		CityName:       getPropertyCityName(doc),
		Price:          getPropertyPrice(doc),
		Description:    getPropretyDescription(doc),
		BedroomNumber:  getPropertyRoomNumber(doc, "cac"),
		RoomNumber:     getPropertyRoomNumber(doc, "piece"),
		BathroomNumber: getPropertyRoomNumber(doc, "sdb"),
		Latitude:       getPropertyCoordinate(doc, "latitude"),
		Longitude:      getPropertyCoordinate(doc, "longitude"),
		CreatedAt:      &currentTime,
		UpdatedAt:      &currentTime,
	}
}

func getPropertyExpenses(doc *html.Node, propertyId int64) []repository.PropertyExpense {
	currentTime := time.Now()
	containerElement := FindElementByClassNode(doc, "div", "row financial-details-tables")
	monthlyTables := FindElementsByAttribute(containerElement, "class", "financial-details-table financial-details-table-monthly")
	var propertyExpenses []repository.PropertyExpense

	for _, element := range monthlyTables {
		tableTitle := FindElementByClass(element, "th", "col pl-0 financial-details-table-title")
		tableBody := FindElementByTagName(element, "tbody")
		tableRows := FindElementsByTagName(tableBody, "tr")
		var expenseType string
		for _, row := range tableRows {
			rowTitle := ExtractText(FindElementByTagName(row, "td"))
			rowValue := strings.TrimSpace(strings.ReplaceAll(FindElementByClass(row, "td", "text-right"), "$", ""))
			floatValue, err := strconv.ParseFloat(rowValue, 32)
			if err != nil {
				continue
			}

			if tableTitle == "Taxes" {
				expenseType = tableTitle + " " + rowTitle
			} else if tableTitle == "Dépenses" {
				expenseType = rowTitle
			}

			float32Value := float32(floatValue)
			annualValue := (float32Value * 12)

			propertyExpense := repository.PropertyExpense{
				ID:           uuid.New(),
				PropertyID:   propertyId,
				Type:         expenseType,
				MonthlyPrice: float32Value,
				AnnualPrice:  annualValue,
				CreatedAt:    &currentTime,
			}

			propertyExpenses = append(propertyExpenses, propertyExpense)
		}

	}

	return propertyExpenses
}

func getPropertyFeatures(doc *html.Node, propertyId int64) []repository.PropertyFeature {
	elements := FindElementsByAttribute(doc, "class", "col-lg-3 col-sm-6 carac-container")
	var propertyFeatures []repository.PropertyFeature

	for _, element := range elements {
		elementTitle := FindElementByClass(element, "div", "carac-title")
		elementValue := FindElementByClass(element, "div", "carac-value")

		if elementTitle != "" && elementValue != "" {
			propertyFeature := repository.PropertyFeature{
				ID:         uuid.New(),
				PropertyID: propertyId,
				Title:      elementTitle,
				Value:      elementValue,
			}
			propertyFeatures = append(propertyFeatures, propertyFeature)
		}
	}

	return propertyFeatures
}

func getPropertyPhotos(propertyId int64) []repository.PropertyPhoto {
	currentTime := time.Now()
	url := "https://www.centris.ca/Property/PhotoViewerDataListing"
	var propertyPhotos []repository.PropertyPhoto

	// Prepare request body
	body := RequestBodyPhoto{
		Lang:                   "fr",
		CentrisNo:              propertyId,
		Track:                  true,
		AuthorizationMediaCode: "999",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	// Parse response
	var photoResponse PhotoResponse
	err = json.NewDecoder(resp.Body).Decode(&photoResponse)
	if err != nil {
		return nil
	}

	for _, photo := range photoResponse.PhotoList {

		photoUrl := photo.UrlThumb
		photoUrlQueryParameterIndex := strings.Index(photoUrl, "&t")
		photoUrl = photoUrl[:photoUrlQueryParameterIndex] + "&t=pi&f=I"

		propertyPhoto := repository.PropertyPhoto{
			ID:          uuid.New(),
			PropertyID:  propertyId,
			Link:        photoUrl,
			Description: &photo.Desc,
			CreatedAt:   &currentTime,
		}

		propertyPhotos = append(propertyPhotos, propertyPhoto)
	}

	return propertyPhotos
}

func getPropertyBroker(doc *html.Node, propertyId int64) []repository.BrokerProperty {
	currentTime := time.Now()
	elements := FindElementsByAttribute(doc, "class", "broker-info legacy-reset  ")
	brokerProperties := []repository.BrokerProperty{}

	for _, element := range elements {
		brokerIdString := FindElementAttribute(element, "div", "class", "broker-info legacy-reset  ", "data-broker-id")
		brokerId, err := strconv.ParseInt(brokerIdString, 10, 64)
		if err != nil {
			continue
		}

		brokerProperty := repository.BrokerProperty{
			ID:         uuid.New(),
			BrokerID:   brokerId,
			PropertyID: propertyId,
			CreatedAt:  &currentTime,
		}

		brokerProperties = append(brokerProperties, brokerProperty)
	}

	return brokerProperties
}

func (s *Server) uploadPropertiesToDB(properties []repository.Property, propertiesExpenses [][]repository.PropertyExpense, propertiesFeatures [][]repository.PropertyFeature, propertiesPhotos [][]repository.PropertyPhoto, brokersProperties [][]repository.BrokerProperty) {
	ctx := context.Background()

	s.queries.DeleteAllProperties(ctx)
	SendNotification("Process Complete", "All properties deleted.")

	allProperties := []repository.CreateAllPropertiesParams{}
	for _, property := range properties {
		propertyParams := repository.CreateAllPropertiesParams{
			ID:             property.ID,
			Title:          property.Title,
			Category:       property.Category,
			Address:        property.Address,
			CityName:       property.CityName,
			Price:          property.Price,
			Description:    property.Description,
			BedroomNumber:  property.BedroomNumber,
			RoomNumber:     property.RoomNumber,
			BathroomNumber: property.BathroomNumber,
			Latitude:       property.Latitude,
			Longitude:      property.Longitude,
		}
		allProperties = append(allProperties, propertyParams)
		fmt.Println("Property" + string(propertyParams.ID) + "Added")
	}
	count, err := s.queries.CreateAllProperties(ctx, allProperties)
	if err != nil {
		log.Printf("Failed to insert %d properties: %s", count, err)
		SendNotification("Failed to Insert", "An error occured while inserting all properties:\n"+err.Error())
	} else {
		fmt.Printf("Successfully inserted %d properties\n", count)
		SendNotification("Process Complete", "All properties successfully inserted.")
	}

	flatPropertiesExpenses := flattenArray(propertiesExpenses)
	allPropertiesExpenses := []repository.CreateAllPropertiesExpensesParams{}
	for _, propertyExpense := range flatPropertiesExpenses {
		allPropertiesExpenses = append(allPropertiesExpenses, repository.CreateAllPropertiesExpensesParams(propertyExpense))
		fmt.Println("Property expense" + propertyExpense.ID.String() + "Added")
	}
	count, err = s.queries.CreateAllPropertiesExpenses(ctx, allPropertiesExpenses)
	if err != nil {
		log.Printf("Failed to insert %d property expenses: %s", count, err)
		SendNotification("Failed to Insert", "An error occured while inserting all property expenses:\n"+err.Error())
	} else {
		fmt.Printf("Successfully inserted %d property expenses\n", count)
		SendNotification("Process Complete", "All property expenses successfully inserted.")
	}

	flatPropertiesFeatures := flattenArray(propertiesFeatures)
	allPropertiesFeatures := make([]repository.CreateAllPropertiesFeaturesParams, 0, len(flatPropertiesFeatures))
	for _, propertyFeature := range flatPropertiesFeatures {
		allPropertiesFeatures = append(allPropertiesFeatures, repository.CreateAllPropertiesFeaturesParams(propertyFeature))
		fmt.Println("Property feature" + propertyFeature.ID.String() + "Added")
	}
	count, err = s.queries.CreateAllPropertiesFeatures(ctx, allPropertiesFeatures)
	if err != nil {
		log.Printf("Failed to insert %d property features: %s", count, err)
		SendNotification("Failed to Insert", "An error occured while inserting all property features:\n"+err.Error())
	} else {
		fmt.Printf("Successfully inserted %d properties features\n", count)
		SendNotification("Process Complete", "All property features successfully inserted.")
	}

	flatPropertiesPhotos := flattenArray(propertiesPhotos)
	allPropertiesPhotos := []repository.CreateAllPropertiesPhotosParams{}
	for _, propertyPhoto := range flatPropertiesPhotos {
		allPropertiesPhotos = append(allPropertiesPhotos, repository.CreateAllPropertiesPhotosParams(propertyPhoto))
		fmt.Println("Property photo" + propertyPhoto.ID.String() + "Added")
	}
	count, err = s.queries.CreateAllPropertiesPhotos(ctx, allPropertiesPhotos)
	if err != nil {
		log.Printf("Failed to insert %d property photos: %s", count, err)
		SendNotification("Failed to Insert", "An error occured while inserting all property photos:\n"+err.Error())
	} else {
		fmt.Printf("Successfully inserted %d properties\n", count)
		SendNotification("Process Complete", "All property photos successfully inserted.")
	}

	flatBrokerProperties := flattenArray(brokersProperties)
	allBrokersProperties := []repository.CreateAllBrokersPropertiesParams{}
	for _, brokerProperty := range flatBrokerProperties {
		allBrokersProperties = append(allBrokersProperties, repository.CreateAllBrokersPropertiesParams(brokerProperty))
		fmt.Println("Broker property" + brokerProperty.ID.String() + "Added")
	}
	count, err = s.queries.CreateAllBrokersProperties(ctx, allBrokersProperties)
	if err != nil {
		log.Printf("Failed to insert %d broker properties: %s", count, err)
		SendNotification("Failed to Insert", "An error occured while inserting all broker properties:\n"+err.Error())
	} else {
		fmt.Printf("Successfully inserted %d broker properties\n", count)
		SendNotification("Process Complete", "All broker properties successfully inserted.")
	}
}
