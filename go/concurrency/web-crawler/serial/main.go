package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/net/html"
)

var (
	fetched  = make(map[string]bool)
	urlRegex = regexp.MustCompile(`(https|http)://.*`)
)

func main() {
	describe()
	start := time.Now()
	crawl("https://github.com", 2)
	elapsed := time.Since(start)
	fmt.Println("Total time taken : ", elapsed)
}

// crawl a url upto a given depth
// To convert this to a more parallel program, we need to convert the recursion to an iteration and use channels for communication
func crawl(url string, depth int) {
	if depth < 0 {
		return
	}

	urls, err := visit(url)
	if err != nil {
		log.Printf("Encountered error visiting url %s : %v", url, err)
		return
	}
	fetched[url] = true // mark this url as fetched so that this link is not revisited again
	log.Println("Fetched url : ", url)

	// range over all the new links that were generated from visiting the given url and recursively call crawl on them
	for _, u := range urls {
		if !fetched[u] {
			crawl(u, depth-1) // Remember to reduce the depth by 1
		}
	}

	return
}

// visit a url, get the http response, parse the html and call traverse function to extract all the links from that page
func visit(url string) ([]string, error) {
	var links []string

	if !urlRegex.MatchString(url) {
		return links, fmt.Errorf("Not a valid url to visit : %s", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return links, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return links, fmt.Errorf("Received status code %s when fetching url %s", resp.Status, url)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return links, err
	}

	links = traverse(doc, nil)
	return links, nil
}

// traverse html nodes recursively to keep finding the link tags
func traverse(node *html.Node, links []string) []string {
	// If node is a link tag then extract the link and update the links slice and then continue further
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" && urlRegex.MatchString(a.Val) {
				links = append(links, a.Val)
			}
		}
	}

	// If the current html node is not a link tag, then keep traversing down to all it's children
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		links = traverse(c, links)
	}

	return links
}

func describe() {
	str := `
Recursively crawl web pages till a depth of 2.

_____________________
	`
	fmt.Println(str)
}
