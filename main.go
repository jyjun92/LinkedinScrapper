package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//https://www.linkedin.com/jobs/search?keywords=
//python&location=Downers%20Grove%2C%20Illinois
//%2C%20United%20States&trk=homepage-jobseeker_jobs-search-bar_search-submit&redirect=false&position=1&pageNum=0
var location_city = "Downers%20Grove"
var location_state = "Illinois"
var keyword = "python"

func main() {
	strings.ReplaceAll(strings.TrimSpace(location_city), " ", "%20")
	strings.ReplaceAll(strings.TrimSpace(keyword), " ", "%20")

	baseURL := "https://www.linkedin.com/jobs/search?keywords=" + keyword +
		"&location=" + location_city + "%2C%20" + location_state +
		"%2C%20United%20States&trk=homepage-jobseeker_jobs-search-bar_search-submit&redirect=false&position=1&pageNum=0"

	cards := getCards(baseURL)
	cards.Each(func(i int, card *goquery.Selection) {
		extractJob(card)
	})
}

func getCards(baseURL string) *goquery.Selection {
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	doc, err2 := goquery.NewDocumentFromReader(res.Body)
	checkErr(err2)
	searchCards := doc.Find(".result-card")
	fmt.Println("Number of Cards Found:", searchCards.Length())
	return searchCards
}

func extractJob(card *goquery.Selection) {
	//extract details, hreff
	id, _ := card.Attr("data-id")
	title := CleanString(card.Find(".job-result-card__title").Text())
	location := CleanString(card.Find(".job-result-card__location").Text())
	summary := CleanString(card.Find(".job-result-card__snippet").Text())
	href, _ := card.Find(".result-card__full-card-link").Attr("href")

	fmt.Println("Id: ", id, " Title: ", title, " Location: ", location, " Summary: ", summary)
	fmt.Println("Link: ", href)

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with StatusCode:", res.StatusCode)
	}
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
