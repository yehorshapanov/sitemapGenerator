package sitemap_generator

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type SitemapGenerator struct {
}

func (s SitemapGenerator) Crawl(url string, next_url chan string, r_token chan bool) {

	var documentBasePath = ""

	// Take a token (waits until a new token appears in a
	// empty channel)
	<-r_token
	// Give a token back at function exit
	defer func() {
		r_token <- true
	}()

	// Get the HTML for our url
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Print info
	fmt.Println("Visited", url)

	// Tokenize the HTML
	z := html.NewTokenizer(resp.Body)

	for {

		// If the HTML has ended, we break out of the loop
		token := z.Next()
		if token == html.ErrorToken {
			break
		}

		// New Token started
		if token == html.StartTagToken {

			// Check if the token is a <base> tag
			if name, _ := z.TagName(); string(name) == "base" {
				for {
					name, val, more := z.TagAttr()
					if string(name) == "href" {
						documentBasePath = string(val)
						break
					}

					if !more {
						break
					}
				}
			}

			// Check if the token is an <a> tag
			if name, _ := z.TagName(); string(name) == "a" {

				for {
					// Get the next attribute
					name, val, more := z.TagAttr()

					// Check if the attribute is "href"
					if string(name) == "href" {
						// Cast Url
						url = string(val)
						// Check if the URL is valid
						if !strings.HasPrefix(url, "http://") {
							if !strings.HasPrefix(url, "https://") {
								url = documentBasePath + url
							}
						}
						// The URL is valid so send it to the Url channel
						next_url <- url
					}

					// There are no more attributes so we break out of the
					// attribute search loop.
					if !more {
						break
					}

				}

			}

		}

	}

}

func (s SitemapGenerator) ParseArguments() (int, int, string, bool) {
	// Check for correct length
	if len(os.Args) != 4 {
		// Print usage
		fmt.Println("Invalid arguments!")
		fmt.Printf("Usage: %s $1 $2 $3\n", os.Args[0])
		fmt.Println(" $1: Number of concurrent processes")
		fmt.Println(" $2: Number of URLs to search for")
		fmt.Println(" $3: URL to start crawling on")
		fmt.Printf("\nExample: %s 100 1000 http://example.com\n", os.Args[0])
		return 0, 0, "", true
	}

	// Cast arguments
	first, _ := strconv.Atoi(os.Args[1])
	second, _ := strconv.Atoi(os.Args[2])
	third := os.Args[3]

	return first, second, third, false

}
