package crawler

import "fmt"

// Crawler ...
type Crawler struct {
	// initial list of urls
	URLs []string

	// Depth to which we need to crawl
	depth int

	// Download Location
	download string
}

func (c *Crawler) String() string {
	return fmt.Sprintf("Depth: {%d} Loc: {%s} URLs: {%v}", c.depth, c.download, c.URLs)
}

// Crawl crawls the given URL and saves the downloaded file on
// given location
func (c *Crawler) Crawl() error {
	return nil
}

// NewCrawler returns a new Crawler object
func NewCrawler(
	depth int, 
	download string, 
	urls []string,
) (*Crawler, error) {
	return &Crawler{urls, depth, download}, nil
}
