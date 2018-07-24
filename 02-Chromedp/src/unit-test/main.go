// Command simple is a chromedp example demonstrating how to do a simple google
// search.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

var cliOpts = []runner.CommandLineOption{
	runner.NoDefaultBrowserCheck,
	runner.NoFirstRun,
	runner.Flag("headless", true),
}
var testDataDir string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v", err)
	}
	testDataDir = "file://" + path.Join(wd, "testdata")
}

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt,
		chromedp.WithRunnerOptions(cliOpts...),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(">>> test input")
	var res string
	if err := c.Run(ctxt, testInput("1", &res)); err == nil {
		if res != "chromedp" {
			log.Fatalf(">>> expected input value to be 'chromedp', got %s", res)
		}
		log.Printf(">>> test input success")
	}

	log.Printf(">>> test input with value")
	text := "GoJakarta"
	if err := c.Run(ctxt, testInputWithValue(text, &res)); err == nil {
		if res != text {
			log.Fatalf(">>> expected input value to be '%s', got %s", text, res)
		}
		log.Printf(">>> test input with value success")
	}

	log.Printf(">>> test reset")
	if err := c.Run(ctxt, testReset(&res)); err == nil {
		if res != "chromedp" {
			log.Fatalf(">>> expected input value to be 'chromedp', got %s", res)
		}
		log.Printf(">>> test reset success")
	}

	log.Printf(">>> test submit")
	var ok bool
	if err := c.Run(ctxt, testSubmit(&res, &ok)); err == nil {
		if res != "Brankas - Easy Money Management" {
			log.Fatalf(">>> expected input value to be 'Easy Money Management', got %s", res)
		}
		log.Printf(">>> test submit success")
	}

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

func testInput(id string, res *string) chromedp.Tasks {
	sel := fmt.Sprintf(`#input%s`, id)
	return chromedp.Tasks{
		chromedp.Navigate(testDataDir + "/index.html"),
		chromedp.Value(sel, res, chromedp.ByID),
	}
}

func testInputWithValue(text string, res *string) chromedp.Tasks {
	sel := `#input1`
	return chromedp.Tasks{
		chromedp.SetValue(sel, text, chromedp.ByID),
		chromedp.Value(sel, res, chromedp.ByID),
	}
}

func testReset(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Click(`#btn1`, chromedp.ByID),
		chromedp.Value(`#input1`, res, chromedp.ByID),
	}
}

func testSubmit(res *string, ok *bool) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Click(`#btn2`, chromedp.ByID),
		chromedp.AttributeValue(`icon-brankas`, "alt", res, ok, chromedp.ByID),
	}
}
