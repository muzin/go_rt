package timer

import (
	"fmt"
	"github.com/muzin/go_rt/try"
	"time"
)

var timeoutTaskStatusMap = make(map[int64]bool)

// SetTimeout
//	@param cb function
//	@param t ms
func SetTimeout(cb func(), ms int) int64 {
	id := time.Now().UnixNano()
	timeoutTaskStatusMap[id] = true

	go func() {
		defer try.CatchUncaughtException(func(throwable try.Throwable) {
			fmt.Printf("SetTimeout Uncaught: %v", throwable)
			delete(timeoutTaskStatusMap, id)
		})

		time.Sleep(time.Duration(ms) * time.Millisecond)

		status, statusOk := timeoutTaskStatusMap[id]

		if statusOk && status {
			cb()
		}

		delete(timeoutTaskStatusMap, id)
	}()

	return id
}

func ClearTimeout(id int64) {
	_, statusOk := timeoutTaskStatusMap[id]
	if statusOk {
		timeoutTaskStatusMap[id] = false
	}
}
