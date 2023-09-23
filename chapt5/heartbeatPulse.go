package chapt5

import (
	"fmt"
	"math/rand"
)

// Heartbeat pulse is also useful for tests: if a value is read from the heartbeat channel then we can ensure that the goroutine process has started

func HeartbeatPulse() {
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		heartbeatStream := make(chan interface{}, 1)
		intStream := make(chan int)

		go func() {
			defer close(heartbeatStream)
			defer close(intStream)

			for i := 0; i < 10; i++ {
				select {
				case heartbeatStream <- struct{}{}:
				default:
				}

				select {
				case <-done:
					return
				case intStream <- rand.Intn(10):
				}
			}

		}()
		return heartbeatStream, intStream
	}

	done := make(chan interface{})
	defer close(done)
	heartbeatStream, intStream := doWork(done)
	for {
		select {
		case _, ok := <-heartbeatStream:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-intStream:
			if ok {
				fmt.Printf("value %v\n", r)
			} else {
				return
			}
		}
	}

}
