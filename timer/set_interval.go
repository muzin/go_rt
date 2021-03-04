package timer

import (
	"fmt"
	"github.com/muzin/go_rt/try"
	"time"
)

var intervalTaskStatusMap = make(map[int64]bool)

//
func SetInterval(cb func() bool, ms int) int64 {
	id := time.Now().UnixNano()
	intervalTaskStatusMap[id] = true
	go func() {
		defer try.CatchUncaughtException(func(throwable try.Throwable) {
			fmt.Printf("SetInterval Uncaught: %v", throwable)
			delete(intervalTaskStatusMap, id)
		})

		for {
			status, statusOk := intervalTaskStatusMap[id]
			if statusOk && status {
				time.Sleep(time.Duration(ms) * time.Millisecond)
				if statusOk && status {
					cbret := cb()
					if !cbret {
						delete(intervalTaskStatusMap, id)
						break
					}
				} else {
					delete(intervalTaskStatusMap, id)
					break
				}
			} else {
				delete(intervalTaskStatusMap, id)
				break
			}
		}
	}()
	return id
}

func ClearInterval(id int64) {
	_, statusOk := intervalTaskStatusMap[id]
	if statusOk {
		intervalTaskStatusMap[id] = false
	}
}
