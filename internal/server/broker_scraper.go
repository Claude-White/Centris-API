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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html"
)

const (
	baseUrl = "https://www.centris.ca"

	MarkerInfoUrl  = "/property/GetMarkerInfo"
	brokerEndpoint = "/Broker/GetBrokers"
	BrokerUrl      = "/fr/courtiers-immobiliers?view=Summary&uc=0"
	PropertyMapUrl = "/fr/propriete~a-vendre?view=Map&uc=0"
)

func RunBrokerScraper() {
	aspNetCoreSession, arrAffinitySameSite, numberOfBrokers := GenerateSession(baseUrl + BrokerUrl)
	brokers, brokersPhoneNumbers, brokersExternalLinks := getBrokers(baseUrl+brokerEndpoint, numberOfBrokers, aspNetCoreSession, arrAffinitySameSite)

	conn, dbErr := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if dbErr != nil {
		log.Fatalf("Failed to connect to the database: %v", dbErr)
	}
	defer conn.Close(context.Background())

	dbServer := CreateServer(conn)
	dbServer.uploadBrokersToDB(brokers, brokersPhoneNumbers, brokersExternalLinks)
}

func makeBrokerRequest(url string, startPosition int, aspNetCoreSession string, arrAffinitySameSite string) BrokerResponse {
	// Data to be sent in the JSON body
	requestData := map[string]int{
		"startPosition": startPosition,
	}

	// Convert requestData to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return BrokerResponse{}
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return BrokerResponse{}
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Add custom headers (if needed)
	req.Header.Set("Cookie", ".AspNetCore.Session="+aspNetCoreSession+"; ARRAffinitySameSite="+arrAffinitySameSite+";")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return BrokerResponse{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return BrokerResponse{}
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return BrokerResponse{}
	}

	// Parse the JSON response into the struct
	var brokerResponse BrokerResponse
	err = json.Unmarshal(body, &brokerResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return BrokerResponse{}
	}

	return brokerResponse
}

func parseBrokerName(fullName string) (string, *string, string) {
	parts := strings.Split(fullName, " ")
	var firstName, lastName string
	var middleName *string

	switch len(parts) {
	case 1:
		// Only one part, treat it as the first name
		firstName = parts[0]
		lastName = ""
		middleName = nil
	case 2:
		// Two parts, treat them as first and last name
		firstName = parts[0]
		lastName = parts[1]
		middleName = nil
	default:
		// Three or more parts, assume first name, middle name(s), and last name
		firstName = parts[0]
		lastName = parts[len(parts)-1]
		middle := strings.Join(parts[1:len(parts)-1], " ")
		middleName = &middle
	}

	return firstName, middleName, lastName
}

func getBrokers(url string, numberOfBrokers int, aspNetCoreSession string, arrAffinitySameSite string) ([]repository.CreateAllBrokersParams, [][]repository.CreateAllBrokerPhoneParams, [][]repository.CreateAllBrokerExternalLinkParams) {
	brokerResults := make(chan repository.CreateAllBrokersParams, numberOfBrokers)
	brokerPhoneResults := make(chan []repository.CreateAllBrokerPhoneParams, numberOfBrokers)
	brokerExternalLinkResults := make(chan []repository.CreateAllBrokerExternalLinkParams, numberOfBrokers)
	sem := make(chan struct{}, 50) // Limit to 50 concurrent requests

	for position := 0; position < numberOfBrokers; position++ {
		go func(pos int) {
			sem <- struct{}{}        // Acquire a spot
			defer func() { <-sem }() // Release the spot
			broker, brokerPhones, brokerExternalLinks := getBroker(url, pos, aspNetCoreSession, arrAffinitySameSite)
			brokerResults <- broker

			// Only send non-nil and non-empty arrays
			if len(brokerPhones) > 0 {
				brokerPhoneResults <- brokerPhones
			} else {
				brokerPhoneResults <- nil
			}

			if len(brokerExternalLinks) > 0 {
				brokerExternalLinkResults <- brokerExternalLinks
			} else {
				brokerExternalLinkResults <- nil
			}
		}(position)
	}

	brokers := make([]repository.CreateAllBrokersParams, 0, numberOfBrokers)
	brokersPhoneNumbers := make([][]repository.CreateAllBrokerPhoneParams, 0, numberOfBrokers)
	brokersExternalLinks := make([][]repository.CreateAllBrokerExternalLinkParams, 0, numberOfBrokers)

	var seenIDs sync.Map
	var mu sync.Mutex
	bar := progressbar.Default(int64(numberOfBrokers))
	for i := 0; i < numberOfBrokers; i++ {
		bar.Add(1)
		broker := <-brokerResults
		if _, loaded := seenIDs.LoadOrStore(broker.ID, true); !loaded {
			mu.Lock()
			brokers = append(brokers, broker)

			if phones := <-brokerPhoneResults; phones != nil {
				brokersPhoneNumbers = append(brokersPhoneNumbers, phones)
			}

			if links := <-brokerExternalLinkResults; links != nil {
				brokersExternalLinks = append(brokersExternalLinks, links)
			}
			mu.Unlock()
		}
	}

	close(brokerResults)
	close(brokerExternalLinkResults)
	close(brokerPhoneResults)
	return brokers, brokersPhoneNumbers, brokersExternalLinks
}

func getBroker(url string, startPosition int, aspNetCoreSession string, arrAffinitySameSite string) (repository.CreateAllBrokersParams, []repository.CreateAllBrokerPhoneParams, []repository.CreateAllBrokerExternalLinkParams) {
	brokerResponse := makeBrokerRequest(url, startPosition, aspNetCoreSession, arrAffinitySameSite)
	brokerName := getBrokerName(brokerResponse)
	brokerFirstName, brokerMiddleName, brokerLastName := parseBrokerName(brokerName)
	currentTime := time.Now()

	doc, err := html.Parse(strings.NewReader(brokerResponse.D.Result.Html))
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return repository.CreateAllBrokersParams{}, nil, nil // or handle error appropriately
	}

	if doc == nil {
		return repository.CreateAllBrokersParams{}, nil, nil
	}

	broker := repository.CreateAllBrokersParams{
		ID:                getBrokerID(doc),
		FirstName:         brokerFirstName,
		MiddleName:        brokerMiddleName,
		LastName:          brokerLastName,
		Title:             getBrokerTitle(brokerResponse),
		ProfilePhoto:      getBrokerProfilePhoto(doc),
		ComplementaryInfo: getBrokerComplimentaryInfo(doc),
		ServedAreas:       getBrokerServedAreas(doc),
		Presentation:      getBrokerPresentation(doc),
		CorporationName:   getBrokerCorporationName(doc),
		AgencyName:        FindElementByClass(doc, "h2", "p1 m-0 broker-info__agency-name"),
		AgencyAddress:     FindElementAttribute(doc, "a", "title", "Google Map", "title"),
		AgencyLogo:        getBrokerAgencyLogo(doc),
		CreatedAt:         &currentTime,
		UpdatedAt:         &currentTime,
	}

	// fmt.Println(broker.ID, "-", broker.FirstName, broker.LastName)
	return broker, getBrokerPhoneNumbers(broker.ID, doc), getBrokerExternalLinks(broker.ID, doc)
}

func getBrokerCorporationName(doc *html.Node) *string {
	// Find all text within the column div
	columnContent := FindElementByClassNode(doc, "div", "col-8 col-md-9")
	if columnContent == nil {
		return nil
	}

	// Find the second p1 div
	secondP1Div := FindSecondP1Div(columnContent)
	if secondP1Div != nil {
		corpName := strings.TrimSpace(ExtractText(secondP1Div))
		if corpName != "" {
			return &corpName
		}
	}

	return nil
}

func getBrokerTitle(brokerResponse BrokerResponse) string {
	parts := strings.Split(brokerResponse.D.Result.Title, ", ")
	if len(parts) < 2 {
		// Return a default or empty string if the title doesn't match expected format
		return "Courtier Immobilier"
	}
	return strings.TrimSuffix(parts[1], " - Centris.ca")
}

func getBrokerName(brokerResponse BrokerResponse) string {
	return strings.Split(brokerResponse.D.Result.Title, ", ")[0]
}

func getBrokerAgencyLogo(n *html.Node) *string {
	var brokerAgencyLogo *string
	agencyPhotoUrl := FindElementAttribute(n, "img", "class", "img-fluid broker-info-office-image", "src")
	if agencyPhotoUrl == "" {
		brokerAgencyLogo = nil
	} else {
		brokerAgencyLogo = &agencyPhotoUrl
	}

	return brokerAgencyLogo
}

func getBrokerProfilePhoto(n *html.Node) *string {
	var brokerProfilePhoto *string
	brokerPhotoUrl := FindElementAttribute(n, "img", "class", "img-fluid broker-info-broker-image", "src")
	if brokerPhotoUrl == "https://cdn.centris.ca/public/qc/consumersite/images/no-broker-pictureV3.png" {
		brokerProfilePhoto = nil
	} else {
		brokerProfilePhoto = &brokerPhotoUrl
	}

	return brokerProfilePhoto
}

func getBrokerID(n *html.Node) int64 {
	num, _ := strconv.ParseInt(FindElementAttribute(n, "div", "class", "broker-info legacy-reset  ", "data-broker-id"), 10, 64)
	return num
}

func getBrokerComplimentaryInfo(n *html.Node) *string {
	var info string
	node := FindElementByClassNode(n, "div", "col-lg-6 mb-3")

	if node == nil {
		return nil
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		for _, attr := range c.Attr {
			if attr.Key == "class" && attr.Val == "p2" {
				if c.Data != "Information complémentaire" {
					var sb strings.Builder
					for n := c.FirstChild; n != nil; n = n.NextSibling {
						if n.Type == html.TextNode {
							sb.WriteString(n.Data)
						}
					}

					text := strings.TrimSpace(sb.String())
					if text != "" {
						if info != "" {
							info += ". "
						}
						info += text
					}

					// Get the span text content
					for span := c.FirstChild; span != nil; span = span.NextSibling {
						if span.Type == html.ElementNode && span.Data == "span" {
							spanText := ExtractText(span)
							spanText = strings.TrimSpace(spanText)
							if spanText != "" {
								if info != "" {
									info += " "
								}
								info += spanText
							}
						}
					}
				}
			}
		}
	}

	if info == "" {
		return nil
	}
	return &info
}

func getBrokerServedAreas(n *html.Node) *string {
	servedAreas := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(FindElementByClass(n, "div", "col-lg-6")), "Territoire desservi"))
	if servedAreas == "" {
		return nil
	}
	servedAreasValue := &servedAreas
	return servedAreasValue
}

func getBrokerPresentation(n *html.Node) *string {
	presentation := strings.TrimSpace(FindElementByClass(n, "div", "p2 broker-summary-presentation-text"))
	if presentation == "" {
		return nil
	}
	presentationValue := &presentation
	return presentationValue
}

func getBrokerPhoneNumbers(brokerID int64, n *html.Node) []repository.CreateAllBrokerPhoneParams {
	aTags := FindElementsByAttribute(n, "itemprop", "telephone")
	var phoneNumberArr []repository.CreateAllBrokerPhoneParams
	currentTime := time.Now()

	for _, a := range aTags {
		phoneNumber := strings.ToLower(ExtractText(a))
		phoneNumberTitle := strings.ToLower(FindElementAttribute(a, "a", "class", "dropdown-item", "title"))

		var phoneNumberType string
		if strings.Contains(phoneNumberTitle, "agence") {
			phoneNumberType = "Agence"
		} else if strings.Contains(phoneNumberTitle, "moi") {
			phoneNumberType = "Courtier"
		}

		brokerPhone := repository.CreateAllBrokerPhoneParams{
			ID:        uuid.New(),
			BrokerID:  brokerID,
			Number:    phoneNumber,
			Type:      phoneNumberType,
			CreatedAt: &currentTime,
		}

		phoneNumberArr = append(phoneNumberArr, brokerPhone)
	}
	return phoneNumberArr
}

func getBrokerExternalLinks(brokerID int64, n *html.Node) []repository.CreateAllBrokerExternalLinkParams {
	externalLinks := FindElementsByAttribute(n, "class", "btn btn-outline-icon-only")
	agencyExternalLink := FindElementByClassNode(n, "a", "btn btn-outline-icon-only broker-info-office-info-icon")
	if agencyExternalLink != nil {
		externalLinks = append(externalLinks, agencyExternalLink)
	}

	currentTime := time.Now()
	var brokerExternalLinks []repository.CreateAllBrokerExternalLinkParams

	for _, externalLink := range externalLinks {
		externalLinkTitle := FindElementAttribute(externalLink, "a", "target", "_blank", "title")

		if externalLinkTitle == "" {
			continue
		}

		var brokerLinkType string

		switch externalLinkTitle {
		case "Suivez-moi sur Twitter":
			brokerLinkType = "X"
		case "Suivez-moi sur Youtube":
			brokerLinkType = "Youtube"
		case "Suivez-moi sur Facebook":
			brokerLinkType = "Facebook"
		case "Suivez-moi sur LinkedIn":
			brokerLinkType = "LinkedIn"
		case "Suivez-moi sur Instagram":
			brokerLinkType = "Instagram"
		case "Visiter mon site":
			brokerLinkType = "Site Web Personel"
		case "Visiter le site de l'agence":
			brokerLinkType = "Site Web Agence"
		default:
			continue
		}

		externalLinkHref := FindElementAttribute(externalLink, "a", "target", "_blank", "href")

		brokerExternalLink := repository.CreateAllBrokerExternalLinkParams{
			ID:        uuid.New(),
			BrokerID:  brokerID,
			Type:      brokerLinkType,
			Link:      externalLinkHref,
			CreatedAt: &currentTime,
		}

		brokerExternalLinks = append(brokerExternalLinks, brokerExternalLink)
	}

	return brokerExternalLinks
}

func flattenArray[T any](nested [][]T) []T {
	totalLength := 0
	for _, inner := range nested {
		totalLength += len(inner)
	}

	flat := make([]T, 0, totalLength)
	for _, inner := range nested {
		flat = append(flat, inner...)
	}
	return flat
}

func (s *Server) uploadBrokersToDB(brokers []repository.CreateAllBrokersParams, brokersPhoneNumbers [][]repository.CreateAllBrokerPhoneParams, brokersExternalLinks [][]repository.CreateAllBrokerExternalLinkParams) {
	ctx := context.Background()

	s.queries.DeleteAllBrokers(ctx)
	SendNotification("Process Complete", "All brokers deleted")

	brokerNum, brokerErr := s.queries.CreateAllBrokers(ctx, brokers)
	if brokerErr != nil {
		log.Printf("Failed to insert brokers")
		log.Println("Error: " + brokerErr.Error())
		SendNotification("Process Failed", fmt.Sprintf("Falied to insert %d brokers", len(brokers)))
		return
	}

	flatBrokersPhoneNumbers := flattenArray(brokersPhoneNumbers)
	_, phoneErr := s.queries.CreateAllBrokerPhone(ctx, flatBrokersPhoneNumbers)
	if phoneErr != nil {
		log.Printf("Failed to insert broker phone numbers")
		log.Println("Error: " + phoneErr.Error())
		SendNotification("Process Failed", fmt.Sprintf("Falied to insert %d broker phone numbers", len(flatBrokersPhoneNumbers)))
		return
	}

	flatBrokersExternalLinks := flattenArray(brokersExternalLinks)
	_, linkErr := s.queries.CreateAllBrokerExternalLink(ctx, flatBrokersExternalLinks)
	if linkErr != nil {
		log.Printf("Failed to insert broker phone numbers")
		log.Println("Error: " + linkErr.Error())
		SendNotification("Process Failed", fmt.Sprintf("Falied to insert %d broker external links", len(flatBrokersExternalLinks)))
		return
	}

	SendNotification("Process Complete", fmt.Sprintf("Successfully inserted %d brokers", brokerNum))
}
