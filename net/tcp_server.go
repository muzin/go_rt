package net

import (
	"github.com/muzin/go_rt/events"
	"github.com/muzin/go_rt/lang/str"
	"github.com/muzin/go_rt/try"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	DEFAULT_IPV4_ADDR = "0.0.0.0"
	DEFAULT_IPV6_ADDR = "::"
)

// 服务 Listen 异常
var ServerListenException = try.DeclareException("ServerListenException")

// 服务 Accept 异常
var ServerAcceptException = try.DeclareException("ServerAcceptException")

// 服务 Close 异常
var ServerCloseException = try.DeclareException("ServerCloseException")

type TCPServer struct {
	events.EventEmitter

	net.Listener

	// 主机
	host string

	// 端口
	port int

	// 连接数
	connections int

	// 使用 协程
	// @Deprecated
	usingWorkers bool

	// @Deprecated
	allowHalfOpen bool

	// 暂停 监听 连接
	pauseOnConnect bool

	// 是否结束
	ended bool

	// socket 结束 WaitGroup
	socketsEndWg sync.WaitGroup

	// 是否初始化
	inited bool

	// 创建 socket of server 处理函数
	newSocketHandle func(conn net.Conn) Socket
}

/*
 * 创建server
 */
func NewTCPServer() *TCPServer {
	s := &TCPServer{
		EventEmitter: *events.NewEventEmitter(),
	}
	s.Init()
	return s
}

// 初始化 函数
func (this *TCPServer) Init() {

	if !this.inited {
		this.inited = true
	} else {
		return
	}

	// 设置 创建 socket 处理函数
	this.SetNewSocketHandle(func(conn net.Conn) Socket {
		return newSocketForServer(conn)
	})

	// 默认 监听 一个 空 error 事件
	this.OnError(func(...interface{}) {})
}

func (this *TCPServer) On(t string, listener func(...interface{})) {
	this.EventEmitter.On(t, listener)
}

func (this *TCPServer) Emit(t string, args ...interface{}) {
	this.EventEmitter.Emit(t, args)
}

// 监听端口
// Listen(port, [host]])
func (this *TCPServer) Listen(args ...interface{}) {

	if len(args) >= 1 {
		this.port = args[0].(int)
	}
	if len(args) >= 2 {
		this.host = args[1].(string)
	}
	//if len(args) >= 3 {
	//	this.usingWorkers = args[2].(bool)
	//}
	//if len(args) >= 4 {
	//	this.allowHalfOpen = args[3].(bool)
	//}

	var network = "tcp"
	var address = this.host + ":" + strconv.Itoa(this.port)

	server, listenErr := net.Listen(network, address)
	if listenErr != nil {
		try.Throw(ServerListenException.NewThrow(listenErr.Error()))
	}

	// 加入 服务结束等待组中
	GetSocketWaitGroup("tcp_server Listen() WaitGroup add 1").Add(1)

	this.Listener = server
	// 处理 连接
	go this.ConnectHandle()

	go func() {
		// 发送 监听事件
		this.EmitGo("listen", network, address)
	}()

}

// 处理所有连接函数
func (this *TCPServer) ConnectHandle() {
	server := this.Listener
	for {
		if this.pauseOnConnect == false {
			func() {
				// 如果有异常， 发送 error 事件
				defer try.Catch(ServerAcceptException, func(throwable try.Throwable) {
					this.EmitGo("error", throwable)
				})()

				// 接收来自 client 的连接,会阻塞
				conn, err := server.Accept()

				if err != nil {
					// close network connect
					if str.EndsWith(err.Error(), "use of closed network connection") {
						this.Emit("close", true) // 主动关闭
					} else {
						try.Throw(ServerAcceptException.NewThrow(err.Error()))
					}
					return
				}

				// 发送 连接 事件
				newSocketHandle := this.GetNewSocketHandle()
				newSocket := newSocketHandle(conn)

				// server 当前连接数 +1
				this.connections += 1

				// 连接结束后，socketsEndWg.Done()
				GetSocketWaitGroup("tcp_server ConnectHandle() socket connect WaitGroup add 1").Add(1)
				newSocket.On("end", func(args ...interface{}) {
					// 当 socket 关闭时， 当前连接数 -1
					this.connections -= 1
					GetSocketWaitGroup("tcp_server ConnectHandle() socket connect WaitGroup done 1").Done()
				})

				// 发送 有新 socket 连接事件
				this.EmitGo("connect", newSocket)

			}()
		} else {
			// 如果 还没有结束， 休息 50ms，继续监听
			if this.ended == false {
				time.Sleep(50 * time.Millisecond)
			} else {
				// 如果结束，跳出循环
				break
			}
		}
	}
}

// OnListen
//
// @param listener func(network string, address string)
//
func (this *TCPServer) OnListen(listener func(...interface{})) {
	this.EventEmitter.Once("listen", listener)
}

// OnConnect
//
// @param listener func(socket net.Socket)
//
func (this *TCPServer) OnConnect(listener func(...interface{})) {
	this.On("connect", listener)
	this.AddAppendListener("connect", func(args ...interface{}) {
		socket := args[0].(Socket)
		go socket.ConnectHandle()
	})
}

// 当关闭时 触发
func (this *TCPServer) OnClose(listener func(...interface{})) {
	this.Once("close", listener)
}

// 当结束时 触发
func (this *TCPServer) OnEnd(listener func(...interface{})) {
	this.Once("end", listener)
}

// 监听 error 事件
// // @param listener func(throwable try.throwable)
func (this *TCPServer) OnError(listener func(...interface{})) {
	this.On("error", listener)
}

// 关闭Server，等待socket全部关闭
func (this *TCPServer) Close() {

	// 暂停接受连接
	this.pauseOnConnect = true

	err := this.Listener.Close()
	if err != nil {
		this.Emit("error", ServerCloseException.NewThrow(err.Error()))
	}

	// 关闭后，发送关闭事件
	//go this.Emit("close")

	// go 等待 socket都结束后，发送 结束事件
	go func() {
		//this.socketsEndWg.Wait()

		// 已结束
		this.ended = true

		go this.Emit("end")
	}()

	// 记录 server 结束
	GetSocketWaitGroup("tcp_server Close() WaitGroup done 1").Done()

}

// 结束Server，不等待socket全部关闭，强行关闭
func (this *TCPServer) End() {
	if !this.ended {
		this.Close()
	}
}

func (this *TCPServer) SetPauseOnConnect(status bool) {
	this.pauseOnConnect = status
}

func (this *TCPServer) GetPauseOnConnect() bool {
	return this.pauseOnConnect
}

func (this *TCPServer) SetNewSocketHandle(f func(conn net.Conn) Socket) {
	this.newSocketHandle = f
}

func (this *TCPServer) GetNewSocketHandle() func(conn net.Conn) Socket {
	return this.newSocketHandle
}

func (this *TCPServer) SetPort(port int) {
	this.port = port
}

func (this *TCPServer) GetPort() int {
	return this.port
}

func (this *TCPServer) SetHost(host string) {
	this.host = host
}

func (this *TCPServer) GetHost() string {
	return this.host
}

// 销毁
func (this *TCPServer) Destroy() {

}
