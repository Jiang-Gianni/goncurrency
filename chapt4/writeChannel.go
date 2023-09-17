package chapt4

import (
	"fmt"
	"math/rand"
	"time"
)

func WriteChannel() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("exiting newRandStream")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}
	done := make(chan interface{})
	randStream := newRandStream(done)
	for i := 0; i < 3; i++ {
		fmt.Println(i, <-randStream)
	}
	close(done)
	time.Sleep(1 * time.Second)
}
