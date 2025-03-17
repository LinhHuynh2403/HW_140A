package bug1

import "sync"

// Counter stores a count.
type Counter struct {
	mu sync.Mutex
	n  int64
}

// Inc increments the count in the Counter.
func (c *Counter) Inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}