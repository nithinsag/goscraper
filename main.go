package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Listing struct {
	name, location, company, region string
}

func parseJobPage(url string) Listing {
	url = "https://weworkremotely.com/" + url
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	// dataInBytes, err := ioutil.ReadAll(response.Body)
	// pageContent := string(dataInBytes)

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

func processElement(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element

	href, exists := element.Attr("href")
	if exists {
		if strings.Contains(href, "remote-jobs") {
			listing := parseJobPage(href)
			fmt.Println(listing)
		}
	}

}

func getContent(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	// dataInBytes, err := ioutil.ReadAll(response.Body)
	// pageContent := string(dataInBytes)

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find all links and process them with the function
	// defined earlier
	document.Find(".jobs li > a").Each(processElement)
	return "finished"

}

func main() {
	fmt.Println("About to start parsing jobs")
	const BASE_URL = "https://weworkremotely.com/categories/remote-programming-jobs"
	mainPage := getContent(BASE_URL)
	fmt.Println(mainPage)

}
