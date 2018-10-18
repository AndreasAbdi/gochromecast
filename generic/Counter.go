package generic

import "sync"

//Counter struct that concurrent incrementing integer.
type Counter struct {
	value int
	mutex sync.Mutex
}

//GetAndIncrement increments the counter and then immediately returns it
func (counter *Counter) GetAndIncrement() int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()
	val := counter.value
	counter.value++
	return val
}

//Get returns the current value of the counter
func (counter *Counter) Get() int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()
	return counter.value
}

//Reset sets the counter back to 0
func (counter *Counter) Reset() {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()
	counter.value = 0
}

//ChCounter is an atomic counter that requires users to get the value via subscribing to an outputs channel.
type ChCounter struct {
	value     int
	Outputs   chan int
	increment chan bool
}
