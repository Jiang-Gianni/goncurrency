package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Jiang-Gianni/goncurrency/chapt4"
)

func main() {
	defer func(start time.Time) {
		fmt.Println("since start (ms): ", time.Since(start))
	}(time.Now())
	chapt4.ExampleRepeatFn()

	fmt.Println(runtime.NumCPU())
}
