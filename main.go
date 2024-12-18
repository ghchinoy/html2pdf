package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	pageURL   string
	outputDir string
)

func init() {
	flag.StringVar(&pageURL, "url", "", "webpage to PDF")
	flag.StringVar(&outputDir, "dir", ".", "output directory")
	flag.Parse()
}

func main() {
	// guard against blank page URL
	if pageURL == "" {
		log.Print("requires a --url")
		os.Exit(1)
	}
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// get an html page and create pdf bytes
	var title string
	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL), // navigate to a page,
		chromedp.Title(&title),     // get the title,
		printToPDF(pageURL, &buf),  // obtain a pdf of the page
	)
	if err != nil {
		log.Fatal(err)
	}

	// create a filename from the title
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, " ", "")
	if len(title) > 15 {
		title = string([]rune(title)[:15])
	}
	filename := fmt.Sprintf("%s.pdf", title)

	// write the pdf output file
	fullPath := filepath.Join(outputDir, filename)
	if err := os.WriteFile(fullPath, buf, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s\n", fullPath)
}

// printToPDF print a specific pdf page.
func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
