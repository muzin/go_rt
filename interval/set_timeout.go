package interval

import (
	"fmt"
	"github.com/muzin/go_rt/try"
	"time"
)

var timeoutTaskStatusMap = make(map[int]bool)
var intervalTaskStatusMap = make(map[int]bool)

func SetTimeout(cb func(), t int) int {
	id := time.Now().Nanosecond()
	timeoutTaskStatusMap[id] = true

	go func() {
		defer try.CatchUncaughtException(func(throwable try.Throwable) {
			fmt.Printf("SetTimeout Uncaught: %v", throwable)
		})

		time.Sleep(time.Duration(t) * time.Millisecond)

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

func SetInterval(cb func(), t int) int {
	id := time.Now().Nanosecond()
	intervalTaskStatusMap[id] = true
	go func() {
		defer try.CatchUncaughtException(func(throwable try.Throwable) {
			fmt.Printf("SetInterval Uncaught: %v", throwable)
		})

		for {
			status, statusOk := intervalTaskStatusMap[id]
			if statusOk && status {
				time.Sleep(time.Duration(t) * time.Millisecond)
				if statusOk {
					if status {
						cb()
					} else {
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
