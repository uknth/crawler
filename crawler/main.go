package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"git.sr.ht/~uknth/crawler"
)

var (
	file       = flag.String("file", "urls.txt", "File containing initial list of URls")
	depth      = flag.Int("depth", 3, "depth to which the application needs to crawl")
	count      = flag.Int("count", 4, "worker count")
	inactivity = flag.Int("inactivity", 15, "default time worker remain idle")
)

func urls(filePath string) ([]string, error) {
	var urls []string

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func main() {
	flag.Parse()

	uris, err := urls(*file)
	if err != nil {
		log.Fatal(err)
		return
	}

	var wg sync.WaitGroup

	cr := crawler.NewCrawler(
		*depth,
		*count,
		uris,
		time.Duration(*inactivity)*time.Second,
		&wg,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(cr.String())

	endch, err := cr.Crawl()
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	endch <- true
}
