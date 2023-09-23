package chapt5

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func ReqReplica() {
	doWork := func(done <-chan interface{}, id int, wg *sync.WaitGroup, result chan<- int) {
		var simulatedLoadTime time.Duration
		defer func(start time.Time) {
			if time.Since(start) < simulatedLoadTime {
				fmt.Printf("%v took %v\n", id, simulatedLoadTime)
			} else {
				fmt.Printf("%v took %v\n", id, time.Since(start))
			}
			wg.Done()
		}(time.Now())

		simulatedLoadTime = time.Duration(1+rand.Intn(5)) * time.Second
		select {
		case <-done:
		case <-time.After(simulatedLoadTime):
		}
		select {
		case <-done:
		case result <- id:
		}
	}

	done := make(chan interface{})
	result := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go doWork(done, i, wg, result)
	}
	firstReturned := <-result
	close(done)
	wg.Wait()
	fmt.Printf("Received an answer from %v\n", firstReturned)
}
