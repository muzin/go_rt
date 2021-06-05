package timer

import (
	"sync"
	"time"
)

var timerIdIncr int64 = 0
var mu sync.Mutex

func ApplyTimerId() int64 {
	mu.Lock()
	defer mu.Unlock()

	if timerIdIncr == 0 {
		timerIdIncr = time.Now().UnixNano()
	} else {
		timerIdIncr = timerIdIncr + 1
	}
	return timerIdIncr
}
