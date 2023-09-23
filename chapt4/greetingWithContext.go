package chapt4

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func GreetingWithContext() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := PrintGreetingWithContext(ctx); err != nil {
			fmt.Printf("cannot print greeting %v\n", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := PrintFarewellWithContext(ctx); err != nil {
			fmt.Printf("cannot print farewell %v\n", err)
			cancel()
		}
	}()

	wg.Wait()

}

func PrintGreetingWithContext(ctx context.Context) error {
	greeting, err := GenGreetingWithContext(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func PrintFarewellWithContext(ctx context.Context) error {
	farewell, err := GenFarewellWithContext(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func GenGreetingWithContext(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale, err := LocaleWithContext(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func GenFarewellWithContext(ctx context.Context) (string, error) {
	switch locale, err := LocaleWithContext(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func LocaleWithContext(ctx context.Context) (string, error) {

	// If deadline is less than 1 minute from now, then return immediately
	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}
