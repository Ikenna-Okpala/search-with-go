package webcrawler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/benjaminestes/robots"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var robotsMap map[string] *robots.Robots


func init(){
	robotsMap = make(map[string] *robots.Robots)
}

func isAllowedToScrape(robot * robots.Robots, url string) bool{
	return robot.Test("*", url)
}

func connectToWebsite(website string) (*http.Response, error){

	
	// fmt.Println("Beginning to scrape website, ", website)
	robotUrl, err1 := robots.Locate(website)

	shouldCrawl:= true

	if err1 != nil {
		return nil, errors.New("bad URL format")
	}

	if robot, ok:=robotsMap[robotUrl]; ok {

		if robot != nil && !isAllowedToScrape(robot, website){
			shouldCrawl = false
		}
	}else{
		respRobot, err2 := http.Get(robotUrl)

		if err2 == nil {
	
			defer respRobot.Body.Close()
	
			robot, err3 := robots.From(respRobot.StatusCode, respRobot.Body)
	
			if err3 != nil {
				robotsMap[robotUrl] = nil
				fmt.Printf("could not parse robots.txt of %s because of %v", website, err3)
			}else{
	
				if !isAllowedToScrape(robot, website){
					shouldCrawl = false
				}
			}
		}else{
			robotsMap[robotUrl] = nil
		}
	}
	

	if !shouldCrawl {
		return nil, errors.New("robot.txt forbids crawling")
	}

	statusChanel := make(chan int)

	// fmt.Println("Finished scraping robot.txt, ", website)

	ctx, cancel:= chromedp.NewContext(context.Background())

	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Second * 30)

	defer cancel()

	
	// fmt.Println("Scraping..... ", website)
	chromedp.ListenTarget(ctx, func(ev interface {}){

		switch event := ev.(type){
		case * network.EventResponseReceived: 
			response:= event.Response
			//assume the index.html is the first thing returned

			go func(){
				statusChanel <- int(response.Status)
			}()
				
		}
	})	


	var htmlContent string

	err:= chromedp.Run(
		ctx,
		network.Enable(),
		chromedp.Navigate(website),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {

		fmt.Println("ERRROR: ", err)
		return nil, err
	}

	// fmt.Println("waiting for response from channel , ", website)
	statusCode:=  <- statusChanel

	// fmt.Println("status code: ", statusCode, "for website: ", website)

	// fmt.Println("received from channel, ", website)

	if statusCode < 200 || statusCode > 299 {
		fmt.Println("Error scraping website: ", website, "with status code: ", statusCode)
		return nil, errors.New("website did not respond correctly")
	}

	htmlContentReader := io.NopCloser(strings.NewReader(htmlContent))

	resp:= &http.Response{
		Body: htmlContentReader,
	}

	// fmt.Println("Finished scraping web , ", website)


	return resp, nil
}