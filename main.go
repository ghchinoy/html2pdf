package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	pageURL string
)

func init() {
	flag.StringVar(&pageURL, "url", "", "webpage to PDF")
	flag.Parse()
}

func main() {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var title string
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL),
		chromedp.Title(&title),
	)
	if err != nil {
		log.Fatal(err)
	}

	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, " ", "")

	if len(title) > 15 {
		title = string([]rune(title)[:15])
	}

	// capture pdf
	var buf []byte
	if err := chromedp.Run(ctx, printToPDF(pageURL, &buf)); err != nil {
		log.Fatal(err)
	}

	filename := fmt.Sprintf("%s.pdf", title)
	if err := os.WriteFile(filename, buf, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s\n", filename)
}

// print a specific pdf page.
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
