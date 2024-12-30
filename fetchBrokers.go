package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

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

func getBrokerName(c *colly.Collector) string {
	var fullName string

	c.OnHTML("h1.broker-info__broker-title.h5.mb-0", func(e *colly.HTMLElement) {
		fullName = strings.TrimSpace(e.Text)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Failed to get broker name:", err)
		return
	})

	return fullName
}

func getBroker(domain string, path string) {
	done := make(chan bool)

	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)

	getBrokerName(c)

	go func() {
		c.Visit("https://" + domain + path)
	}()
	<-done
}

// func main() {
// 	domain := "www.centris.ca"
// 	path := "/fr/courtiers-immobiliers"
//
// 	brokers := getTotalNumberOfBrokers(domain, path)
//
// }
