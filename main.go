package main

import (
	"fmt"
	"time"

	"github.com/Jiang-Gianni/goncurrency/chapt5"
)

func main() {
	defer func(start time.Time) {
		fmt.Println("since start: ", time.Since(start))
	}(time.Now())
	chapt5.ExampleHealing()
}
