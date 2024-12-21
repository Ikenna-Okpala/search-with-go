package webcrawler

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/chromedp/chromedp"
)

func connectToWebsite(website string) (*http.Response, error){

	ctx:= context.Background()

	ctx, cancel:= chromedp.NewContext(ctx)

	defer cancel()

	var htmlContent string

	err:= chromedp.Run(
		ctx,
		chromedp.Navigate(website),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return nil, err
	}

	htmlContentReader := io.NopCloser(strings.NewReader(htmlContent))

	resp:= &http.Response{
		Body: htmlContentReader,
	}


	return resp, nil
}