package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"webcrawler/internal/database"
	"webcrawler/internal/tfidf"
	"webcrawler/internal/webcrawler"

	"github.com/joho/godotenv"
)

func getAllWebsites(websiteChannel chan webcrawler.Website) map[string] webcrawler.Website{

	websites:= make(map[string] webcrawler.Website)

	for website:= range websiteChannel {

		// fmt.Println("website: ", website)
		websites[website.URL] = website
	}

	return websites
}
func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err:= godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Env set up error")
	}

	crawlCapacity, _ := strconv.Atoi(os.Args[1:][0])

	sitesToCrawlChannel := make(chan string)
	linksCrawledChannel := make(chan string)
	pendingSitesChannel  := make(chan int)
	websitesChannel := make(chan webcrawler.Website)
	

	go func() {

		linksCrawledChannel <- "https://en.wikipedia.org/"
	}()

	go webcrawler.ProcessLinksCrawled(sitesToCrawlChannel, linksCrawledChannel, pendingSitesChannel)
	go webcrawler.MonitorCrawling(sitesToCrawlChannel, linksCrawledChannel, pendingSitesChannel, websitesChannel)

	const nThreads = 5

	for range nThreads{
		go webcrawler.CrawlWebsite(sitesToCrawlChannel, linksCrawledChannel, pendingSitesChannel, websitesChannel, crawlCapacity)
	}

	// fmt.Println("Started crawling....")

	websites:= getAllWebsites(websitesChannel)

	// for key := range maps.Keys(websites){
	// 	fmt.Println("web: ", key, " value: ", websites[key])
	// }

	db:= database.DB()

	// websites:= map[string][]string{
	// 	"d1": {"new", "york", "time"},
	// 	"d2": {"new", "york", "post"},
	// 	"d3": {"los", "angeles", "time"},
	// }

	tf:= tfidf.TF(websites)


	idf:= tfidf.IDF(websites)


	tfidf := tfidf.TFIDF(tf, idf)

	// fmt.Println("TFIDF: ", tfidf)


	for website, wordsMap := range tfidf {

		website_id:= ""

		website:= websites[website]
		err := db.QueryRow(database.InsertWebsite, website.URL, website.Title, website.Icon, website.Name, website.Description).Scan(&website_id)

		if err != nil {
			fmt.Printf("***DID NOT ADD WEBSITE: %v BECAUSE OF: %v ***\n", website, err)
		}

		for word, tfidf := range wordsMap {
			
			keyword_id := ""

			err = db.QueryRow(database.InsertKeyword, word).Scan(&keyword_id)

			if(err != nil){
				// fmt.Printf("***DID NOT ADD KEYWORD: %v BECAUSE OF: %v ***\n", word, err)
			}

			_, err = db.Exec(database.InsertWebsiteKeyword, website_id, keyword_id, tfidf, idf[word])

			if(err != nil){
				// fmt.Printf("***DID NOT ADD TFIDF for WEBSITE: %v, KEYWORD: %v BECAUSE OF: %v ***\n", website, word, err)
			}
		}
	}
}