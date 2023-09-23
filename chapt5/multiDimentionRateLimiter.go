package chapt5

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func OpenMultiDimension() *APIConnectionMultiDimension {
	return &APIConnectionMultiDimension{
		apiLimit: MultiLimiter(
			rate.NewLimiter(Per(2, time.Second), 2),
			rate.NewLimiter(Per(10, time.Minute), 10),
		),
		diskLimit:    MultiLimiter(rate.NewLimiter(rate.Limit(1), 1)),
		networkLimit: MultiLimiter(rate.NewLimiter(rate.Limit(3), 3)),
	}
}

type APIConnectionMultiDimension struct {
	networkLimit,
	diskLimit,
	apiLimit *MultiLimiterStruct
}

func (a *APIConnectionMultiDimension) ReadFile(ctx context.Context) error {
	err := MultiLimiter(a.apiLimit, a.diskLimit).Wait(ctx)
	if err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

func (a *APIConnectionMultiDimension) ResolveAddress(ctx context.Context) error {
	err := MultiLimiter(a.apiLimit, a.networkLimit).Wait(ctx)
	if err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

func MultiDimensionRateLimiter() {
	defer log.Println("Done. ")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := OpenMultiDimension()
	var wg = sync.WaitGroup{}
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v\n", err)
			}
			log.Println("ReadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v\n", err)
			}
			log.Println("ResolveAddress")
		}()
	}

	wg.Wait()
}
