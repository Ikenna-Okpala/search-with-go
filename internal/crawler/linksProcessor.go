package webcrawler

func ProcessLinksCrawled(sitesToCrawlChannel chan string, linksCrawledChannel chan string, pendingSitesChannel chan int) {

	visitedSites := make(map[string]bool)

	for link := range linksCrawledChannel {

		if !visitedSites[link] {

			sitesToCrawlChannel <- link
			pendingSitesChannel <- 1
			visitedSites[link] = true
		}
	}
}