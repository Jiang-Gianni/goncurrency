package chapt4

import (
	"fmt"
	"math/rand"
)

func RepeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	fnChan := make(chan interface{})
	go func() {
		defer close(fnChan)
		for {
			select {
			case <-done:
				return
			case fnChan <- fn():
			}
		}
	}()
	return fnChan
}

func ExampleRepeatFn() {
	done := make(chan interface{})
	fn := func() interface{} { return rand.Int() }
	fnChan := RepeatFn(done, fn)
	for num := range Take(done, fnChan, 10) {
		fmt.Println(num.(int))
	}
}
