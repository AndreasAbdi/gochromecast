package generic

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterSerial(t *testing.T) {
	counter := Counter{}
	repeats := 10
	for i := 0; i < repeats; i++ {
		assert.Equal(t, counter.GetAndIncrement(), i, "Should be equal")
	}
}

func TestCounterSerialReset(t *testing.T) {
	counter := Counter{}
	repeats := 10
	for i := 0; i < repeats; i++ {
		counter.GetAndIncrement()
	}
	counter.Reset()
	assert.Equal(t, counter.GetAndIncrement(), 0, "Should be equal")
}

func TestCounterConcurrent(t *testing.T) {
	counter := Counter{}
	repeats := 10
	allHit := sync.Map{}
	wg := sync.WaitGroup{}
	wg.Add(repeats)
	for i := 0; i < repeats; i++ {
		go func() {
			defer wg.Done()
			val := counter.GetAndIncrement()
			t.Log(val)
			allHit.Store(val, true)
		}()
	}
	wg.Wait()
	for i := 0; i < repeats; i++ {
		if _, ok := allHit.Load(i); !ok {
			t.Error("Not all values were generated", i)
		}
	}
}
