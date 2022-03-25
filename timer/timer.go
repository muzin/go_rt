package timer

import (
	"time"
)

// 定时器
type Timer interface {
	Open()

	Close()

	Destroy()

	Id() int64

	ExpireTimeStamp() int64

	IsOpened() bool
}

type TimeoutTimer struct {

	// 处理函数
	handle func()

	// 延迟事件  单位毫秒
	delay time.Duration

	timerid int64

	ticker *time.Ticker

	isOpened bool
}

// 创建 Timeout 定时器
func NewTimeoutTimer(cb func(), ms int) Timer {
	delay := time.Duration(ms)
	timeoutTimer := &TimeoutTimer{
		handle:   cb,
		delay:    delay,
		isOpened: false,
	}
	return timeoutTimer
}

func (this *TimeoutTimer) openHandle() {
	this.timerid = ApplyTimerId()
	this.ticker = SetTimeout(this.handle, int(this.delay))
}

func (this *TimeoutTimer) Open() {
	this.isOpened = true
	this.openHandle()
}

func (this *TimeoutTimer) Close() {
	ClearTimeout(this.ticker)
}

func (this *TimeoutTimer) Destroy() {
	this.Close()
	this.handle = nil
	this.ticker = nil
}

func (this *TimeoutTimer) IsOpened() bool {
	return this.isOpened
}

func (this *TimeoutTimer) Id() int64 {
	return this.timerid
}

// 过期时间
// -1 永不过期
func (this *TimeoutTimer) ExpireTimeStamp() int64 {
	if this.IsOpened() {
		return this.timerid + int64(this.delay*time.Millisecond)
	} else {
		return -1
	}
}

type IntervalTimer struct {

	// 处理函数
	// 如果 返回值为 false，定时器结束
	handle func() bool

	// 延迟事件  单位毫秒
	delay time.Duration

	timerid int64

	ticker *time.Ticker

	isOpened bool
}

// 创建定时器
//
// 如果 cb 返回 false，则 interval 结束
func NewIntervalTimer(cb func() bool, ms int) Timer {
	delay := time.Duration(ms)
	intervalTimer := &IntervalTimer{
		handle:   cb,
		delay:    delay,
		isOpened: false,
	}
	return intervalTimer
}

func (this *IntervalTimer) openHandle() {
	this.timerid = ApplyTimerId()
	this.ticker = SetInterval(this.handle, int(this.delay))
}

func (this *IntervalTimer) Open() {
	//if this.timerid > 0 {
	//	this.Close()
	//}
	this.isOpened = true
	this.openHandle()
}

func (this *IntervalTimer) Close() {
	ClearInterval(this.ticker)
}

func (this *IntervalTimer) Destroy() {
	this.Close()
	this.handle = nil
	this.ticker = nil
}

func (this *IntervalTimer) IsOpened() bool {
	return this.isOpened
}

func (this *IntervalTimer) Id() int64 {
	return this.timerid
}

// 过期时间
// -1 永不过期
func (this *IntervalTimer) ExpireTimeStamp() int64 {
	return -1
}
