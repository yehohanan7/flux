package utils

import (
	"time"
)

func WaitUntil(fn func() bool, timeout time.Duration) {
	ch := make(chan struct{})

	go func() {
		for {
			if fn() {
				ch <- struct{}{}
			}
		}
	}()

	select {
	case <-time.After(timeout):
		return
	case <-ch:
		return
	}
}
