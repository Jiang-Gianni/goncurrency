package chapt5

import (
	"fmt"
	"time"
)

func Heartbeat() {
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)

		go func() {
			defer close(heartbeat)
			defer close(results)
			pulse := time.NewTicker(pulseInterval)
			workGen := time.NewTicker(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}

			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return
					case <-pulse.C:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for {
				select {
				case <-done:
					return
				case <-pulse.C:
					sendPulse()
				case r := <-workGen.C:
					sendResult(r)
				}
			}
		}()

		return heartbeat, results
	}

	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, result := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-result:
			if !ok {
				return
			}
			fmt.Printf("result %v\n", r.Second())
			// case <-time.After(timeout):
			// 	return
		}
	}
}
