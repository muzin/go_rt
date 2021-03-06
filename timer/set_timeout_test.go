package timer

import (
	"sync"
	"testing"
)

func TestSetTimeout(t *testing.T) {

	t.Run("Test 设置定时器", func(t *testing.T) {

		var wg sync.WaitGroup

		wg.Add(1)

		timer := SetTimeout(func() {

			t.Logf("SetTimeout done")

			wg.Done()

		}, 1000)

		wg.Wait()

		ClearTimeout(timer)
	})
}

func TestTimeoutTimer(t *testing.T) {

	t.Run("Test 设置定时器 TimeoutTimer", func(t *testing.T) {

		var wg sync.WaitGroup

		wg.Add(1)

		timer := NewTimeoutTimer(func() {

			t.Logf("SetTimeout done")

			wg.Done()

		}, 1000)
		timer.Open()

		wg.Wait()

		timer.Close()
		timer.Destroy()

	})
}
