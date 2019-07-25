package crawler

import (
	"log"
	"sync"
	"time"
)

type Timer struct {
	timestamp time.Time

	end chan bool

	mu sync.Mutex
}

func (t *Timer) Start() {
	t.mu.Lock()
	t.timestamp = time.Now()
	t.mu.Unlock()
}

func (t *Timer) Reset() {
	t.mu.Lock()
	t.timestamp = time.Now()
	t.mu.Unlock()
}

func (t *Timer) control(inactivity time.Duration) {
	for {
		t.mu.Lock()
		if time.Now().After(t.timestamp.Add(inactivity)) {
			log.Println("Inactive for :", inactivity.Seconds())
			t.end <- true
		}
		t.mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func NewTimer(inactivity time.Duration, end chan bool) *Timer {
	timer := &Timer{
		end: end,
	}

	go timer.control(inactivity)

	timer.Start()
	return timer
}
