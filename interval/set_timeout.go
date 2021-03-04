package interval

import (
	"fmt"
	"github.com/muzin/go_rt/try"
	"time"
)

var timeoutTaskStatusMap = make(map[int]bool)

// SetTimeout
//	@param cb function
//	@param t ms
func SetTimeout(cb func(), ms int) int {
	id := time.Now().Nanosecond()
	timeoutTaskStatusMap[id] = true

	go func() {
		defer try.CatchUncaughtException(func(throwable try.Throwable) {
			fmt.Printf("SetTimeout Uncaught: %v", throwable)
		})

		time.Sleep(time.Duration(ms) * time.Millisecond)

		status, statusOk := timeoutTaskStatusMap[id]

		if statusOk {
			if status {
				cb()
			} else {
				delete(timeoutTaskStatusMap, id)
			}
		}
	}()

	return id
}

func ClearTimeout(id int) {
	_, statusOk := timeoutTaskStatusMap[id]
	if statusOk {
		timeoutTaskStatusMap[id] = false
	}
}
