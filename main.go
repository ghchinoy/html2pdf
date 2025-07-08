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
	// pageURL is the url to be converted to PDF
	pageURL string
	// outputDir is the directory where the PDF will be saved
	outputDir string
)

// init initializes the command-line flags.
func init() {
	flag.StringVar(&pageURL, "url", "", "webpage to PDF")
	flag.StringVar(&outputDir, "dir", ".", "output directory")
	flag.Parse()
}

// main is the entry point of the program.
// It parses the command-line flags, calls the html2pdf function to generate the PDF,
// and writes the PDF to a file.
func main() {
	// guard against blank page URL
	if pageURL == "" {
		log.Print("requires a --url")
		os.Exit(1)
	}

	buf, filename, err := html2pdf(pageURL)
	if err != nil {
		log.Fatal(err)
	}

	// write the pdf output file
	fullPath := filepath.Join(outputDir, filename)
	if err := os.WriteFile(fullPath, buf, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s\n", fullPath)
}

// html2pdf converts a given URL to a PDF.
// It uses chromedp to navigate to the URL, get the title, and print the page to PDF.
// It returns the PDF content as a byte slice, the generated filename, and an error if any.
func html2pdf(pageURL string) ([]byte, string, error) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// get an html page and create pdf bytes
	var title string
	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL),
		chromedp.Title(&title),
		printToPDF(&buf),
	)
	if err != nil {
		return nil, "", err
	}

	// create a filename from the title
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, " ", "")
	if len(title) > 15 {
		title = string([]rune(title)[:15])
	}
	filename := fmt.Sprintf("%s.pdf", title)

	return buf, filename, nil
}

// printToPDF prints the current page to PDF.
// It takes a byte slice pointer as input and populates it with the PDF content.
func printToPDF(res *[]byte) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
		if err != nil {
			return err
		}
		*res = buf
		return nil
	})
}