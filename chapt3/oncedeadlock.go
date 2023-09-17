package chapt3

import (
	"fmt"
	"sync"
)

func OnceDeadlock() {
	var onceA, onceB sync.Once
	var initB func()
	initA := func() {
		fmt.Println("initA")
		onceB.Do(initB)
	}
	initB = func() {
		fmt.Println("initB")
		onceA.Do(initA)
	}
	onceA.Do(initA)
}
