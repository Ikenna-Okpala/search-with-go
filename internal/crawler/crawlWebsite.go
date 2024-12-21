package webcrawler

import (
	"fmt"
	"strings"
	"webcrawler/internal/processor"
)

type Website struct {
	URL string
	Words [] string
}


func CrawlWebsite(sitesToCrawlChannel chan string, linksCrawledChannel chan string, pendingSiteChannel chan int, websitesChannel chan Website) {

	crawled:= 0



	for website := range sitesToCrawlChannel{

		webText:= parseWebpage(website, linksCrawledChannel)
		processor.Process(webText)

		// fmt.Println("Crawling ", website)
		fmt.Println("Found:", webText)

		if webText != ""{
			websitesChannel <- Website{
				URL: website,
				Words: strings.Split(webText, " "),
			}
		}

		pendingSiteChannel <- -1

		crawled ++

	}

	// fmt.Println("Crawled: ", crawled)
}