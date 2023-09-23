package chapt4

import (
	"fmt"
	"sync"
	"time"
)

func GreetingWithoutContext() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := PrintGreetingWithoutContext(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := PrintFarewellWithoutContext(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Wait()

}

func PrintGreetingWithoutContext(done <-chan interface{}) error {
	greeting, err := GenGreetingWithoutContext(done)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func PrintFarewellWithoutContext(done <-chan interface{}) error {
	farewell, err := GenFarewellWithoutContext(done)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func GenGreetingWithoutContext(done <-chan interface{}) (string, error) {
	switch locale, err := LocaleWithoutContext(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func GenFarewellWithoutContext(done <-chan interface{}) (string, error) {
	switch locale, err := LocaleWithoutContext(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func LocaleWithoutContext(done <-chan interface{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("canceled")
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}
