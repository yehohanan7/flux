package utils

import "time"

func Every(duration time.Duration, fn func()) {
	tick := time.Tick(duration)

	for {
		select {
		case <-tick:
			fn()
		}
	}

}
