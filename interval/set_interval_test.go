package interval

import (
	"sync"
	"testing"
)

func TestSetInterval(t *testing.T) {
	t.Run("Test 设置定时器", func(t *testing.T) {

		var wg sync.WaitGroup

		loopCount := 3

		count := 0

		wg.Add(loopCount)

		timer := SetInterval(func() bool {

			t.Logf("SetTimeout done")

			wg.Done()

			count += 1

			if count == loopCount {
				return false
			} else {
				return true
			}

		}, 1000)

		wg.Wait()

		ClearInterval(timer)

	})
}
