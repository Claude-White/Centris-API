package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
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

// GetBroker collects all broker information from a single POST request
func GetBroker(startPosition int) (Broker, error) {
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

	return broker, nil
}

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

func main() {
	domain := "www.centris.ca"
	path := "/fr/courtiers-immobiliers"
	numberOfBrokers := getTotalNumberOfBrokers(domain, path)

	for i := range numberOfBrokers {
		broker, err := GetBroker(i)
		if err != nil {
			log.Fatalf("Error getting broker info: %v", err)
		}

		fmt.Println("Id:", broker.Id)
		fmt.Println("Name:", broker.Name)
		fmt.Println("Title:", broker.Title)
		fmt.Println("ProfilePhoto:", broker.ProfilePhoto)
		fmt.Println()
	}
}
