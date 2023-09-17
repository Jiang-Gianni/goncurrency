package chapt4

import (
	"fmt"
	"time"
)

func ReadChannel() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}
	done := make(chan interface{})
	strings := make(chan string)
	terminated := doWork(done, strings)

	strings <- "Hello there"

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine")
		close(done)
	}()

	<-terminated

	fmt.Println("Done")
}
