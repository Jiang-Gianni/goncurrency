package chapt3

import (
	"fmt"
	"sync"
	"time"
)

func Cond() {
	c := sync.NewCond(&sync.Mutex{})
	queue := []int{}
	removeFromQueue := func(delay time.Duration, i int) {
		time.Sleep(delay)
		c.L.Lock()
		fmt.Println(queue)
		queue = queue[1:]
		fmt.Println("Removed from queue ", i)
		fmt.Println(queue)
		c.L.Unlock()
		c.Signal()
	}
	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Println("Adding to queue")
		queue = append(queue, i)
		go removeFromQueue(1*time.Second, i)
		c.L.Unlock()
	}
}
