package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Listing struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Company  string `json:"company"`
	Region   string `json:"region"`
}

func parseJobPage(url string) Listing {
	url = "https://weworkremotely.com/" + url
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find all links and process them with the function
	// defined earlier
	name, _ := document.Find("body > div > div.content > div.listing-header > div.listing-header-container > h1").First().Html()
	location, _ := document.Find("body > div > div.content > div.listing-header > div.listing-header-container > h2 > span.location").First().Html()
	company, _ := document.Find("body > div > div.content > div.listing-header > div.listing-header-container > h2 > span.company").First().Html()
	region, _ := document.Find("body > div > div.content > div.listing-header > div.listing-header-container > h2 > span.region").First().Html()

	listing := Listing{name, location, company, region}
	return listing
}

func getListings(url string) []Listing {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	remoteJobUrls := make([]string, 0)

	processElement := func(index int, element *goquery.Selection) {
		href, exists := element.Attr("href")
		if exists {
			if strings.Contains(href, "remote-jobs") {
				//	listing = parseJobPage(href)
				remoteJobUrls = append(remoteJobUrls, href)
			}
		}

	}

	document.Find(".jobs li > a").Each(processElement)

	var wg sync.WaitGroup
	goroutines := make(chan struct{}, 100)
	listings := make([]Listing, 0)
	for _, url := range remoteJobUrls {
		wg.Add(1) // increasing wait group size to the no of urls
		goroutines <- struct{}{}
		go func(url string) {
			listing := parseJobPage(url)
			// fmt.Println(listing)
			<-goroutines
			listings = append(listings, listing)
			wg.Done()
		}(url)
	}
	wg.Wait()
	return listings
}

func main() {
	fmt.Println("About to start parsing jobs")
	const BASE_URL = "https://weworkremotely.com/categories/remote-programming-jobs"
	mainPage := getListings(BASE_URL)
	result, _ := json.Marshal(mainPage)
	ioutil.WriteFile("listings.json", result, 0644)
}
