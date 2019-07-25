package crawler

import (
	"log"
	"sync"
	"time"
)

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
	result chan []string

	inactiveTimer *Timer

	wg *sync.WaitGroup

	// dispatcher
	dispatcher ExecutorDispatcher
}

func (dw *defaultWorker) Start() {
	log.Printf("Starting Worker: %d\n", dw.id)
	dw.wg.Add(1)
	go func() {
		for {
			dw.cont <- dw.work
			select {
			case task := <-dw.work:
				dw.inactiveTimer.Reset()

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

				dw.result <- result
			case <-dw.end:
				return
			}
		}
	}()
}

func (dw *defaultWorker) janitor(end chan bool) {
	go func() {
		for {
			select {
			case <-end:
				dw.wg.Done()
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
	inactivity time.Duration,
	wg *sync.WaitGroup,
) Worker {
	ender := make(chan bool)

	dw := &defaultWorker{
		id:            id,
		cont:          contc,
		work:          make(chan Task),
		end:           make(chan bool),
		result:        result,
		inactiveTimer: NewTimer(inactivity, ender),
		wg:            wg,
		dispatcher:    NewExecutorDispatcher(),
	}

	dw.janitor(ender)
	return dw
}
