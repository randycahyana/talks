// Command simple is a chromedp example demonstrating how to do a simple google
// search.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

var cliOpts = []runner.CommandLineOption{
	runner.NoDefaultBrowserCheck,
	runner.NoFirstRun,
	runner.Flag("headless", true),
}

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt,
		chromedp.WithLog(log.Printf),
		chromedp.WithRunnerOptions(cliOpts...),
	)
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var site, res string
	err = c.Run(ctxt, googleSearch("site:brank.as", "Home", &site, &res))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Brankas: %s", res)

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func googleSearch(q, text string, site, res *string) chromedp.Tasks {
	sel := fmt.Sprintf(`//a[text()[contains(., '%s')]]`, text)
	return chromedp.Tasks{
		chromedp.Navigate(`https://www.google.com`),
		chromedp.WaitVisible(`#hplogo`, chromedp.ByID),
		chromedp.SendKeys(`#lst-ib`, q+"\n", chromedp.ByID),
		chromedp.WaitVisible(`#res`, chromedp.ByID),
		chromedp.Click(sel),
		chromedp.Text(`/html/body/div[2]/div/div[2]/section[1]/div[2]/div[1]/h1`, res),
	}
}
