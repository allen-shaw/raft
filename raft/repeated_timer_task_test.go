package raft

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestUnlock(t *testing.T) {
	var m sync.Mutex
	m.Lock()
	defer m.Unlock()

	fmt.Println("Hello World")

	m.Unlock()
}

func TestTimer(t *testing.T) {
	timer := time.NewTimer(200 * time.Millisecond)

	go func() {
		for {
			select {
			case c := <-timer.C:
				fmt.Println(c)
				timer.Reset(100 * time.Millisecond)
			}
		}
	}()

	time.Sleep(2 * time.Second)

	timer.Stop()
	time.Sleep(10 * time.Second)
}
