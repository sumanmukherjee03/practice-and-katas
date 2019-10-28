package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type repoData struct {
	id           int
	name         string
	url          string
	languagesURL string
}

var githubNetTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 3 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 3 * time.Second,
}

var githubHTTPClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: githubNetTransport,
}

var wg sync.WaitGroup

func requestGithubAPI(since int, ch chan<- repoData) {
	url := fmt.Sprintf("https://api.github.com/repositories?since=%d", since)
	resp, err := githubHTTPClient.Get(url)
	if err != nil {
		log.Fatalln("ERROR :>> Encountered error making request to the github api", err)
	}
	defer resp.Body.Close() // Dont forget to do this to prevent memory leaks

	var result []map[string]interface{} // array of maps
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("ERROR :>> Encountered error reading response from api request", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalln("ERROR :>> Encountered error in unmarshalling json response from api", err)
	}

	var nextSince int
	for _, m := range result {
		tempId, ok1 := m["id"].(float64)
		if !ok1 {
			log.Fatalln("Could not perform type assertion for the data received from the api for id")
		}
		id := int(tempId)

		name, ok2 := m["full_name"].(string)
		if !ok2 {
			log.Fatalln("Could not perform type assertion for the data received from the api for full_name")
		}

		url, ok3 := m["url"].(string)
		if !ok3 {
			log.Fatalln("Could not perform type assertion for the data received from the api for url")
		}

		languagesURL, ok4 := m["languages_url"].(string)
		if !ok4 {
			log.Fatalln("Could not perform type assertion for the data received from the languages_url")
		}

		nextSince = id
		ch <- repoData{id, name, url, languagesURL}
	}
	defer close(ch)

	fmt.Println("The next slice of repo data to fetch is since ", nextSince)
	wg.Done()
}

func populateLanguagesForRepo(ch <-chan repoData) {
	for r := range ch {
		fmt.Println(r)
	}
	wg.Done()
}

func main() {
	since, err := strconv.ParseInt(os.Args[1], 0, 0)
	if err != nil {
		panic(err)
	}

	ch := make(chan repoData)
	wg.Add(2)
	go requestGithubAPI(int(since), ch)
	go populateLanguagesForRepo(ch)
	wg.Wait()
}
