package chapt3

import (
	"fmt"
	"sync"
)

// ouptup: "goodday" printed thrice
// salutation is assigned to the next string and each goroutine refers to that variable: the runtime is looping over before any goroutine ends executing
func ForLoop1() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()
}

// ouptup: "good day, hello, greetings"
// each goroutines takes an input value -> a copy of salutation is passed
func ForLoop2() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(value string) {
			defer wg.Done()
			fmt.Println(value)
		}(salutation)
	}
	wg.Wait()
}
