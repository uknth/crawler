package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"git.sr.ht/~uknth/crawler"
)

var (
	file     = flag.String("file", "urls.txt", "File containing initial list of URls")
	depth    = flag.Int("depth", 3, "depth to which the application needs to crawl")
	download = flag.String("download", "/tmp/crawler", "location to download URL contents")
	count    = flag.Int("count", 4, "worker count")
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

	cr := crawler.NewCrawler(
		*depth, *download, *count, uris,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(cr.String())

	err = cr.Crawl()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Duration(3) * time.Second)
}
