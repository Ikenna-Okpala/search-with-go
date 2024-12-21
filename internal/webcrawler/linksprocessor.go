package webcrawler

import (
	"strings"
)

func ProcessLinksCrawled(sitesToCrawlChannel chan string, linksCrawledChannel chan string, pendingSitesChannel chan int) {

	visitedSites := make(map[string]bool)

	for link := range linksCrawledChannel {

		link = strings.TrimSuffix(link, "/")

		if !visitedSites[link]{

			sitesToCrawlChannel <- link
			pendingSitesChannel <- 1
			visitedSites[link] = true
		}
	}
}