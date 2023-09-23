package chapt4

func OrDone(done, c <-chan interface{}) <-chan interface{} {
	output := make(chan interface{})
	go func() {
		defer close(output)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				// If c has been closed then return
				if !ok {
					return
				}
				// Select in case the c chan's reader is blocked and return if there is a signal from the done chan
				select {
				case output <- v:
				case <-done:
					return
				}
			}
		}
	}()
	return output
}
