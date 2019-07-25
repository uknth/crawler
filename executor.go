package crawler

import (
	"errors"
	"log"
)

var (
	errExecutorNotFound = errors.New("executor not found")
)

// Executor defines the task executor
type Executor interface {
	Execute(task Task) ([]string, error)
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

type downloadExecutor struct{}

func (de *downloadExecutor) Execute(task Task) ([]string, error) {
	log.Println("Download Task Received:", task.ID, task.Type, task.Data)
	return []string{}, nil
}

type parseExecutor struct{}

func (pe *parseExecutor) Execute(task Task) ([]string, error) {
	log.Println("Parse Task Received:", task.ID, task.Type, task.Data)
	return []string{"http://new_url.com"}, nil
}

// NewDownloadExecutor returns an executor which performs the download task
func NewDownloadExecutor() (Executor, error) {
	return &downloadExecutor{}, nil
}

// NewParseExecutor returns an executor which performs the parsing task
func NewParseExecutor() (Executor, error) {
	return &parseExecutor{}, nil
}
