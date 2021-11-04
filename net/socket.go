package net

import (
	"net"
)

type Socket interface {

	// 初始化
	Init()

	// 连接
	Connect(args ...interface{})

	// 重连
	Reconnect()

	// 连接 处理器
	ConnectHandle()

	// 写数据
	Write(args ...interface{}) int

	// 监听 连接 事件
	OnConnect(listener func(...interface{}))

	// 监听 接收数据 事件
	OnData(listener func(...interface{}))

	// 监听 错误 事件
	OnError(listener func(...interface{}))

	// 监听 关闭 事件
	OnClose(listener func(...interface{}))

	// 监听 结束 事件
	OnEnd(listener func(...interface{}))

	// 监听 超时 事件
	OnTimeout(listener func(...interface{}))

	// 监听 事件
	On(t string, listener func(...interface{}))

	// 仅监听一次 事件
	Once(t string, listener func(...interface{}))

	// 发射 事件
	Emit(t string, args ...interface{}) bool

	// go 发射 事件
	EmitGo(t string, args ...interface{}) bool

	// 获取本地地址
	LocalAddr() net.Addr

	// 获取远程地址
	RemoteAddr() net.Addr

	// 关闭
	Close()

	// 是否 关闭
	IsClose() bool

	// 设置 打开中 状态
	SetOpeningStatus()

	// 设置 打开 状态
	SetOpenStatus()

	// 结束
	End()

	// 暂停
	Pause()

	// 是否连接
	IsConnected() bool

	// 设置 流 超时时间
	SetStreamTimeout(msecs int)

	// 获取 流 超时时间
	GetStreamTimeout() int

	// 设置 缓冲区大小
	SetBufferSize(size int)

	// 获取  缓冲区大小
	GetBufferSize() int

	// 设置 端口
	SetPort(port int)

	// 获取 端口
	GetPort() int

	// 设置 主机
	SetHost(host string)

	// 获取 主机
	GetHost() string

	// 销毁
	Destroy()
}
