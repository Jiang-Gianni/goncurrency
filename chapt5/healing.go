package chapt5

import (
	"log"
	"os"
	"time"

	"github.com/Jiang-Gianni/goncurrency/chapt4"
)

type startGoroutineFn func(
	done <-chan interface{},
	pulseInterval time.Duration,
) (heartbeat <-chan interface{})

func NewSteward(timeout time.Duration, startGoroutineFn startGoroutineFn) startGoroutineFn {
	return func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) <-chan interface{} {
		heartbeat := make(chan interface{})

		go func() {
			defer close(heartbeat)
			var wardDone chan interface{}
			var wardHeartbeat <-chan interface{}
			startWard := func() {
				wardDone = make(chan interface{})
				wardHeartbeat = startGoroutineFn(chapt4.OrDone(wardDone, done), timeout/2)
			}
			startWard()
			pulse := time.NewTicker(pulseInterval)

		monitorLoop:
			for {
				timeoutSignal := time.After(timeout)
				for {
					select {
					case <-pulse.C:
						select {
						case heartbeat <- struct{}{}:
						default:
						}
					case <-wardHeartbeat:
						continue monitorLoop
					case <-timeoutSignal:
						log.Println("steward: ward unhealthy; restarting")
						close(wardDone)
						startWard()
						continue monitorLoop
					case <-done:
						return
					}
				}
			}
		}()

		return heartbeat
	}
}

func ExampleHealing() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWork := func(done <-chan interface{}, _ time.Duration) <-chan interface{} {
		log.Println("ward: Hello, I'm irresponsible!")
		go func() {
			<-done
			log.Println("ward: I'm halting")
		}()
		return nil
	}

	doWorkWithSteward := NewSteward(4*time.Second, doWork)
	done := make(chan interface{})
	time.AfterFunc(9*time.Second, func() {
		log.Println("main: halting steward and ward")
		close(done)
	})
	for range doWorkWithSteward(done, 4*time.Second) {
	}
	log.Println("Done")
}
