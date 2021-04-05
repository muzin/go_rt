package timer

import (
	"fmt"
	"github.com/muzin/go_rt/collection/hash_map"
	"github.com/muzin/go_rt/try"
	"time"
)

//var timeoutTaskStatusMap = make(map[int64]bool)
var timeoutTaskStatusMap = *hash_map.NewHashMap() // map[int64]bool

// SetTimeout
//	@param cb function
//	@param t ms
func SetTimeout(cb func(), ms int) int64 {
	id := time.Now().UnixNano()

	timeoutTaskStatusMap.Put(id, true)

	go func() {
		defer try.CatchUncaughtException(func(throwable try.Throwable) {
			fmt.Printf("SetTimeout Uncaught: %v", throwable)
			timeoutTaskStatusMap.Remove(id)
		})()

		time.Sleep(time.Duration(ms) * time.Millisecond)

		statusOk := timeoutTaskStatusMap.ContainsKey(id)
		obj := timeoutTaskStatusMap.Get(id)
		var status bool = false
		if obj != nil {
			status = obj.(bool)
		}

		if statusOk && status {
			cb()
		}

		timeoutTaskStatusMap.Remove(id)
	}()

	return id
}

func ClearTimeout(id int64) {
	statusOk := timeoutTaskStatusMap.ContainsKey(id)
	if statusOk {
		timeoutTaskStatusMap.Put(id, false)
	}
}
