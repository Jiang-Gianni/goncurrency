package chapt4

import (
	"fmt"
	"time"
)

func OrChannel() {
	var or func(chans ...<-chan interface{}) <-chan interface{}
	or = func(chans ...<-chan interface{}) <-chan interface{} {
		switch len(chans) {
		case 0:
			return nil
		case 1:
			return chans[0]
		}
		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			switch len(chans) {
			case 2:
				select {
				case <-chans[0]:
				case <-chans[1]:
				}
			default:
				select {
				case <-chans[0]:
				case <-chans[1]:
				case <-chans[2]:
				case <-or(append(chans[3:], orDone)...):
				}
			}
		}()
		return orDone
	}
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}
