package timer

import (
	"fmt"
	"github.com/muzin/go_rt/try"
	"time"
)

//var intervalTaskStatusMap = make(map[int64]bool)
//var intervalTaskStatusMap = *hash_map.NewHashMap() // map[int64]bool
//
////
//func SetInterval(cb func() bool, ms int) int64 {
//	id := ApplyTimerId()
//	intervalTaskStatusMap.Put(id, true)
//	go func() {
//		defer try.CatchUncaughtException(func(throwable try.Throwable) {
//			fmt.Printf("SetInterval Uncaught: %v\n", throwable)
//			throwable.PrintStackTrace()
//			intervalTaskStatusMap.Remove(id)
//		})()
//
//		for {
//			statusOk := intervalTaskStatusMap.ContainsKey(id)
//			obj := intervalTaskStatusMap.Get(id)
//			var status bool = false
//			if obj != nil {
//				status = obj.(bool)
//			}
//
//			if statusOk && status {
//				time.Sleep(time.Duration(ms) * time.Millisecond)
//				if statusOk && status {
//					cbret := cb()
//					if !cbret {
//						intervalTaskStatusMap.Remove(id)
//						break
//					}
//				} else {
//					intervalTaskStatusMap.Remove(id)
//					break
//				}
//			} else {
//				intervalTaskStatusMap.Remove(id)
//				break
//			}
//		}
//	}()
//	return id
//}
//
//func ClearInterval(id int64) {
//	statusOk := intervalTaskStatusMap.ContainsKey(id)
//	if statusOk {
//		intervalTaskStatusMap.Put(id, false)
//	}
//}

func SetInterval(cb func() bool, ms int) *time.Ticker {
	ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
	go func() {
		defer try.CatchUncaughtException(func(throwable try.Throwable) {
			fmt.Printf("SetInterval Uncaught: %v\n", throwable)
			throwable.PrintStackTrace()
			ticker.Stop()
		})()

		for {
		_:
			<-ticker.C
			cbret := cb()
			if !cbret {
				ticker.Stop()
				break
			}
		}
	}()
	return ticker
}

func ClearInterval(ticker *time.Ticker) {
	if ticker != nil {
		ticker.Stop()
	}
}
