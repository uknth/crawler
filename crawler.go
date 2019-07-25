package crawler

import (
	"fmt"
	"sync"
	"time"
)

// Crawler ...
type Crawler struct {
	// initial list of urls
	URLs []string

	// Depth to which we need to crawl
	depth int

	// Download Location
	download string

	dispatcher  Dispatcher
	taskBuilder TaskBuilder

	wg *sync.WaitGroup
}

func (c *Crawler) String() string {
	return fmt.Sprintf("Depth: {%d} Loc: {%s} URLs: {%v}", c.depth, c.download, c.URLs)
}

// Crawl crawls the given URL and saves the downloaded file on
// given location
func (c *Crawler) Crawl() (chan bool, error) {
	collector, err := c.dispatcher.Dispatch()
	if err != nil {
		return nil, err
	}

	tasks, err := c.taskBuilder(c.URLs)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		collector.Work <- task
	}
	return collector.End, nil
}

// NewCrawler returns a new Crawler object
func NewCrawler(
	depth int,
	download string,
	wc int,
	urls []string,
	inactivity time.Duration,
	wg *sync.WaitGroup,
) Crawler {
	return Crawler{
		urls,
		depth,
		download,
		NewDispatcher(wc, inactivity, wg),
		NewTaskBuilder(),
		wg,
	}
}
