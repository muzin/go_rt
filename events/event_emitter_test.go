package events

import (
	"fmt"
	"github.com/muzin/go_rt/try"
	"testing"
	"time"
)

func TestEventEmitter_New(t *testing.T) {
	t.Run("测试 EventTmitter New", func(t *testing.T) {

		eventemitter := NewEventEmitter()

		if nil == eventemitter {
			t.Error("NewEventEmitter emitter is not null are expect")
		}

	})
}

func TestEventEmitter_On(t *testing.T) {
	t.Run("测试 EventTmitter On", func(t *testing.T) {

		try.CatchUncaughtException(func(err try.Throwable) {
			fmt.Errorf("UnhandledError: %v", err.GetStackTrace())
		})

		eventemitter := NewEventEmitter()

		if nil == eventemitter {
			t.Error("NewEventEmitter emitter is not null are expect")
		}

		eventemitter.On("error", func(args ...interface{}) {
			err := args[0]
			fmt.Printf("error %v\n", err)
		})

		except := 10000000
		count := 0
		eventemitter.On("data", func(args ...interface{}) {
			//t.Logf("args: %v", args)
			count++

			//if count > 1 {
			//	try.Throw(try.UnhandledError.NewThrow("index: " + strconv.Itoa(count) + " exist unhandled error"))
			//}

		})

		for i := 0; i < except; i++ {
			eventemitter.Emit("data", i)
		}

		time.Sleep(50 * time.Millisecond)

		//time.Sleep(10 * time.Second)
		if count != except {
			t.Errorf("NewEventEmitter expect: %v but: %v", except, count)
		} else {
			t.Logf("NewEventEmitter expect: %v result: %v", except, count)
		}

		listenerCount := eventemitter.ListenerCount("data")

		if listenerCount != 1 {
			t.Errorf("listenerCount expect: %v but: %v", 1, listenerCount)
		} else {
			t.Logf("listenerCount expect: %v result: %v", 1, listenerCount)
		}

	})
}

func TestEventEmitter_Once(t *testing.T) {
	t.Run("测试 EventTmitter Once", func(t *testing.T) {

		try.CatchUncaughtException(func(err try.Throwable) {
			fmt.Errorf("UnhandledError: %v", err.GetStackTrace())
		})

		eventemitter := NewEventEmitter()

		if nil == eventemitter {
			t.Error("NewEventEmitter emitter is not null are expect")
		}

		except := 100
		count := 0
		eventemitter.On("data", func(args ...interface{}) {
			t.Logf("args: %v", args)
			count++
		})

		for i := 0; i < 100; i++ {
			func(v int) {
				eventemitter.Emit("data", v)
			}(i)
		}

		eventemitter.Stop()

		eventemitter.Emit("data", 501)

		time.Sleep(1 * time.Second)

		if count != except {
			t.Errorf("NewEventEmitter expect: %v but: %v", except, count)
		} else {
			t.Logf("NewEventEmitter expect: %v result: %v", except, count)
		}

	})
}
