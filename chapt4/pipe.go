package chapt4

import "fmt"

func Pipe() {
	generator := func(done chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, v := range integers {
				select {
				case <-done:
					return
				case intStream <- v:
				}
			}
		}()
		return intStream
	}
	operator := func(done chan interface{}, intStream <-chan int, op func(int) int) <-chan int {
		outputStream := make(chan int)
		go func() {
			defer close(outputStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case outputStream <- op(i):
				}
			}
		}()
		return outputStream
	}
	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)
	multiply := func(multiplier int, intStream <-chan int) <-chan int {
		return operator(done, intStream, func(i int) int { return i * multiplier })
	}
	add := func(additive int, intStream <-chan int) <-chan int {
		return operator(done, intStream, func(i int) int { return i + additive })
	}
	pipe := multiply(3, add(2, intStream))
	for v := range pipe {
		fmt.Println(v)
	}

}
