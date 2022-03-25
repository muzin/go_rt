package timer

import (
	"fmt"
	"github.com/muzin/go_rt/try"
	"time"
)

//var timeoutTaskStatusMap = make(map[int64]bool)
//var timeoutTaskStatusMap = *hash_map.NewHashMap() // map[int64]bool
//
//// SetTimeout
////	@param cb function
////	@param t ms
//func SetTimeout(cb func(), ms int) int64 {
//	id := ApplyTimerId()
//	timeoutTaskStatusMap.Put(id, true)
//
//	go func() {
//		defer try.CatchUncaughtException(func(throwable try.Throwable) {
//			fmt.Printf("SetTimeout Uncaught: %v", throwable)
//			throwable.PrintStackTrace()
//			timeoutTaskStatusMap.Remove(id)
//		})()
//
//		time.Sleep(time.Duration(ms) * time.Millisecond)
//
//		statusOk := timeoutTaskStatusMap.ContainsKey(id)
//		if statusOk {
//			obj := timeoutTaskStatusMap.Get(id)
//			var status bool = false
//			if obj != nil {
//				status = obj.(bool)
//			}
//			if status {
//				cb()
//			}
//		}
//
//		timeoutTaskStatusMap.Remove(id)
//	}()
//
//	return id
//}
//
//func ClearTimeout(id int64) {
//	statusOk := timeoutTaskStatusMap.ContainsKey(id)
//	if statusOk {
//		timeoutTaskStatusMap.Put(id, false)
//	}
//}

// SetTimeout
//	@param cb function
//	@param t ms
func SetTimeout(cb func(), ms int) *time.Ticker {
	ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
	go func() {
		defer try.CatchUncaughtException(func(throwable try.Throwable) {
			fmt.Printf("SetTimeout Uncaught: %v", throwable)
			throwable.PrintStackTrace()
			ticker.Stop()
		})()

	_:
		<-ticker.C
		cb()
		ticker.Stop()
	}()

	return ticker
}

func ClearTimeout(ticker *time.Ticker) {
	if ticker != nil {
		ticker.Stop()
	}
}
