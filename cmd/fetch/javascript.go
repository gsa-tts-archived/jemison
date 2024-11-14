package main

import (
	"context"
	"log"
	"net/url"

	kv "github.com/GSA-TTS/jemison/internal/kv"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

// WARNING
// This doesn't work. It was left in the tree for reference.
// It may want to come out until such time as it wants to be tackled.
func extractJavascriptHTML(obj kv.Object) string {

	// initialize a chrome instance
	// ctx, cancel := chromedp.NewContext(
	// 	context.Background(),
	// 	chromedp.WithLogf(zap.L()), /*log.Printf*/
	// )
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	u := url.URL{
		Scheme: obj.GetValue("scheme"),
		Host:   obj.GetValue("host"),
		Path:   obj.GetValue("path"),
	}

	// https://github.com/chromedp/examples/blob/master/text/main.go
	// https://pkg.go.dev/github.com/chromedp/chromedp#section-readme
	// https://stackoverflow.com/questions/76359161/whats-the-best-way-to-get-just-the-user-readble-word-content-of-a-page
	var res []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(u.String()),
		chromedp.WaitReady("body"),
		chromedp.Nodes("//p[text()] | //li[text()] | //div[text()] | //td[text()]", &res))

	if err != nil {
		log.Fatal(err)
	}

	content := ""
	for _, item := range res {

		var innerHTML string
		chromedp.Run(ctx,
			chromedp.InnerHTML(item.FullXPath(), &innerHTML),
		)

		zap.L().Debug("innerHTML from chromedp", zap.String("innerHTML", innerHTML))

		content += innerHTML
	}

	return content
}
