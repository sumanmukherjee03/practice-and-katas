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

type visitationResult struct {
	url   string
	urls  []string
	err   error
	depth int
}

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
	ch := make(chan *visitationResult)

	fetch := func(u string, d int) {
		urls, err := visit(u)
		r := visitationResult{
			url:   u,
			urls:  urls,
			err:   err,
			depth: d,
		}
		ch <- &r
	}

	go fetch(url, depth)
	fetched[url] = true // make sure to not update the fetched map inside the goroutine to avoid creating a mutex

	// fetching is the variable that keeps track of how many goroutines are still running
	// the concept is similar to adding to a waitgroup so that we know how many goroutines we need to wait for to finish
	// since when this loop begins, we only have 1 goroutine, we start with 1 and start counting down
	for fetching := 1; fetching > 0; fetching-- {
		res := <-ch // This blocks until the channel has been read from

		if res.err != nil {
			log.Printf("Encountered error visiting url %s : %v", res.url, res.err)
			continue
		}

		log.Println("Fetched url : ", res.url)

		if res.depth > 0 {
			for _, u := range res.urls {
				if !fetched[u] {
					// Every time to add a new goroutine to fetch a url, we increment fetching.
					// That way, we know we want to run the loop one more time to consume from the channel
					fetching++
					go fetch(u, res.depth-1)
					fetched[u] = true
				}
			}
		}
	}

	close(ch)
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
