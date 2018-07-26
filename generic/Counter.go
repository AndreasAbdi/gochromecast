package generic

import "sync"

//Counter struct that concurrent incrementing integer.
type Counter struct {
	value int
	mutex sync.Mutex
}

//IncrementAndGet increments the counter and then immediately returns it
func (counter *Counter) IncrementAndGet() int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()
	counter.value++
	return counter.value
}

//Get returns the current value of the counter
func (counter *Counter) Get() int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()
	return counter.value
}
