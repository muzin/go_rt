package interval

import (
	"fmt"
	"github.com/muzin/go_rt/try"
	"time"
)

var intervalTaskStatusMap = make(map[int]bool)

//
func SetInterval(cb func() bool, ms int) int {
	id := time.Now().Nanosecond()
	intervalTaskStatusMap[id] = true
	go func() {
		defer try.CatchUncaughtException(func(throwable try.Throwable) {
			fmt.Printf("SetInterval Uncaught: %v", throwable)
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

func ClearInterval(id int) {
	_, statusOk := intervalTaskStatusMap[id]
	if statusOk {
		intervalTaskStatusMap[id] = false
	}
}
