# Problem Statment

[![builds.sr.ht status](https://builds.sr.ht/~uknth/crawler/.build.yml.svg)](https://builds.sr.ht/~uknth/crawler/.build.yml?)

Given a list of URLS download and contents of URL and store it at a location. Parse the downloaded content for more URLs and keep on downloading till a given depth.

## Input

- initial list of URLs
- depth to which file needs to be downloaded


## Output

Downloaded contents of the URLs crawled

## Build

Ensure that dependencies are loaded 
- `go mod tidy`

Build the Binary
- ` go build -o crawler.bin git.sr.ht/~uknth/crawler/crawler`

This creates a binary `crawler.bin` in the working directory

## Run the Binary

`./crawler.bin -file="sample/urls.txt" -depth=4 -download=/tmp/crawler -count=5`

CLI parameters in the Binary

```
➜  $:(master) ./crawler.bin --help
Usage of ./crawler.bin:
  -count int
    	worker count (default 4)
  -depth int
    	depth to which the application needs to crawl (default 3)
  -download string
    	download location (default "/tmp/crawler")
  -file string
    	File containing initial list of URls (default "urls.txt")
  -inactivity int
    	default time worker remain idle (default 15)
```
