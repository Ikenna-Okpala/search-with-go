package webcrawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
	"golang.org/x/net/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var favIconMap map[string] string
var crawlCount int
var crawlMutex sync.Mutex
var mutex sync.Mutex

type WebsiteData struct {
	IconUrl string
	Title string
	Name string
	Description string
	Content string
}


func init() {
	favIconMap = make(map[string]string)
}

func findDescription(tokentype html.TokenType, token html.Token) string{

	if tokentype != html.StartTagToken || token.DataAtom.String() != "meta"{
		return ""
	}

	var isValidDescrptionTag bool

	for _, attr:= range token.Attr {

		if attr.Key == "name" && attr.Val == "description"{
			isValidDescrptionTag = true
		}
	}

	if !isValidDescrptionTag{
		return ""
	}

	index:= slices.IndexFunc(token.Attr, func(attr html.Attribute) bool{
		return attr.Key == "content"
	})

	if index < 0 {
		return ""
	}

	return token.Attr[index].Val
}

func findFavIcon(tokenType html.TokenType, token html.Token) (string){

	if tokenType != html.StartTagToken || token.DataAtom.String() != "link"{
		return ""
	}

	var isValidFavIconLink bool

	for _, attr:= range token.Attr{

		if attr.Key == "rel" {

			if attr.Val == "icon" || attr.Val == "shortcut icon" || attr.Val == "apple-touch-icon" || attr.Val == "apple-touch-icon-precomposed"{
				isValidFavIconLink = true
				break
			}
		}
	}

	if !isValidFavIconLink{
		return ""
	}

	index:= slices.IndexFunc(token.Attr, func(attr html.Attribute) bool {
		return attr.Key == "href"
	})

	if index < 0 {
		return ""
	}

	return token.Attr[index].Val
}

func isAnchorTag(tokenType html.TokenType, token html.Token) bool{

	return tokenType == html.StartTagToken && token.DataAtom.String() == "a"
}

func removeURLParams(url string) string{

	return strings.Split(url, "?") [0];
}

func formatUrl(base string, linkedUrl string) string{

	base = strings.TrimSuffix(base, "/")

	switch{
	case strings.HasPrefix(linkedUrl, "http://"):
	case strings.HasPrefix(linkedUrl, "https://"):
		if strings.Contains(linkedUrl, base){
			return removeURLParams(linkedUrl)
		}

		return ""
		

	case strings.HasPrefix(linkedUrl, "/"):

		urlParsed, err:= url.Parse(base)

		if err != nil {
			return ""
		}

		return urlParsed.Scheme + "://" + urlParsed.Hostname() + removeURLParams(linkedUrl)
	
	default:

		return base + "/" + removeURLParams(linkedUrl)

	}

	return ""
}

func getLink(website string, token html.Token) (string, bool){

	for _, attr := range token.Attr {

		if attr.Key == "href" {
			link := attr.Val

			link = formatUrl(website, link)

			if link == ""{
				break
			}


			return link, true
		}
	}

	return "", false
}

func getDomainFromUrl(websiteUrl string) string {

	parsedUrl, err:= url.Parse(websiteUrl)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	domains := strings.Split(parsedUrl.Hostname(), ".")

	return domains[len(domains) - 2] + "." + domains[len(domains) - 1]
}

func getNameFromUrl(url string) string {

	domain:= getDomainFromUrl(url)

	if domain == ""{
		return ""
	}

	domain = strings.Split(domain, ".")[0]

	return cases.Title(language.English).String(domain)
}

func parseFavIcon(website string, iconUrl string) string{

	domain:= getDomainFromUrl(website)


	if domain == "" || iconUrl == ""{
		return ""
	}

	mutex.Lock()

	if _,ok:= favIconMap[domain]; ok {

		mutex.Unlock()
		return favIconMap[domain]
	}

	

	mutex.Unlock()

	favIconUrl:= ""

	switch{
	case strings.HasPrefix(iconUrl, "http"):
		favIconUrl = iconUrl
	default:
		parsedUrl, err := url.Parse(website)
		if err != nil {
			return ""
		}
		favIconUrl = parsedUrl.Scheme + "://" + domain + iconUrl
	}
	
	resp, err:= http.Get(favIconUrl)

	if err != nil {
		return ""
	}

	client, err2 := storage.NewClient(context.Background())

	if err2 != nil {
		return ""
	}

	favIconName:= fmt.Sprintf("%s.ico", domain)

	wc:= client.Bucket(os.Getenv("BUCKET_NAME")).Object(favIconName).NewWriter(context.Background())


	if _, err:= io.Copy(wc, resp.Body); err != nil {
		return ""
	}

	if err := wc.Close(); err != nil {
		return ""
	}


	iconBucketUrl:= fmt.Sprintf("https://storage.googleapis.com/search-with-go/%s", favIconName)


	mutex.Lock()

	favIconMap[domain] = iconBucketUrl

	mutex.Unlock()

	return iconBucketUrl
	
}


func parseWebpage(website string, linksCrawledChannel chan string, crawlCapacity int)  WebsiteData{
	
	websiteData:= WebsiteData{}

	resp, err := connectToWebsite(website)

	// fmt.Println("Finished connecting to website, ", website)

	if err != nil {
		return websiteData
	}

	websiteData.Name = getNameFromUrl(website)

	defer resp.Body.Close()

	// fmt.Printf("Status of %s: %v \n", website, resp.StatusCode)


	var builder strings.Builder
	var iconUrl string
	var title string
	var description string


	tokenizer := html.NewTokenizer(resp.Body)
	previousToken := tokenizer.Token()

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken{
			break
		}

		token:= tokenizer.Token()

		if tokenType == html.StartTagToken{
			previousToken = token
		}

		if iconUrl == ""{
			iconUrl =  findFavIcon(tokenType, token)
		}

		if description == ""{
			description = findDescription(tokenType, token)
		}
		

		if isAnchorTag(tokenType, token){
			link, ok := getLink(website, token)

			if !ok{
				continue
			}

			go func(){

				crawlMutex.Lock()

				if crawlCount < crawlCapacity{
					linksCrawledChannel <- link
					crawlCount++
					fmt.Println("Crawled count: ", crawlCount)

				}
				crawlMutex.Unlock()
				
			}()
		}

		if tokenType == html.TextToken{

			if previousToken.Data != "script" && previousToken.Data != "style"{
				builder.WriteString(token.Data + " ")
			}

			if previousToken.Data == "title" && title == "" {
				title = token.Data
			}
		}

	}

	websiteData.Content = builder.String()
	websiteData.IconUrl = iconUrl
	websiteData.Title = title
	websiteData.Description = description

	return websiteData
}