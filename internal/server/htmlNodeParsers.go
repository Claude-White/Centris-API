package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(strings.TrimSpace(extractText(c)))
	}
	return sb.String()
}

func findElementByClass(n *html.Node, tagName, className string) string {
	if n.Type == html.ElementNode && n.Data == tagName {
		// Check if the element has the desired class attribute
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == className {
				return extractText(n)
			}
		}
	}

	// Traverse child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findElementByClass(c, tagName, className); result != "" {
			return result
		}
	}

	return ""
}

func findElementAttribute(n *html.Node, tagName, attrName, attrValue, targetAttr string) string {
	if n.Type == html.ElementNode && n.Data == tagName {
		// First find if this element has the attribute we're searching for
		hasMatchingAttr := false
		for _, attr := range n.Attr {
			if attr.Key == attrName && attr.Val == attrValue {
				hasMatchingAttr = true
				break
			}
		}

		// If we found the matching attribute, return the target attribute's value
		if hasMatchingAttr {
			if attrValue == "Google Map" {
				var sb strings.Builder
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					sb.WriteString(strings.TrimSpace(extractText(c)))
				}
				return sb.String()
			}

			for _, attr := range n.Attr {
				if attr.Key == targetAttr {
					return attr.Val
				}
			}
		}
	}

	// Recursively search child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findElementAttribute(c, tagName, attrName, attrValue, targetAttr); result != "" {
			return result
		}
	}

	return ""
}

func findSecondP1Div(n *html.Node) *html.Node {
	var p1Count int
	var findP1 func(*html.Node) *html.Node

	findP1 = func(node *html.Node) *html.Node {
		if node.Type == html.ElementNode && node.Data == "div" {
			for _, attr := range node.Attr {
				if attr.Key == "class" && attr.Val == "p1" {
					p1Count++
					if p1Count == 2 {
						return node
					}
					break
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if result := findP1(c); result != nil {
				return result
			}
		}

		return nil
	}

	return findP1(n)
}

func findElementByClassNode(n *html.Node, tagName, className string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tagName {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, className) {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findElementByClassNode(c, tagName, className); result != nil {
			return result
		}
	}

	return nil
}

func findElementByTagName(n *html.Node, tagName string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tagName {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findElementByTagName(c, tagName); result != nil {
			return result
		}
	}

	return nil
}

func findElementsByAttribute(n *html.Node, attrName, attrValue string) []*html.Node {
	var matches []*html.Node

	if n == nil {
		return matches
	}

	// Check if current node has the attribute with specified value
	for _, attr := range n.Attr {
		if attr.Key == attrName && attr.Val == attrValue {
			matches = append(matches, n)
			break // Break after finding the first matching attribute
		}
	}

	// Recursively search through child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		childMatches := findElementsByAttribute(c, attrName, attrValue)
		matches = append(matches, childMatches...)
	}

	return matches
}

func generateSession(url string) (string, string, int) {
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
			aspNetCoreSession = extractCookieValue(cookie, ".AspNetCore.Session")
		} else if strings.Contains(cookie, "ARRAffinitySameSite=") {
			arrAffinitySameSite = extractCookieValue(cookie, "ARRAffinitySameSite")
		}
	}

	numberOfBrokers := 0
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return aspNetCoreSession, arrAffinitySameSite, numberOfBrokers
	}

	numberOfBrokers, _ = strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(findElementByClass(doc, "span", "resultCount"), "\u00a0", ""), " ", ""))

	return aspNetCoreSession, arrAffinitySameSite, numberOfBrokers
}

func extractCookieValue(cookie string, key string) string {
	parts := strings.Split(cookie, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, key+"=") {
			return strings.TrimPrefix(part, key+"=")
		}
	}
	return ""
}
