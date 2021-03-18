package timer

import "time"

// 定时器
type Timer interface {
	Open()

	Close()

	Destroy()
}

type TimeoutTimer struct {

	// 处理函数
	handle func()

	// 延迟事件  单位毫秒
	delay time.Duration

	timerid int64
}

// 创建 Timeout 定时器
func NewTimeoutTimer(cb func(), ms int) Timer {
	delay := time.Duration(ms)
	timeoutTimer := &TimeoutTimer{
		handle: cb,
		delay:  delay,
	}
	return timeoutTimer
}

func (this *TimeoutTimer) openHandle() {
	this.timerid = SetTimeout(this.handle, int(this.delay))
}

func (this *TimeoutTimer) Open() {
	//if this.timerid > 0 {
	//	this.Close()
	//}
	this.openHandle()
}

func (this *TimeoutTimer) Close() {
	ClearTimeout(this.timerid)
}

func (this *TimeoutTimer) Destroy() {
	this.Close()
	this.handle = nil
}

type IntervalTimer struct {

	// 处理函数
	// 如果 返回值为 false，定时器结束
	handle func() bool

	// 延迟事件  单位毫秒
	delay time.Duration

	timerid int64
}

// 创建定时器
//
// 如果 cb 返回 false，则 interval 结束
func NewIntervalTimer(cb func() bool, ms int) Timer {
	delay := time.Duration(ms)
	intervalTimer := &IntervalTimer{
		handle: cb,
		delay:  delay,
	}
	return intervalTimer
}

func (this *IntervalTimer) openHandle() {
	this.timerid = SetInterval(this.handle, int(this.delay))
}

func (this *IntervalTimer) Open() {
	//if this.timerid > 0 {
	//	this.Close()
	//}
	this.openHandle()
}

func (this *IntervalTimer) Close() {
	ClearInterval(this.timerid)
}

func (this *IntervalTimer) Destroy() {
	this.Close()
	this.handle = nil
}
