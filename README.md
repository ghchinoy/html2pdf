# html2pdf

A small example of using chromedp to create a PDF from an HTML web page, given an URL, get a PDF file.

For use in applications like [Fabulae](https://github.com/ghchinoy/fabulae).

## flags

* `url` (required) HTTP url to PDF
* `dir` (optional) output directory, defaults to "." 


## example use

```
go run . --url https://cloud.google.com/transform/gen-ai-kpis-measuring-ai-success-deep-dive --dir samples
```
