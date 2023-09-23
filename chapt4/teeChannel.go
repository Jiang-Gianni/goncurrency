package chapt4

import "fmt"

func TeeChannel(done, input <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range OrDone(done, input) {
			var out1, out2 = out1, out2
			// 2 iteration: once per output channel
			// nil after pushing the value so that the local copy of the channel it is blocked on next iteration
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

func ExampleTeeChannel() {
	done := make(chan interface{})
	defer close(done)
	out1, out2 := TeeChannel(done, Take(done, Repeat(done, 1, 2), 3))
	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
