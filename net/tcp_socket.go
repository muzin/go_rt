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

var SocketConnectException = try.DeclareException("SocketConnectException")
var SocketReadException = try.DeclareException("SocketReadException")
var SocketWriteException = try.DeclareException("SocketWriteException")
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

	// 是否关闭
	closed bool

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

	// 写 通道
	writeChannel chan []byte

	// 写通道关闭后 为 true
	writeChannelClosed bool

	//写 完毕
	writeChannelFinished bool

	// 读 通道
	readChannel chan []byte

	// 读通道关闭后 为 true
	readChannelClosed bool

	// 读 完毕
	readChannelFinished bool

	// 声明 的 处理器， 再重新加载时加载
	declareHanlders map[string][]func(...interface{})

	// 锁
	mu sync.Mutex

	// socket 地址缓存
	localAddrCache  net.Addr
	remoteAddrCache net.Addr
}

func NewTCPSocket() *TCPSocket {
	s := &TCPSocket{
		EventEmitter: *events.NewEventEmitter(),
	}
	s.Init()
	return s
}

func (this *TCPSocket) Init() {

	//if !this.inited {
	//	this.inited = true
	//} else {
	//	return
	//}

	this.connecting = false
	this.pending = false
	this.readable = false
	this.writeable = false
	this.destroyed = false
	this.readyState = ""
	this.inited = false
	this.suspend = false
	this.closed = false

	this.pending = true

	// 设置默认 读缓冲区尺寸
	if this.bufferSize <= 0 {
		this.bufferSize = 4096
	}

	// 读写 缓冲 通道
	this.readChannel = make(chan []byte, 100)
	this.writeChannel = make(chan []byte, 100)

	// 初始化 读写 缓冲区状态
	this.readChannelFinished = false
	this.readChannelFinished = false

	this.readChannelClosed = false
	this.writeChannelClosed = false

	// 如果没有声明过的处理器函数 创建
	if nil == this.declareHanlders {
		this.declareHanlders = make(map[string][]func(...interface{}))
	}

	// 如果有事件， 移除全部监听的事件，重新加载
	if len(this.EventNames()) > 0 {
		this.RemoveAllListener()
	}

	// 默认 监听 一个 空 error 事件
	this.On("error", func(...interface{}) {})

	// 监听 内部 超时事件
	this.Once("_timeout", func(args ...interface{}) {
		// 当 超时后 关闭 连接
		// 等待 连接 关闭完成
		// 发送超时事
		GetSocketWaitGroup("tcp_socket [event]_timeout WaitGroup add 1").Add(1)
		this.AddAppendListener("end", func(...interface{}) {
			// 当结束后 调 结束事件
			this.Emit("timeout", args[0])
			GetSocketWaitGroup("tcp_socket [event]_timeout WaitGroup done 1").Done()
		})
		this.Close()
	})

	// 默认关闭事件
	this.On("close", func(...interface{}) {
		// 如果已关闭，不在向下执行
		if this.writeChannelClosed && this.readChannelClosed {
			return
		}

		// 设置 关闭 状态
		this.SetCloseStatus()

		var wg sync.WaitGroup

		wg.Add(2)

		this.Once("writeChannelFinished", func(...interface{}) {
			this.writeChannelFinished = true
			wg.Done()
		})
		this.Once("readChannelFinished", func(...interface{}) {
			this.readChannelFinished = true
			wg.Done()
		})

		// 关闭 读写通道
		if !this.writeChannelClosed {
			this.writeChannelClosed = true
			close(this.writeChannel)
		}
		if !this.readChannelClosed {
			this.readChannelClosed = true
			close(this.readChannel)
		}

		wg.Wait()

		// 发送 结束事件
		this.EmitGo("_end")
	})

	// 默认_结束事件
	this.Once("_end", func(...interface{}) {

		// 添加 end 事件, end 执行完后， waitGroup 设置为 完成
		this.AddAppendListener("end", func(...interface{}) {
			// 结束后，从等待组中 标记为 done
			GetSocketWaitGroup("tcp_socket [event]_end WaitGroup done 1").Done()
		})
		// 发射 事件
		this.Emit("end")
	})

	// 如果 之前声明过 声明函数 直接加载
	if this.declareHanlders != nil && len(this.declareHanlders) > 0 {
		this.reloadDeclareHandlers()
	}

	// 处理 读写数据缓冲区 的 数据
	go this.writeConsumer()
	go this.readConsumer()

}

// connect(port, host, options)
func (this *TCPSocket) Connect(args ...interface{}) {
	if len(args) >= 1 {
		this.port = args[0].(int)
	}
	if len(args) >= 2 {
		this.host = args[1].(string)
	}

	//if this.Conn != nil {
	//	return
	//}

	this.SetOpeningStatus()

	network := "tcp"
	address := this.host + ":" + strconv.Itoa(this.port)

	// 在 SocketWaitGroup 中标记 进行中
	GetSocketWaitGroup("tcp_socket Connect() WaitGroup add 1").Add(1)

	conn, err := net.Dial(network, address)
	if err != nil {
		this.Emit("error", SocketConnectException.NewThrow(err.Error()))
		this.Emit("close", true)
		return
	}

	this.SetOpenStatus()

	this.Conn = conn

	go this.ConnectHandle()

	this.Emit("connect", this)

}

// 重连
func (this *TCPSocket) Reconnect() {
	//this.SetCloseStatus()
	this.Init()

	GetSocketWaitGroup("tcp_socket Reconnect() Connect WaitGroup add 1").Add(1)
	go func() {
		this.Connect(this.port, this.host)
		GetSocketWaitGroup("tcp_socket Reconnect() Connect WaitGroup done 1").Done()
	}()
}

// 连接处理
func (this *TCPSocket) ConnectHandle() {

	buf := make([]byte, this.GetBufferSize())

	for {
		if this.suspend == false {
			func() {
				// 捕获异常， 发送error事件
				defer try.CatchUncaughtException(func(throwable try.Throwable) {
					this.EmitGo("error", throwable)
					// 有错误 关闭 连接
					//this.Close()
				})()

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
						this.Emit("close", true) // 主动关闭
					} else if err.Error() == "EOF" { // 结束
						this.Emit("close", false) // 被动关闭
					} else if str.StartsWith(err.Error(), "read ") &&
						str.EndsWith(err.Error(), " i/o timeout") {
						this.Emit("_timeout", "read")
					} else if str.StartsWith(err.Error(), "write ") &&
						str.EndsWith(err.Error(), " i/o timeout") {
						this.Emit("_timeout", "write")
					} else {
						try.Throw(SocketReadException.NewThrow(err.Error()))
					}
				} else {
					//this.Emit("data", buf[0:cnt])
					// 将 数据 写入 读缓冲区
					if !this.readChannelClosed {
						this.readChannel <- append(make([]byte, 0), buf[0:cnt]...)
					}
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

func (this *TCPSocket) Write(args ...interface{}) int {
	var data []byte
	var length int
	var index int
	if len(args) >= 1 {
		data = args[0].([]byte)
		length = len(data)
	}
	if len(args) >= 2 {
		length = args[1].(int)
		index = 0
	}
	if len(args) >= 3 {
		index = args[2].(int)
	}

	// 如果 没有 关闭加入到 写缓冲区
	if !this.writeChannelClosed {
		this.writeChannel <- data[index:(index + length)]
		data = nil
	}

	//if !this.closed {

	//this.writeChannel <- append(make([]byte, 0), data[index:(index + length)]...)
	//}

	//cnt, err := this.write(data[index:(index + length)])

	return length
}

// write(data []byte, len int, index int)
func (this *TCPSocket) write(data []byte) (int, error) {
	if this.Conn != nil && this.suspend == false {
		if this.timeout > 0 {
			timeoutDuration := time.Duration(this.timeout) * time.Millisecond
			this.Conn.SetWriteDeadline(time.Now().Add(timeoutDuration))
		}
		cnt, err := this.Conn.Write(data)

		//fmt.Printf("Conn Write: from: %v dest: %v cnt: %v, err: %v, data:%v\n",
		//	this.LocalAddr(), this.RemoteAddr(), cnt, err, string(data))

		if nil != err {
			this.Emit("error", SocketWriteException.NewThrow(err.Error()))
			// 有错误 关闭 连接
			this.Close()
		}
		return cnt, err
	} else {
		return 0, nil
	}
}

func (this *TCPSocket) writeConsumer() {
	for {
		if this.pending {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		if this.connecting && !this.suspend {
			data, isOpen := <-this.writeChannel
			if isOpen {
				if !this.closed {
					_, err := this.write(data)
					if nil != err {
						for v := range this.writeChannel {
							func(v []byte) {}(v)
						}
						this.EmitGo("writeChannelFinished")
						break
					}
				} else {
					for v := range this.writeChannel {
						func(v []byte) {}(v)
					}
					this.EmitGo("writeChannelFinished")
					break
				}
				data = nil
			} else {
				this.EmitGo("writeChannelFinished")
				break
			}
		} else {
			if !this.closed {
				time.Sleep(10 * time.Millisecond)
			} else {
				for v := range this.writeChannel {
					func(v []byte) {}(v)
				}
				this.EmitGo("writeChannelFinished")
				break
			}
		}

	}
}

func (this *TCPSocket) readConsumer() {
	for {
		if this.pending {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if this.connecting && !this.suspend {
			data, isOpen := <-this.readChannel
			if isOpen {
				if !this.closed {
					this.Emit("data", data)
				} else {
					for v := range this.readChannel {
						func(v []byte) {}(v)
					}
					this.EmitGo("readChannelFinished")
					break
				}
			} else {
				for v := range this.readChannel {
					func(v []byte) {}(v)
				}
				this.EmitGo("readChannelFinished")
				break
			}
		} else {
			if !this.closed {
				time.Sleep(10 * time.Millisecond)
			} else {
				for v := range this.readChannel {
					func(v []byte) {}(v)
				}
				this.EmitGo("readChannelFinished")
				break
			}
		}
	}
}

func (this *TCPSocket) OnConnect(listener func(...interface{})) {
	this.mu.Lock()
	defer this.mu.Unlock()

	eventName := "connect"
	_, ok := this.declareHanlders[eventName]
	if !ok {
		this.declareHanlders[eventName] = make([]func(...interface{}), 0)
	}
	this.declareHanlders[eventName] = append(this.declareHanlders[eventName], listener)

	this.Once(eventName, listener)
}

func (this *TCPSocket) OnData(listener func(...interface{})) {
	this.mu.Lock()
	defer this.mu.Unlock()

	eventName := "data"
	_, ok := this.declareHanlders[eventName]
	if !ok {
		this.declareHanlders[eventName] = make([]func(...interface{}), 0)
	}
	this.declareHanlders[eventName] = append(this.declareHanlders[eventName], listener)

	this.On(eventName, listener)
}

func (this *TCPSocket) OnError(listener func(...interface{})) {
	this.mu.Lock()
	defer this.mu.Unlock()

	eventName := "error"
	_, ok := this.declareHanlders[eventName]
	if !ok {
		this.declareHanlders[eventName] = make([]func(...interface{}), 0)
	}
	this.declareHanlders[eventName] = append(this.declareHanlders[eventName], listener)

	this.On(eventName, listener)
}

func (this *TCPSocket) OnClose(listener func(...interface{})) {
	this.mu.Lock()
	defer this.mu.Unlock()

	eventName := "close"
	_, ok := this.declareHanlders[eventName]
	if !ok {
		this.declareHanlders[eventName] = make([]func(...interface{}), 0)
	}
	this.declareHanlders[eventName] = append(this.declareHanlders[eventName], listener)

	this.Once(eventName, listener)
}

func (this *TCPSocket) OnEnd(listener func(...interface{})) {
	this.mu.Lock()
	defer this.mu.Unlock()

	eventName := "end"
	_, ok := this.declareHanlders[eventName]
	if !ok {
		this.declareHanlders[eventName] = make([]func(...interface{}), 0)
	}
	this.declareHanlders[eventName] = append(this.declareHanlders[eventName], listener)

	this.Once(eventName, listener)
}

func (this *TCPSocket) OnTimeout(listener func(...interface{})) {
	this.mu.Lock()
	defer this.mu.Unlock()

	eventName := "timeout"
	_, ok := this.declareHanlders[eventName]
	if !ok {
		this.declareHanlders[eventName] = make([]func(...interface{}), 0)
	}
	this.declareHanlders[eventName] = append(this.declareHanlders[eventName], listener)

	this.Once(eventName, listener)
}

func (this *TCPSocket) On(t string, listener func(...interface{})) {
	this.EventEmitter.On(t, listener)
}

func (this *TCPSocket) Once(t string, listener func(...interface{})) {
	this.EventEmitter.Once(t, listener)
}

func (this *TCPSocket) LocalAddr() net.Addr {
	if this.localAddrCache == nil {
		if this.Conn != nil {
			this.localAddrCache = this.Conn.LocalAddr()
		}
	}
	return this.localAddrCache
}

func (this *TCPSocket) RemoteAddr() net.Addr {
	if this.remoteAddrCache == nil {
		if this.Conn != nil {
			this.remoteAddrCache = this.Conn.RemoteAddr()
		}
	}
	return this.remoteAddrCache
}

func (this *TCPSocket) Close() {
	if this.closed {
		return
	}
	this.closed = true
	if this.Conn != nil {
		err := this.Conn.Close()
		if err == nil {
			//this.Emit("close", true)
		} else {
			// this.EmitGo("error", SocketCloseException.NewThrow(err.Error()))
		}
	}
}

func (this *TCPSocket) IsClose() bool {
	return this.closed
}

func (this *TCPSocket) SetOpeningStatus() {
	// 打开中
	this.pending = false
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
	this.closed = true
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

// 重新加载声明过的事件
func (this *TCPSocket) reloadDeclareHandlers() {
	for k, v := range this.declareHanlders {
		if nil != v {
			for i := 0; i < len(v); i++ {
				if nil != v[i] {
					this.On(k, v[i])
				}
			}
		}
	}
}

// 销毁
func (this *TCPSocket) Destroy() {
	go func() {
		// 关闭 读写通道
		if !this.writeChannelClosed {
			close(this.writeChannel)
		}
		if !this.readChannelClosed {
			close(this.readChannel)
		}

		this.EventEmitter.Destory()
		this.Conn = nil
		this.writeChannel = nil
		this.readChannel = nil

		this.declareHanlders = nil
	}()
}

// connect(port [, host])
func Connect(port int, host string) Socket {
	socket := NewTCPSocket()

	GetSocketWaitGroup("tcp_socket Connect() Listen WaitGroup add 1").Add(1)
	go func() {
		socket.Connect(port, host)
		GetSocketWaitGroup("tcp_socket Connect() Listen WaitGroup done 1").Done()
	}()

	return socket
}

// 创建 socket for server
func newSocketForServer(conn net.Conn) *TCPSocket {
	s := &TCPSocket{
		EventEmitter: *events.NewEventEmitter(),
		Conn:         conn,
	}
	s.Init()

	// 在 SocketWaitGroup 中标记 进行中
	GetSocketWaitGroup("tcp_socket newSocketForServer() WaitGroup add 1").Add(1)

	s.SetOpeningStatus()
	s.SetOpenStatus()
	return s
}
