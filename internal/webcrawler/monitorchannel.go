package webcrawler

func MonitorCrawling(sitesToCrawlChannel chan string, linksCrawledChannel chan string, pendingSitesChannel chan int, websiteChannel chan Website) {

	total := 0

	for count := range pendingSitesChannel {

		total += count

		if total == 0 {
			close(sitesToCrawlChannel)
			close(linksCrawledChannel)
			close(pendingSitesChannel)
			close(websiteChannel)
		}
	}
}