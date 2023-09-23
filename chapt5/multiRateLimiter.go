package chapt5

import (
	"context"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiterStruct interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

type MultiLimiterStruct struct {
	limiters []RateLimiterStruct
}

func MultiLimiter(limiters ...RateLimiterStruct) *MultiLimiterStruct {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit)
	return &MultiLimiterStruct{limiters: limiters}
}

func (l *MultiLimiterStruct) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *MultiLimiterStruct) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

type APIConnectionMulti struct {
	rateLimiter *MultiLimiterStruct
}

func OpenMulti() *APIConnectionMulti {
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10)
	return &APIConnectionMulti{
		rateLimiter: MultiLimiter(secondLimit, minuteLimit),
	}
}

func (a *APIConnectionMulti) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// ReadFile simulation
	return nil
}

func (a *APIConnectionMulti) ResolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// ResolveAddress simulation
	return nil
}

func MultiRateLimiter() {
	defer log.Println("Done. ")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := OpenMulti()
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
