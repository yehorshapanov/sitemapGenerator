package main

import (
	"fmt"
	"sitemapGenerator/sitemap_generator"
)

func main() {
	s := sitemap_generator.SitemapGenerator{}
	// Parse arguments
	num_routines, num_urls, start_url, err := s.ParseArguments()
	if err {
		return
	}

	// Make buffered channel with strings for the urls
	next_url := make(chan string, num_routines)
	next_url <- start_url

	// Create buffered channel for tokens
	r_tokens := make(chan bool, num_routines)
	for i := 0; i < num_routines; i++ {
		r_tokens <- true
	}

	// Create map (string -> bool) to check if a url has
	// been visited
	m := map[string]bool{}

	for i := 0; i < num_urls; {
		// Get next url
		url := <-next_url

		// Check if the url has been visited
		_, found := m[url]

		// If yes, go to the next url
		if found {
			continue
		}

		// If not, add the url to the map, increment
		// the counter, insert it into the database
		// and start new crawl goroutine with that
		// url.
		m[url] = true
		i++
		fmt.Println(url)
		go s.Crawl(url, next_url, r_tokens)
	}

}
