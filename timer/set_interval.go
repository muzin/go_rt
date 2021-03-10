package timer

import (
	"fmt"
	"github.com/muzin/go_rt/collection/hash_map"
	"github.com/muzin/go_rt/try"
	"time"
)

//var intervalTaskStatusMap = make(map[int64]bool)
var intervalTaskStatusMap = *hash_map.NewHashMap() // map[int64]bool

//
func SetInterval(cb func() bool, ms int) int64 {
	id := time.Now().UnixNano()
	intervalTaskStatusMap.Put(id, true)
	go func() {
		defer try.CatchUncaughtException(func(throwable try.Throwable) {
			fmt.Printf("SetInterval Uncaught: %v", throwable)
			intervalTaskStatusMap.Remove(id)
		})()

		for {
			statusOk := intervalTaskStatusMap.ContainsKey(id)
			status := intervalTaskStatusMap.Get(id).(bool)
			if statusOk && status {
				time.Sleep(time.Duration(ms) * time.Millisecond)
				if statusOk && status {
					cbret := cb()
					if !cbret {
						intervalTaskStatusMap.Remove(id)
						break
					}
				} else {
					intervalTaskStatusMap.Remove(id)
					break
				}
			} else {
				intervalTaskStatusMap.Remove(id)
				break
			}
		}
	}()
	return id
}

func ClearInterval(id int64) {
	statusOk := intervalTaskStatusMap.ContainsKey(id)
	if statusOk {
		intervalTaskStatusMap.Put(id, false)
	}
}
