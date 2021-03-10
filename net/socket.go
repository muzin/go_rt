package net

import (
	"net"
)

type Socket interface {
	Init()

	Connect(args ...interface{})

	Reconnect()

	ConnectHandle()

	Write(args ...interface{}) int

	OnConnect(listener func(...interface{}))

	OnData(listener func(...interface{}))

	OnError(listener func(...interface{}))

	OnClose(listener func(...interface{}))

	OnEnd(listener func(...interface{}))

	OnTimeout(listener func(...interface{}))

	On(t string, listener func(...interface{}))

	Once(t string, listener func(...interface{}))

	Emit(t string, args ...interface{}) bool

	EmitGo(t string, args ...interface{}) bool

	LocalAddr() net.Addr

	RemoteAddr() net.Addr

	Close()

	IsClose() bool

	SetOpeningStatus()

	SetOpenStatus()

	End()

	Pause()

	SetStreamTimeout(msecs int)

	GetStreamTimeout() int

	SetBufferSize(size int)

	GetBufferSize() int

	SetPort(port int)

	GetPort() int

	SetHost(host string)

	GetHost() string

	Destroy()
}
