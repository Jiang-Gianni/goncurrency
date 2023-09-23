package chapt4

import "fmt"

func BridgeChannel(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	output := make(chan interface{})
	go func() {
		defer close(output)
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range OrDone(done, stream) {
				select {
				case output <- val:
				case <-done:
				}
			}
		}
	}()
	return output
}

func ExampleBridgeChannel() {
	genVals := func() <-chan <-chan interface{} {
		output := make(chan (<-chan interface{}))
		go func() {
			defer close(output)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				output <- stream
			}
		}()
		return output
	}

	for v := range BridgeChannel(nil, genVals()) {
		fmt.Printf("%v", v)
	}
}
