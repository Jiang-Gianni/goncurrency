package chapt4

import "sync"

func FanInWG(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	wg := &sync.WaitGroup{}
	output := make(chan interface{})

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case output <- <-c:
			}
		}
	}

	wg.Add(len(channels))

	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		defer close(output)
		wg.Wait()
	}()

	return output
}
