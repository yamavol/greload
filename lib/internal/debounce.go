package internal

import (
	"sync"
	"time"
)

type debouncer struct {
	mu       sync.Mutex
	interval time.Duration
	timer    *time.Timer
}

// NewDebouncer returns a function debouncer that delays execution.
// The second returned function is to cancel the pending exceution.
func NewDebouncer(interval time.Duration) (func(fn func()), func()) {
	d := &debouncer{
		interval: interval,
	}
	return func(fn func()) {
		d.add(fn)
	}, d.cancel
}

func (d *debouncer) add(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.interval, fn)
}

func (d *debouncer) cancel() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
