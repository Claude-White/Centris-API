package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var httpClient = &http.Client{
	Timeout: 600 * time.Second, // Reuse HTTP client with a timeout
}

// Fetch broker HTML for a given start position
func getBrokerHTML(startPosition int) (string, error) {
	url := "https://www.centris.ca/Broker/GetBrokers"
	requestBody := map[string]int{
		"startPosition": startPosition,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var brokerResp BrokerResponse
	if err := json.Unmarshal(body, &brokerResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if !brokerResp.D.Succeeded {
		return "", fmt.Errorf("broker request was not successful")
	}

	return brokerResp.D.Result.Html, nil
}

func (broker *Broker) getBrokerName(doc *goquery.Document) {
	doc.Find("h1.broker-info__broker-title.h5.mb-0").Each(func(i int, s *goquery.Selection) {
		broker.Name = strings.TrimSpace(s.Text())
	})
}

func (broker *Broker) getBrokerImage(doc *goquery.Document) {
	doc.Find("img.img-fluid.broker-info-broker-image").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			broker.ProfilePhoto = strings.TrimSpace(src)
		}
	})
}

func (broker *Broker) getBrokerTitle(doc *goquery.Document) {
	doc.Find("div.p1[itemprop='jobTitle']").Each(func(i int, s *goquery.Selection) {
		broker.Title = strings.TrimSpace(s.Text())
	})
}

func (broker *Broker) getBrokerId(doc *goquery.Document) {
	doc.Find("div.broker-info.legacy-reset").Each(func(i int, s *goquery.Selection) {
		if id, exists := s.Attr("data-broker-id"); exists {
			broker.Id, _ = strconv.Atoi(id)
		}
	})
}

// Fetch broker information for a single position
func GetBroker(startPosition int, ch chan<- Broker, wg *sync.WaitGroup) (Broker, error) {
	defer wg.Done()
	// Get HTML from POST request
	html, err := getBrokerHTML(startPosition)
	if err != nil {
		return Broker{}, fmt.Errorf("error getting broker HTML: %v", err)
	}

	// Use goquery to parse the HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(html)))
	if err != nil {
		return Broker{}, fmt.Errorf("error parsing HTML: %v", err)
	}

	// Initialize Broker object
	broker := Broker{}

	// Extract broker information using goquery selectors
	broker.getBrokerId(doc)
	broker.getBrokerName(doc)
	broker.getBrokerTitle(doc)
	broker.getBrokerImage(doc)

	fmt.Println(broker)

	ch <- broker

	return broker, nil
}

func GetAllBrokers(startPositions []int, maxWorkers int) []Broker {
	brokers := []Broker{}
	ch := make(chan Broker, maxWorkers)
	var wg sync.WaitGroup

	for _, startPosition := range startPositions {
		wg.Add(1)
		go GetBroker(startPosition, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for broker := range ch {
		brokers = append(brokers, broker)
	}

	return brokers
}

// Fetch brokers concurrently
// func GetBrokersConcurrently(startPositions []int, maxWorkers int) ([]Broker, error) {
// 	var wg sync.WaitGroup
// 	var mu sync.Mutex
// 	brokers := []Broker{}
// 	errChan := make(chan error, len(startPositions))
//
// 	// Limit concurrency by using a buffered channel for workers
// 	workerChan := make(chan struct{}, maxWorkers)
//
// 	for _, startPosition := range startPositions {
// 		wg.Add(1)
// 		workerChan <- struct{}{} // Block if we've reached max concurrency
//
// 		go func(pos int) {
// 			defer wg.Done()
// 			defer func() { <-workerChan }() // Release a worker slot
//
// 			broker, err := GetBroker(pos)
// 			if err != nil {
// 				errChan <- fmt.Errorf("error fetching broker at position %d: %v", pos, err)
// 				return
// 			}
//
// 			// Use mutex to prevent race conditions while appending to the slice
// 			mu.Lock()
// 			brokers = append(brokers, broker)
// 			mu.Unlock()
// 		}(startPosition)
// 	}
//
// 	// Wait for all goroutines to finish
// 	wg.Wait()
// 	close(errChan)
//
// 	// Check if any errors occurred
// 	if len(errChan) > 0 {
// 		return nil, <-errChan
// 	}
//
// 	return brokers, nil
// }

// Fetch total number of brokers (helper function)
func getTotalNumberOfBrokers(domain string, path string) int {
	var num int
	done := make(chan bool)

	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)

	c.OnHTML("p.results-nb span.resultCount", func(e *colly.HTMLElement) {
		rawText := e.Text
		cleanedText := strings.ReplaceAll(rawText, "\u00a0", "")
		var err error
		num, err = strconv.Atoi(cleanedText)
		if err != nil {
			log.Println("Error converting string to int:", err)
		}
		done <- true
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request failed:", err)
		done <- true
	})

	go func() {
		c.Visit("https://" + domain + path)
	}()

	<-done
	return num
}

// Example usage
// func main() {
// 	// Assume you want brokers starting from positions 0, 10, 20, etc.
// 	startPositions := []int{0, 10, 20, 30, 40, 50} // Add as many as needed
//
// 	start := time.Now()
// 	brokers, err := GetBrokersConcurrently(startPositions)
// 	if err != nil {
// 		log.Fatalf("Error fetching brokers: %v", err)
// 	}
//
// 	// Print results
// 	for _, broker := range brokers {
// 		fmt.Printf("Broker ID: %d, Name: %s, Title: %s, Profile Photo: %s\n",
// 			broker.Id, broker.Name, broker.Title, broker.ProfilePhoto)
// 	}
//
// 	fmt.Printf("Fetched %d brokers in %v\n", len(brokers), time.Since(start))
// }
