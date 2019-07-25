package crawler

import "log"

// Worker takes a task and executes it
type Worker interface {
	// Starts the Worker
	Start()
	// Stops the Worker
	Stop()
}

type defaultWorker struct {
	id int
	// controls the availability of the worker
	cont chan chan Task
	// channel which recieves the actual task
	work chan Task
	// terminates the worker loop
	end chan bool
	// results is channel where we return the output of execution
	results chan []string

	// dispatcher
	dispatcher ExecutorDispatcher
}

func (dw *defaultWorker) Start() {
	log.Printf("Starting Worker: %d\n", dw.id)
	go func() {
		for {
			dw.cont <- dw.work
			select {
			case task := <-dw.work:
				exec, err := dw.dispatcher(task)
				if err != nil {
					log.Println("DISPATCHER ERROR:", err.Error())
					continue
				}

				result, err := exec.Execute(task)
				if err != nil {
					log.Println("TASK ERROR: ", err.Error())
					continue
				}
				dw.results <- result
			case <-dw.end:
				return
			}
		}
	}()
}

func (dw *defaultWorker) Stop() {
	log.Printf("Stopping Worker: [%d]", dw.id)
	dw.end <- true
}

// NewDefaultWorker returns the default implementation of Worker
func NewDefaultWorker(
	id int,
	contc chan chan Task,
	result chan []string,
) Worker {
	return &defaultWorker{
		id:         id,
		cont:       contc,
		work:       make(chan Task),
		end:        make(chan bool),
		dispatcher: NewExecutorDispatcher(),
	}
}
