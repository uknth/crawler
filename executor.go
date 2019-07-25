package crawler

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	errExecutorNotFound = errors.New("executor not found")

	download = flag.String("download", "/tmp/crawler", "download location")
)

// Executor defines the task executor
type Executor interface {
	Execute(task Task) (*Result, error)
}

// ExecutorDispatcher returns dispatcher based on the Task type
type ExecutorDispatcher func(task Task) (Executor, error)

// NewExecutorDispatcher returns a Dispatcher which redirects a task
// to a right Executor
func NewExecutorDispatcher() ExecutorDispatcher {
	return func(task Task) (Executor, error) {
		switch task.Type {
		case "download":
			return NewDownloadExecutor()
		case "parse":
			return NewParseExecutor()
		default:
			return nil, errExecutorNotFound
		}
	}
}

type downloadExecutor struct {
	client *http.Client

	dir string
}

func (de *downloadExecutor) dial(req *http.Request) (*http.Response, error) {
	res, err := de.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (de *downloadExecutor) url(unparsed string) (string, error) {
	_, err := url.Parse(unparsed)
	if err != nil {
		return *new(string), err
	}
	return unparsed, nil
}

func (de *downloadExecutor) newRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	return req, err
}

func (de *downloadExecutor) generateFile(taskID int) string {
	return de.dir + strconv.Itoa(taskID) + ".html"
}

func (de *downloadExecutor) Execute(task Task) (*Result, error) {
	log.Println("Download Task Received:", task.ID, task.Type, task.Depth, task.Data)
	var results []string

	u, err := de.url(task.Data)
	if err != nil {
		return nil, err
	}

	req, err := de.newRequest(u)
	if err != nil {
		return nil, err
	}

	res, err := de.dial(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fn := de.generateFile(task.ID)
	err = ioutil.WriteFile(
		fn,
		bts,
		0644,
	)
	if err != nil {
		return nil, err
	}
	results = append(results, fn)
	return &Result{task.Depth, results}, nil
}

type parseExecutor struct{}

func (pe *parseExecutor) validateFile(fn string) error {
	_, err := os.Stat(fn)
	if err != nil {
		return err
	}
	return nil
}

func (pe *parseExecutor) queryDocument(fn string) (*goquery.Document, error) {
	fl, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(fl)
}

func (pe *parseExecutor) Execute(task Task) (*Result, error) {
	log.Println("Parse Task Received:", task.ID, task.Type, task.Depth, task.Data)
	var results []string

	err := pe.validateFile(task.Data)
	if err != nil {
		return nil, err
	}

	doc, err := pe.queryDocument(task.Data)
	if err != nil {
		return nil, err
	}

	doc.Find("a").Each(func(i int, sl *goquery.Selection) {
		ln, exists := sl.Attr("href")
		if !exists {
			return
		}
		if strings.HasPrefix(ln, "http") {
			results = append(results, ln)
		}
	})

	return &Result{task.Depth, results}, nil
}

// NewDownloadExecutor returns an executor which performs the download task
func NewDownloadExecutor() (Executor, error) {
	var dir = *download

	if !strings.HasSuffix(*download, "/") {
		dir = dir + "/"
	}

	return &downloadExecutor{
		client: &http.Client{},
		dir:    dir,
	}, nil
}

// NewParseExecutor returns an executor which performs the parsing task
func NewParseExecutor() (Executor, error) {
	return &parseExecutor{}, nil
}
