package server

import (
	"strings"

	"golang.org/x/net/html"
)

func ExtractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(strings.TrimSpace(ExtractText(c)))
	}
	return sb.String()
}

func FindElementByClass(n *html.Node, tagName, className string) string {
	if n.Type == html.ElementNode && n.Data == tagName {
		// Check if the element has the desired class attribute
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == className {
				return ExtractText(n)
			}
		}
	}

	// Traverse child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := FindElementByClass(c, tagName, className); result != "" {
			return result
		}
	}

	return ""
}

func FindElementAttribute(n *html.Node, tagName, attrName, attrValue, targetAttr string) string {
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
					sb.WriteString(strings.TrimSpace(ExtractText(c)))
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
		if result := FindElementAttribute(c, tagName, attrName, attrValue, targetAttr); result != "" {
			return result
		}
	}

	return ""
}

func FindSecondP1Div(n *html.Node) *html.Node {
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

func FindElementByClassNode(n *html.Node, tagName, className string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tagName {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, className) {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := FindElementByClassNode(c, tagName, className); result != nil {
			return result
		}
	}

	return nil
}

func FindElementByTagName(n *html.Node, tagName string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tagName {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := FindElementByTagName(c, tagName); result != nil {
			return result
		}
	}

	return nil
}

func FindElementsByAttribute(n *html.Node, attrName, attrValue string) []*html.Node {
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
		childMatches := FindElementsByAttribute(c, attrName, attrValue)
		matches = append(matches, childMatches...)
	}

	return matches
}
