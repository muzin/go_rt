package events

import (
	"fmt"
	str "github.com/muzin/go_rt/lang/str"
	"github.com/muzin/go_rt/print/colorstr"
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

var (
	// 空 监听器集合
	EMPTY_LISTENERS = make([]func(...interface{}), 0)
)

// 事件发射器
type EventEmitter struct {

	// 事件监听数量
	eventsCount int

	//events map[string][]func(...interface{})

	// 事件集合中存放的事件包裹器
	events *sync.Map

	// 最大监听事件数，超出后警告
	maxListeners int

	// 是否允许事件有多个监听函数
	prepend bool

	// 错误处理函数
	errorEventWrap *EventWrap

	// 事件处理通道
	eventChannel chan EventChanWrap

	// 事件处理通道是否结束
	eventChannelFinished bool

	mu sync.RWMutex
}

func NewEventEmitter() *EventEmitter {
	emitter := &EventEmitter{}
	emitter.init()
	return emitter
}

// init
func (this *EventEmitter) init() {
	// 初始化 events
	if nil == this.events {
		//this.events = make(map[string][]func(...interface{}))
		this.events = &sync.Map{}
		this.eventsCount = 0
	}

	// 设置 最大监听数量
	this.maxListeners = DEFAULT_MAX_LISTENERS

	// 默认 新增的监听器向后方
	this.prepend = false

	this.eventChannel = make(chan EventChanWrap, 100)

	this.eventChannelFinished = false

	go this.eventHandler()

}

func (this *EventEmitter) Open() bool {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.eventChannelFinished {
		this.eventChannelFinished = false
		this.eventChannel = make(chan EventChanWrap, 100)
		go this.eventHandler()
		return true
	} else {
		return false
	}
}

func (this *EventEmitter) Close() bool {
	this.mu.Lock()
	defer this.mu.Unlock()

	if !this.eventChannelFinished {
		this.eventChannel <- EventChanWrap{
			t: CloseEventChanType,
		}
		return true
	} else {
		return false
	}
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

	this.mu.RLock()
	defer this.mu.RUnlock()

	// 如果 事件通道已结束 不允许发射任何事件
	if this.eventChannelFinished {
		return false
	}

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

	handlersEventWrapPtr, _ := events.Load(t)

	var handlers []func(...interface{})
	if nil != handlersEventWrapPtr {
		eventWrapPtr := handlersEventWrapPtr.(*EventWrap)
		handlers = eventWrapPtr.GetListeners()

		// 提前将 发射标识设置为已发射，确保只执行一次
		eventWrapPtr.Emitted()

		if len(handlers) > 0 {
			eventWrapPtr.emitcount++
		}

	}

	if nil == handlers {
		return false
	}

	for i := 0; i < len(handlers); i++ {
		handler := handlers[i]
		if handler != nil {
			this.eventChannel <- EventChanWrap{
				t:       NormalEventChanType,
				handler: handler,
				args:    args,
			}
		}
	}

	return true
}

func (this *EventEmitter) eventHandler() {
	for {
		eventChanWrap, isOpen := <-this.eventChannel
		if isOpen {
			chanWrapType := eventChanWrap.t
			handler := eventChanWrap.handler
			args := eventChanWrap.args

			if NormalEventChanType == chanWrapType {
				try.Try(func() {
					handler(args...)
				}, try.CatchUncaughtException(func(throwable try.Throwable) {
					if this.errorEventWrap != nil {
						errListeners := this.errorEventWrap.GetListeners()
						for _, errListener := range errListeners {
							errListener(throwable)
						}
					}
				}))
			} else if CloseEventChanType == chanWrapType {
				close(this.eventChannel)
				this.eventChannelFinished = true
			}

		} else {
			break
		}
	}
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
	return this.addListener(t, listener, this.GetPrepend(), false)
}

func (this *EventEmitter) AddAppendListener(t string, listener func(...interface{})) *EventEmitter {
	return this.addListener(t, listener, false, false)
}

func (this *EventEmitter) addListener(t string, listener func(...interface{}), prepend bool, isOnce bool) *EventEmitter {

	var existing []func(...interface{})

	events := this.events

	if nil == events {
		//this.events = make(map[string][]func(...interface{}))
		this.events = &sync.Map{}
		events = this.events
		this.eventsCount = 0
	}

	_, newOk := events.Load(EVENT_EMITTER_NEW_LISTENER_NAME)
	if newOk {
		this.EmitGo(EVENT_EMITTER_NEW_LISTENER_NAME, t, listener)
		events = this.events
	}
	existingInterfaceEventWrapPtr, _ := events.Load(t)
	if nil != existingInterfaceEventWrapPtr {
		eventWrap := existingInterfaceEventWrapPtr.(*EventWrap)
		existing = eventWrap.GetListeners()
	}

	var eventWrap *EventWrap = nil

	// 如果 没有获取到事件包裹器 添加
	if nil == existingInterfaceEventWrapPtr {
		eventWrap = &EventWrap{
			name:    t,
			isOnce:  isOnce,
			emitted: false,
		}
		eventWrap.AddHandler(listener, prepend)

		this.events.Store(t, eventWrap)
		this.eventsCount += 1

	} else {
		eventWrap = existingInterfaceEventWrapPtr.(*EventWrap)
		eventWrap.AddHandler(listener, prepend)
	}

	m := this.GetMaxListeners()
	if m > 0 && len(existing) > m {
		fmt.Println(colorstr.Yellow(fmt.Sprintf("Possible EventEmitter memory leak detected. "+
			"%v %v listeners "+
			"added. Use emitter.SetMaxListeneres() to increase limit.",
			len(existing), t)))
	}

	if t == "error" {
		this.errorEventWrap = eventWrap
	}

	events = nil

	return this
}

func (this *EventEmitter) ListenerCount(t string) int {

	events := this.events

	if nil != events {

		evlistenerInterfaceEventWrapPtr, ok := events.Load(t)

		if ok {
			eventWrap := evlistenerInterfaceEventWrapPtr.(*EventWrap)
			evlistener := eventWrap.GetListeners()
			return len(evlistener)
		} else {
			return 0
		}
	}
	return 0
}

func (this *EventEmitter) EventNames() []string {
	strings := make([]string, 0)
	if this.eventsCount > 0 && nil != this.events {
		this.events.Range(func(key interface{}, value interface{}) bool {
			strings = append(strings, key.(string))
			return true
		})
	}
	return strings
}

// 只监听一次
func (this *EventEmitter) Once(t string, listener func(...interface{})) *EventEmitter {
	return this.addListener(t, listener, false, true)
}

// 废弃
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

	handlersEventWrapPtr, ok := events.Load(t)

	if ok {
		if (this.eventsCount - 1) <= 0 {
			//this.events = make(map[string][]func(...interface{}))
			this.events = &sync.Map{}
			this.eventsCount = 0
		} else {
			//delete(this.events, t)
			this.events.Delete(t)
			this.eventsCount -= 1
			_, rmOk := events.Load(EVENT_EMITTER_REMOVE_LISTENER_NAME)
			if rmOk {
				this.EmitGo(EVENT_EMITTER_REMOVE_LISTENER_NAME, t, handlersEventWrapPtr)
			}
		}
	}

	return this
}

// 移除事件
func (this *EventEmitter) RemoveAllListener() *EventEmitter {

	events := this.events

	this.events = &sync.Map{}
	this.eventsCount = 0

	if nil != events {
		eventNames := this.EventNames()
		for i := 0; i < len(eventNames); i++ {
			eventName := eventNames[i]
			events.Delete(eventName)
		}
	}

	return this
}

// 移除事件
func (this *EventEmitter) ResetOnce() *EventEmitter {

	events := this.events
	if nil != events {
		eventNames := this.EventNames()
		for i := 0; i < len(eventNames); i++ {
			eventName := eventNames[i]
			evlistenerInterfaceEventWrapPtr, ok := events.Load(eventName)

			if ok {
				eventWrap := evlistenerInterfaceEventWrapPtr.(*EventWrap)
				eventWrap.emitted = false
				eventWrap.emitcount = 0
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

// 销毁
func (this *EventEmitter) Destory() {

}

// 事件包裹器
type EventWrap struct {
	// 事件名称
	name string

	// 监听事件
	listeners []func(...interface{})

	// 是否只能执行一次 默认：false
	isOnce bool

	// 是否触发过 默认：false
	emitted bool

	emitcount int

	mu sync.Mutex
}

// 获取名称
func (this *EventWrap) GetName() string {
	return this.name
}

// 获取事件
func (this *EventWrap) GetListeners() []func(...interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	// 如果 是只触发一次
	if this.isOnce && this.emitted {
		return EMPTY_LISTENERS
	}

	return this.listeners
}

// 获取事件
func (this *EventWrap) AddHandler(handler func(...interface{}), prepend bool) {
	if prepend {
		this.listeners = append(append(make([]func(...interface{}), 0), handler), this.listeners...)
	} else {
		this.listeners = append(this.listeners, handler)
	}
}

func (this *EventWrap) Emitted() {
	if !this.emitted {
		this.emitted = true
	}
}

func (this *EventWrap) IsOnce() bool {
	return this.isOnce
}

type EventChanWrapType int

var (
	NormalEventChanType EventChanWrapType = 0 // 事件通道Wrap
	CloseEventChanType  EventChanWrapType = 1 // 关闭事件通道Wrap
)

type EventChanWrap struct {
	t EventChanWrapType

	handler func(...interface{})

	args []interface{}
}
