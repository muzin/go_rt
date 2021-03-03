package net

import (
	"github.com/muzin/go_rt/events"
	"github.com/muzin/go_rt/lang/str"
	"github.com/muzin/go_rt/system"
	"github.com/muzin/go_rt/try"
	"net"
	"strconv"
	"time"
)

var SocketConnectException = try.DeclareException("SocketConnectException")
var SocketReadException = try.DeclareException("SocketReadException")
var SocketCloseException = try.DeclareException("SocketCloseException")

type TCPSocket struct {
	events.EventEmitter

	net.Conn

	port int

	host string

	// 是否连接中
	connecting bool

	// 等待中
	pending bool

	// 是否可读
	readable bool

	// 是否可写
	writeable bool

	// 是否已销毁
	destroyed bool

	// 超时时间 毫秒
	timeout int

	// 设置读缓冲区尺寸
	bufferSize int

	// 状态
	// connecting => opening
	// readable && writable => open
	// readable && !writable => readOnly
	// !readable && writable => writeOnly
	// else => closed
	readyState string

	// 是否初始化过
	inited bool

	// 暂停，挂起
	suspend bool
}

func NewTCPSocket() *TCPSocket {
	s := &TCPSocket{
		EventEmitter: *events.NewEventEmitter(),
	}
	s.Init()
	return s
}

func (this *TCPSocket) Init() {

	if !this.inited {
		this.inited = true
	} else {
		return
	}

	// 设置默认 读缓冲区尺寸
	if this.bufferSize <= 0 {
		this.bufferSize = 4096
	}

	// 在 SocketWaitGroup 中标记 进行中
	GetSocketWaitGroup().Add(1)

	// 默认 监听 一个 空 error 事件
	this.OnError(func(...interface{}) {})

	// 默认关闭事件
	this.OnClose(func(...interface{}) {
		this.SetCloseStatus()
		// 发送 结束事件
		this.EmitGo("_end")
	})

	// 默认_结束事件
	this.Once("_end", func(...interface{}) {
		this.Emit("end")
		// 结束后，从等待组中 标记为 done
		GetSocketWaitGroup().Done()
	})

}

// connect(port, host, options)
func (this *TCPSocket) Connect(args ...interface{}) {
	if len(args) >= 1 {
		this.port = args[0].(int)
	}
	if len(args) >= 2 {
		this.host = args[1].(string)
	}

	if this.Conn != nil {
		return
	}

	this.SetOpeningStatus()

	network := "tcp"
	address := this.host + ":" + strconv.Itoa(this.port)

	conn, err := net.Dial(network, address)
	if err != nil {
		try.Throw(SocketConnectException.NewThrow(err.Error()))
	}

	this.SetOpenStatus()

	this.Conn = conn
	go this.ConnectHandle()
}

func (this *TCPSocket) Reconnect() {
	//this.SetCloseStatus()
	this.Init()
	this.Connect(this.port, this.host)
}

func (this *TCPSocket) ConnectHandle() {

	buf := make([]byte, this.GetBufferSize())

	for {
		if this.suspend == false {
			func() {
				// 捕获异常， 发送error事件
				defer try.CatchUncaughtException(func(throwable try.Throwable) {
					this.EmitGo("error", throwable)
				})
				// 如果有超时时间，设置超时时间
				if this.timeout > 0 {
					timeoutDuration := time.Duration(this.timeout) * time.Millisecond
					this.Conn.SetReadDeadline(time.Now().Add(timeoutDuration))
				}

				cnt, err := this.Conn.Read(buf)
				if err != nil {
					this.suspend = true   // 暂停
					this.destroyed = true // 销毁
					if str.EndsWith(err.Error(), "use of closed network connection") {
						this.EmitGo("close", true) // 主动关闭
					} else if err.Error() == "EOF" { // 结束
						this.EmitGo("close", false) // 被动关闭
					} else if str.StartsWith(err.Error(), "read ") &&
						str.StartsWith(err.Error(), " i/o timeout") {
						this.Emit("timeout", "read")
					} else if str.StartsWith(err.Error(), "write ") &&
						str.StartsWith(err.Error(), " i/o timeout") {
						this.Emit("timeout", "write")
					} else {
						try.Throw(SocketReadException.NewThrow(err.Error()))
					}
				} else {
					this.Emit("data", buf[0:cnt])
				}
			}()
		} else {
			if this.destroyed == false {
				time.Sleep(50 * time.Millisecond)
			} else {
				break
			}
		}
	}
}

// write(data []byte, len int, index int)
func (this *TCPSocket) Write(args ...interface{}) int {
	if this.Conn != nil {
		var data []byte
		var length int
		var index int
		if len(args) >= 1 {
			data = args[0].([]byte)

			this.Conn.Write(data)
			return len(data)
		}
		if len(args) >= 2 {
			length = args[1].(int)
			index = 0
		}
		if len(args) >= 3 {
			index = args[2].(int)
		}

		newbytes := make([]byte, length)
		system.ByteArrayCopy(&data, index, &newbytes, 0, length)
		if this.timeout > 0 {
			timeoutDuration := time.Duration(this.timeout) * time.Millisecond
			this.Conn.SetWriteDeadline(time.Now().Add(timeoutDuration))
		}
		this.Conn.Write(newbytes)
		return length

	} else {
		return 0
	}
}

func (this *TCPSocket) OnData(listener func(...interface{})) {
	this.On("data", listener)
}

func (this *TCPSocket) OnError(listener func(...interface{})) {
	this.On("error", listener)
}

func (this *TCPSocket) OnClose(listener func(...interface{})) {
	this.Once("close", listener)
}

func (this *TCPSocket) OnEnd(listener func(...interface{})) {
	this.Once("end", listener)
}

func (this *TCPSocket) OnTimeout(listener func(...interface{})) {
	this.On("timeout", listener)
}

func (this *TCPSocket) On(t string, listener func(...interface{})) {
	this.EventEmitter.On(t, listener)
}

func (this *TCPSocket) Once(t string, listener func(...interface{})) {
	this.EventEmitter.Once(t, listener)
}

func (this *TCPSocket) LocalAddr() net.Addr {
	return this.Conn.LocalAddr()
}

func (this *TCPSocket) RemoteAddr() net.Addr {
	return this.Conn.RemoteAddr()
}

func (this *TCPSocket) Close() {
	if this.Conn != nil {
		err := this.Conn.Close()
		if err != nil {
			this.EmitGo("error", SocketCloseException.NewThrow(err.Error()))
		} else {
			this.EmitGo("close")
		}
	}
}

func (this *TCPSocket) SetOpeningStatus() {
	// 打开中
	this.connecting = true
	this.updateReadyStatus()
}

func (this *TCPSocket) SetOpenStatus() {
	// 状态变更为可读可写
	this.readable = true
	this.writeable = true
	this.updateReadyStatus()
}

func (this *TCPSocket) SetCloseStatus() {
	// 更新状态
	this.connecting = false
	this.readable = false
	this.writeable = false
	this.destroyed = true
	this.inited = false
	this.updateReadyStatus()
}

// 结束
func (this *TCPSocket) End() {
	if !this.destroyed {
		this.Close()
	}
}

// 暂停
func (this *TCPSocket) Pause() {
	this.suspend = true
}

func (this *TCPSocket) SetStreamTimeout(msecs int) {
	this.timeout = msecs
}

func (this *TCPSocket) GetStreamTimeout() int {
	return this.timeout
}

func (this *TCPSocket) SetBufferSize(size int) {
	this.bufferSize = size
}

func (this *TCPSocket) GetBufferSize() int {
	return this.bufferSize
}

//  // connecting => opening
//	// readable && writable => open
//	// readable && !writable => readOnly
//	// !readable && writable => writeOnly
//	// else => closed
func (this *TCPSocket) updateReadyStatus() {
	if this.connecting == true {
		this.readyState = "opening"
	} else if this.readable && this.writeable {
		this.readyState = "open"
	} else if this.readable && !this.writeable {
		this.readyState = "readOnly"
	} else if !this.readable && this.writeable {
		this.readyState = "writeOnly"
	} else {
		this.readyState = "closed"
	}
}

func (this *TCPSocket) SetPort(port int) {
	this.port = port
}

func (this *TCPSocket) GetPort() int {
	return this.port
}

func (this *TCPSocket) SetHost(host string) {
	this.host = host
}

func (this *TCPSocket) GetHost() string {
	return this.host
}

// 销毁
func (this *TCPSocket) Destroy() {
	go func() {
		this.EventEmitter.Destory()
		this.Conn = nil
	}()
}

// connect(port [, host])
func Connect(port int, host string) Socket {
	socket := NewTCPSocket()
	socket.Connect(port, host)
	return socket
}

// 创建 socket for server
func newSocketForServer(conn net.Conn) *TCPSocket {
	s := &TCPSocket{
		EventEmitter: *events.NewEventEmitter(),
		Conn:         conn,
	}
	s.Init()
	s.SetOpeningStatus()
	s.SetOpenStatus()
	return s
}
