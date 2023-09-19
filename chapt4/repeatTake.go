package chapt4

import "fmt"

func Repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueChan := make(chan interface{})
	go func() {
		defer close(valueChan)
		for {
			for _, value := range values {
				select {
				case <-done:
					return
				case valueChan <- value:
				}
			}
		}
	}()
	return valueChan
}

func Take(done <-chan interface{}, inputChan <-chan interface{}, qty int) <-chan interface{} {
	takeChan := make(chan interface{})
	go func() {
		defer close(takeChan)
		for i := 0; i < qty; i++ {
			select {
			case <-done:
				return
			case takeChan <- <-inputChan:
			}
		}
	}()
	return takeChan
}

func ExampleRepeatTake() {
	done := make(chan interface{})
	defer close(done)
	for num := range Take(done, Repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}
	fmt.Println()
}
