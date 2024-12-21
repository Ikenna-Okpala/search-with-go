package webcrawler

import (
	"strings"

	"golang.org/x/net/html"
)

func isAnchorTag(tokenType html.TokenType, token html.Token) bool{

	return tokenType == html.StartTagToken && token.DataAtom.String() == "a"
}

func formatUrl(base string, linkedUrl string) string{

	base = strings.TrimSuffix(base, "/")

	switch{
	case strings.HasPrefix(linkedUrl, "http://"):
	case strings.HasPrefix(linkedUrl, "https://"):
		if strings.Contains(linkedUrl, base){
			return linkedUrl
		}

		return ""
		

	case strings.HasPrefix(linkedUrl, "/"):
		return base + linkedUrl

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
func parseWebpage(website string, linksCrawledChannel chan string)  string{
	
	resp, err := connectToWebsite(website)

	if err != nil {
		return ""
	}

	defer resp.Body.Close()


	var builder strings.Builder

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

		if isAnchorTag(tokenType, token){
			link, ok := getLink(website, token)

			if !ok{
				continue
			}

			go func(){
				linksCrawledChannel <- link
			}()
		}

		if tokenType == html.TextToken{

			if previousToken.Data != "script" && previousToken.Data != "style"{
				builder.WriteString(token.Data + " ")
			}
		}

	}

	return builder.String()
}