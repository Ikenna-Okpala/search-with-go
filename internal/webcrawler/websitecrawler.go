package webcrawler

import (
	"fmt"
	"strings"
	"webcrawler/internal/processor"
)

type Website struct {
	URL string
	Words [] string
	Icon string
	Title string
	Name string
	Description string
}


func CrawlWebsite(sitesToCrawlChannel chan string, linksCrawledChannel chan string, pendingSiteChannel chan int, websitesChannel chan Website, crawlCapacity int) {

	for website := range sitesToCrawlChannel{

		fmt.Println("scraping:", website)
		websiteData := parseWebpage(website, linksCrawledChannel, crawlCapacity)
		// fmt.Println("Finished processing data for url: ", website)

		// fmt.Println("icon url: ", websiteData.IconUrl)
		webText := processor.Process(websiteData.Content)
	
		fmt.Println("Crawled ", website)
		// fmt.Println("Found:", webText)
	
		if webText != ""{

			favIconUrl:= parseFavIcon(website, websiteData.IconUrl)

			// fmt.Println("icon: ", favIconUrl)

			// fmt.Println("title: ", websiteData.Title)

			// fmt.Println("name: ", websiteData.Name)

			// fmt.Println("description: ", websiteData.Description)

			websitesChannel <- Website{
				URL: website,
				Words: strings.Split(webText, " "),
				Icon: favIconUrl,
				Title: websiteData.Title,
				Name: websiteData.Name,
				Description: websiteData.Description,
			}
		}else{
			fmt.Println("Skipping: ", website)
		}
		
		pendingSiteChannel <- -1

	}

	// fmt.Println("Crawled: ", crawled)
}