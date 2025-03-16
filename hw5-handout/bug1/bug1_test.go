package bug1

import (
	"sync"
	"testing"
)

func sumUp(n1, n2 int64) Counter {
	var wg sync.WaitGroup
	counter := Counter{n: 0}
	for i := int64(0); i < n1; i++ {
		wg.Add(1)
		go func() {
			for i := int64(0); i < n2; i++ {
				counter.Inc()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return counter
}

func TestBug1(t *testing.T) {
	tests := []struct {
		n1   int64
		n2   int64
		want int64
	}{
		{100, 1000, 100 * 1000},
		{500, 4000, 500 * 4000},
	}
	for _, test := range tests {
		got := sumUp(test.n1, test.n2)
		if test.want != got.n {
			t.Errorf("bug1 failed on (%d, %d); want %d, got %d", test.n1, test.n2, test.want, got.n)
		}
	}
}
