package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
)

// Fetcher is an interface for crawlers
type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

type fetchedUrls struct {
	m   map[string]error
	mux sync.Mutex
}

var fetched = fetchedUrls{m: make(map[string]error)}
var loading = errors.New("currently loading url")

// Fake fetcher and result
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f *fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := (*f)[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = &fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}

	fetched.mux.Lock()
	if _, ok := fetched.m[url]; ok {
		fetched.mux.Unlock()
		fmt.Printf("<- Done with %v, already fetched.\n", url)
		return
	}
	fetched.m[url] = loading
	fetched.mux.Unlock()

	body, urls, err := fetcher.Fetch(url)
	fetched.mux.Lock()
	fetched.m[url] = err
	fetched.mux.Unlock()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	done := make(chan bool)
	for i, childUrl := range urls {
		fmt.Printf("-> Crawling child %v/%v of %v : %v.\n", i, len(urls), url, childUrl)
		go func(u string) {
			Crawl(u, depth-1, fetcher)
			done <- true
		}(childUrl)
	}

	for i, u := range urls {
		fmt.Printf("-> Done crawling child %v/%v of %v : %v.\n", i, len(urls), url, u)
		<-done
	}
	return
}

func main() {
	url := os.Args[1]
	if len(url) == 0 {
		panic("No url provided")
	}
	depth, err := strconv.ParseInt(os.Args[2], 0, 0)
	if err != nil {
		panic(err)
	}
	Crawl(url, int(depth), fetcher)
}
