package crawler

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type Collector struct {
	Work chan Task
	End  chan bool
}

// Dispatcher takes the initial tasks and dispatches it to
// the workers. It also collects the results generated by the
// workers and generates the tasks that
type Dispatcher struct {
	workers     []Worker
	taskbuilder TaskBuilder

	control chan chan Task

	// task controls
	task chan Task
	rslt chan *Result
	end  chan bool
	wg   *sync.WaitGroup

	buffer *list.List
	tkr    *time.Ticker
}

func (d *Dispatcher) dispatcher(end chan bool) {
	go func() {
		for {
			select {
			case <-end:
				for _, w := range d.workers {
					w.Stop()
				}
				return
			case task := <-d.task:
				// wait for worker to be available
				w := <-d.control
				w <- task
			}
		}
	}()
}

func (d *Dispatcher) collector(end chan bool) {
	go func() {
		for {
			select {
			case r := <-d.rslt:
				tasks, err := d.taskbuilder(r)
				if err != nil && err == errUndefinedTask {
					continue
				}

				if err != nil && err != errUndefinedTask {
					log.Println("ERROR IN BUILDING TASK:", err.Error())
					continue
				}

				for _, t := range tasks {
					d.buffer.PushBack(t)
				}
			case <-end:
				return
			}
		}
	}()
}

func (d *Dispatcher) emitter() {
	go func() {
		for {
			select {
			case <-d.tkr.C:
				for e := d.buffer.Front(); e != nil; e = e.Next() {
					task := e.Value.(Task)
					d.task <- task
					d.buffer.Remove(e)
				}
			}
		}
	}()
}

func (d *Dispatcher) janitor(ends ...chan bool) {
	go func() {
		<-d.end

		for _, en := range ends {
			en <- true
		}

		return
	}()
}

// Dispatch starts the dispatcher
func (d *Dispatcher) Dispatch() (Collector, error) {
	var (
		cend chan bool
		dend chan bool
	)

	for _, worker := range d.workers {
		worker.Start()
	}

	d.janitor(cend, dend)
	d.collector(cend)
	d.dispatcher(dend)
	d.emitter()

	return Collector{d.task, d.end}, nil
}

// NewDispatcher returns a dispatcher
func NewDispatcher(wc int, depth int, inactivity time.Duration, wg *sync.WaitGroup) Dispatcher {
	var (
		workers []Worker

		control = make(chan chan Task)
		result  = make(chan *Result)
	)

	for idx := 1; idx <= wc; idx++ {
		worker := NewDefaultWorker(idx, control, result, inactivity, wg)
		workers = append(workers, worker)
	}

	return Dispatcher{
		workers:     workers,
		taskbuilder: NewTaskBuilder(depth),
		control:     control,
		task:        make(chan Task),
		rslt:        result,
		end:         make(chan bool),
		buffer:      list.New(),
		tkr:         time.NewTicker(2 * time.Second),
	}
}
