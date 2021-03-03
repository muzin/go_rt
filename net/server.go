package net

import "net"

type Server interface {
	Init()

	On(t string, listener func(...interface{}))

	Emit(t string, args ...interface{})

	Listen(args ...interface{})

	ConnectHandle()

	OnListen(listener func(...interface{}))

	OnConnect(listener func(...interface{}))

	OnClose(listener func(...interface{}))

	OnEnd(listener func(...interface{}))

	OnError(listener func(...interface{}))

	Close()

	End()

	SetPauseOnConnect(status bool)

	GetPauseOnConnect() bool

	SetNewSocketHandle(f func(net.Conn) Socket)

	GetNewSocketHandle() func(net.Conn) Socket

	SetPort(port int)

	GetPort() int

	SetHost(host string)

	GetHost() string
}
