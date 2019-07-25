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

	Depth int

	Data string
}

type Result struct {
	Depth int
	Vals  []string
}

// TaskBuilder builds new tasks from given results
type TaskBuilder func(*Result) ([]Task, error)

func NewTaskBuilder(maxDepth int) TaskBuilder {
	var (
		taskID int

		processor = func(en string, depth int) Task {
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
				depth + 1,
				en,
			}
		}
	)

	return func(result *Result) ([]Task, error) {
		entries := result.Vals

		if len(entries) == 0 {
			return nil, errUndefinedTask
		}

		var tasks []Task

		if result.Depth >= maxDepth {
			return nil, errUndefinedTask
		}

		for _, en := range entries {
			tasks = append(
				tasks,
				processor(en, result.Depth),
			)
		}

		return tasks, nil
	}
}
