package events

import (
	"fmt"
	colorstr "github.com/muzin/go_rt/collection/print/color_string"
	str "github.com/muzin/go_rt/lang/str"
	"github.com/muzin/go_rt/try"
	"sync"
)

const (
	// 默认最大事件监听数量
	DEFAULT_MAX_LISTENERS = 10

	EVENT_EMITTER_ERROR_LISTENER_NAME = "error"

	EVENT_EMITTER_NEW_LISTENER_NAME = "newListener"

	EVENT_EMITTER_REMOVE_LISTENER_NAME = "removeListener"
)

// 事件发射器
type EventEmitter struct {
	eventsCount int

	//events map[string][]func(...interface{})

	events *sync.Map

	maxListeners int

	// 是否允许事件有多个监听函数
	prepend bool
}

func NewEventEmitter() *EventEmitter {
	emitter := &EventEmitter{}
	emitter.init()
	return emitter
}

// init
func (this *EventEmitter) init() {
	if nil == this.events {
		//this.events = make(map[string][]func(...interface{}))
		this.events = &sync.Map{}
		this.eventsCount = 0
	}
	this.maxListeners = DEFAULT_MAX_LISTENERS
	this.prepend = true
}

// SetMaxListeners
func (this *EventEmitter) SetMaxListeners(n uint) *EventEmitter {
	this.maxListeners = int(n)
	return this
}

// GetMaxListeners
func (this *EventEmitter) GetMaxListeners() int {
	return this.maxListeners
}

// emit event
func (this *EventEmitter) Emit(t string, args ...interface{}) bool {
	return this.emit(t, false, args...)
}

func (this *EventEmitter) EmitGo(t string, args ...interface{}) bool {
	go this.emit(t, true, args...)
	return true
}

func (this *EventEmitter) emit(t string, rungo bool, args ...interface{}) bool {

	doError := false
	if t == "error" {
		doError = true
	}

	events := this.events

	if nil != events {
		errorHandles, _ := events.Load("error")
		if doError == true && nil == errorHandles {
			doError = true
		} else {
			doError = false
		}
	} else if !doError {
		return false
	}

	if doError {
		var er try.Throwable
		if len(args) > 0 {
			switch args[0].(type) {
			case try.Throwable:
				throwable := (args[0]).(try.Throwable)
				er = throwable
				break
			default:
				break
			}
		}
		// 如果 在这里 抛出异常，说明没有监听 error 事件，监听后 不会 panic
		if nil != er {
			try.Throw(try.UnhandledError.NewThrow(er.Error()))
		} else {
			err := str.Strval(args[0])
			try.Throw(try.UnhandledError.NewThrow(err))
		}
	}

	handlersInterface, _ := events.Load(t)
	var handlers []func(...interface{})
	if nil != handlersInterface {
		handlers = handlersInterface.([]func(...interface{}))
	}

	if nil == handlers {
		return false
	}

	var errorHandler func(...interface{})

	// 获取 异常处理函数
	errorHandlersInterface, errOk := events.Load("error")
	var errorHandlers []func(...interface{})
	if nil != errorHandlersInterface {
		errorHandlers = errorHandlersInterface.([]func(...interface{}))
	}

	if errOk && len(errorHandlers) > 0 {
		errorHandler = errorHandlers[0]
	}

	// 是否 并行 调用 回调函数
	if rungo == true {
		for i := 0; i < len(handlers); i++ {
			handler := handlers[i]
			go func() {
				if handler != nil {
					defer try.CatchUncaughtException(errorOfCallerHandle(errorHandler))()
					handler(args...)
				}
			}()
		}
	} else {
		for i := 0; i < len(handlers); i++ {
			handler := handlers[i]
			func() {
				if handler != nil {
					defer try.CatchUncaughtException(errorOfCallerHandle(errorHandler))()
					handler(args...)
				}
			}()
		}
	}

	return true
}

// 调用函数时出现错误的处理函数
func errorOfCallerHandle(errorHandler func(...interface{})) func(try.Throwable) {
	return func(throwable try.Throwable) {
		if nil != errorHandler {
			errorHandler(throwable)
		} else {
			fmt.Println(colorstr.Red("The EventEmitter is missing listening for the 'error' event. " +
				"Use emitter.OnError(func) or emitter.On('error', func)."))
			try.Throw(throwable)
		}
	}
}

// listen event
func (this *EventEmitter) On(t string, listener func(...interface{})) *EventEmitter {
	return this.AddListener(t, listener)
}

func (this *EventEmitter) OnError(listener func(...interface{})) *EventEmitter {
	return this.On(EVENT_EMITTER_ERROR_LISTENER_NAME, listener)
}

func (this *EventEmitter) OnNewListener(listener func(...interface{})) *EventEmitter {
	return this.On(EVENT_EMITTER_NEW_LISTENER_NAME, listener)
}

func (this *EventEmitter) OnRemoveListener(listener func(...interface{})) *EventEmitter {
	return this.On(EVENT_EMITTER_REMOVE_LISTENER_NAME, listener)
}

func (this *EventEmitter) AddListener(t string, listener func(...interface{})) *EventEmitter {
	return this.addListener(t, listener, this.GetPrepend())
}

func (this *EventEmitter) addListener(t string, listener func(...interface{}), prepend bool) *EventEmitter {

	var existing []func(...interface{})

	events := this.events

	if nil == events {
		//this.events = make(map[string][]func(...interface{}))
		this.events = &sync.Map{}
		events = this.events
		this.eventsCount = 0
	} else {
		_, newOk := events.Load(EVENT_EMITTER_NEW_LISTENER_NAME)
		if newOk {
			this.EmitGo(EVENT_EMITTER_NEW_LISTENER_NAME, t, listener)
			events = this.events
		}
		existingInterface, _ := events.Load(t)
		if nil != existingInterface {
			existing = existingInterface.([]func(...interface{}))
		}
	}

	if nil == existing {
		handlers := make([]func(...interface{}), 1)
		handlers[0] = listener
		this.events.Store(t, handlers)
		this.eventsCount += 1
	} else {
		if prepend {
			existing = append(append(make([]func(...interface{}), 0), listener), existing...)
		} else {
			existing = append(make([]func(...interface{}), 0), listener)
		}
		this.events.Store(t, existing)
	}

	m := this.GetMaxListeners()
	if m > 0 && len(existing) > m {
		fmt.Println(colorstr.Yellow(fmt.Sprintf("Possible EventEmitter memory leak detected. "+
			"%v %v listeners "+
			"added. Use emitter.SetMaxListeneres() to increase limit.",
			len(existing), t)))
	}

	events = nil

	return this
}

func (this *EventEmitter) ListenerCount(t string) int {

	events := this.events

	if nil != events {

		evlistenerInterface, ok := events.Load(t)

		if ok {
			evlistener := evlistenerInterface.([]func(...interface{}))
			return len(evlistener)
		} else {
			return 0
		}
	}
	return 0
}

func (this *EventEmitter) EventNames() []string {
	strings := make([]string, 0)
	if this.eventsCount > 0 {
		this.events.Range(func(key interface{}, value interface{}) bool {
			strings = append(strings, key.(string))
			return true
		})
	}
	return strings
}

// 监听一次
func (this *EventEmitter) Once(t string, listener func(...interface{})) *EventEmitter {
	return this.On(t, this.onceWrap(t, listener))
}

func (this *EventEmitter) onceWrap(t string, listener func(...interface{})) func(...interface{}) {
	return func(args ...interface{}) {
		this.RemoveListener(t)
		listener(args...)
	}
}

// 移除事件
func (this *EventEmitter) RemoveListener(t string) *EventEmitter {

	events := this.events

	if nil == events {
		return this
	}

	handlers, ok := events.Load(t)

	if ok {
		if (this.eventsCount - 1) == 0 {
			//this.events = make(map[string][]func(...interface{}))
			this.events = &sync.Map{}
			this.eventsCount = 0
		} else {
			//delete(this.events, t)
			this.events.Delete(t)
			this.eventsCount--
			_, rmOk := events.Load(EVENT_EMITTER_REMOVE_LISTENER_NAME)
			if rmOk {
				this.EmitGo(EVENT_EMITTER_REMOVE_LISTENER_NAME, t, handlers)
			}
		}
	}

	return this
}

func (this *EventEmitter) SetPrepend(prepend bool) {
	this.prepend = prepend
}

func (this *EventEmitter) GetPrepend() bool {
	return this.prepend
}

//
//func (this *EventEmitter) EnableThreadSafe() {
//	this.isThreadSafe = true
//}
//
//func (this *EventEmitter) DisableThreadSafe() {
//	this.isThreadSafe = false
//}
//
//func (this *EventEmitter) IsThreadSafe() bool {
//	return this.isThreadSafe
//}

// 销毁
func (this *EventEmitter) Destory() {
	go func() {

		this.eventsCount = 0
		this.maxListeners = 0
		this.prepend = false
		if nil != this.events {
			eventNames := this.EventNames()
			for i := 0; i < len(eventNames); i++ {
				eventName := eventNames[i]
				this.events.Delete(eventName)
			}
			this.events = nil
		}
	}()
}
