package main

import (
	"fmt"
	webcrawler "webcrawler/internal/crawler"
	"webcrawler/internal/database"
)



func getAllWebsites(websiteChannel chan webcrawler.Website) map[string][] string{

	websites:= make(map[string][] string)

	for website:= range websiteChannel {

		// fmt.Println("website: ", website)
		websites[website.URL] = website.Words
	}

	return websites
}
func main() {

	sitesToCrawlChannel := make(chan string)
	linksCrawledChannel := make(chan string)
	pendingSitesChannel  := make(chan int)
	websitesChannel := make(chan webcrawler.Website)

	// var waitGroup sync.WaitGroup
	

	go func() {
		// linksCrawledChannel <- "https://theuselessweb.com/"
		linksCrawledChannel <- "https://ikennadev.netlify.app/"
	}()

	go webcrawler.ProcessLinksCrawled(sitesToCrawlChannel, linksCrawledChannel, pendingSitesChannel)
	go webcrawler.MonitorCrawling(sitesToCrawlChannel, linksCrawledChannel, pendingSitesChannel, websitesChannel)

	const nThreads = 10

	for range nThreads{
		// waitGroup.Add(1)
		go webcrawler.CrawlWebsite(sitesToCrawlChannel, linksCrawledChannel, pendingSitesChannel, websitesChannel)
	}

	// waitGroup.Wait()

	// websites:= getAllWebsites(websitesChannel)

	db:= database.DB()

	websites:= map[string][]string{
		"d1": {"new", "york", "time"},
		"d2": {"new", "york", "post"},
		"d3": {"los", "angeles", "time"},
	}

	tf:= webcrawler.TF(websites)


	idf:= webcrawler.IDF(websites)


	tfidf := webcrawler.TFIDF(tf, idf)


	for website, wordsMap := range tfidf {

		website_id:= ""
		err := db.QueryRow(database.InsertWebsite, website).Scan(&website_id)

		if err != nil {
			fmt.Printf("***DID NOT ADD WEBSITE: %v BECAUSE OF: %v ***:", website, err)
		}

		for word, tfidf := range wordsMap {
			
			keyword_id := ""

			err = db.QueryRow(database.InsertKeyword, word).Scan(&keyword_id)

			if(err != nil){
				fmt.Printf("***DID NOT ADD KEYWORD: %v BECAUSE OF: %v ***", word, err)
			}

			_, err = db.Exec(database.InsertWebsiteKeyword, website_id, keyword_id, tfidf, idf[word])

			if(err != nil){
				fmt.Printf("***DID NOT ADD TFIDF for WEBSITE: %v, KEYWORD: %v BECAUSE OF: %v ***", website, word, err)
			}
		}
	}
}