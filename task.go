package crawler

import (
	"errors"
	"strings"
)

var errUndefinedTask = errors.New("task undefined")

// Task defines the task to be performed
// by the crawler, in our case it can be
// download or parse
type Task struct {
	ID int

	Type string

	Data string
}

// TaskBuilder builds new tasks from given results
type TaskBuilder func([]string) ([]Task, error)

func NewTaskBuilder() TaskBuilder {
	var (
		taskID int

		processor = func(en string) Task {
			var taskType string

			switch {
			case strings.HasPrefix(en, "http"):
				taskType = "download"
			case strings.HasPrefix(en, "/"):
				taskType = "parse"
			default:
				taskType = "undefined"
			}

			taskID = taskID + 1

			return Task{
				taskID,
				taskType,
				en,
			}
		}
	)

	return func(entries []string) ([]Task, error) {
		if len(entries) == 0 {
			return nil, errUndefinedTask
		}

		var tasks []Task

		for _, en := range entries {
			tasks = append(
				tasks,
				processor(en),
			)
		}

		return tasks, nil
	}
}
