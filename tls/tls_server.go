package tls

import (
	"crypto/tls"
	rt_net "github.com/muzin/go_rt/net"
	"github.com/muzin/go_rt/try"
	"net"
	"strconv"
	"time"
)

type TLSServer struct {
	rt_net.TCPServer

	tls.Config
}

func NewTLSServer() *TLSServer {
	t := &TLSServer{
		TCPServer: *rt_net.NewTCPServer(),
	}
	t.init()
	return t
}

func (this *TLSServer) init() {

	// 设置 创建 socket 处理函数
	this.SetNewSocketHandle(func(conn net.Conn) rt_net.Socket {
		return newSocketForTLSServer(conn)
	})

}

// 监听端口
// Listen(port, [host, [tls.Config]])
func (this *TLSServer) Listen(args ...interface{}) {

	if len(args) >= 1 {
		this.TCPServer.SetPort(args[0].(int))
	}
	if len(args) >= 2 {
		this.TCPServer.SetHost(args[1].(string))
	}
	if len(args) >= 3 {
		this.Config = args[2].(tls.Config)
	}

	var network = "tcp"
	var address = this.TCPServer.GetHost() + ":" + strconv.Itoa(this.TCPServer.GetPort())

	server, listenErr := tls.Listen(network, address, &this.Config)
	if listenErr != nil {
		try.Throw(rt_net.ServerListenException.NewThrow(listenErr.Error()))
	}

	// 加入 服务结束等待组中
	rt_net.GetSocketWaitGroup("tls_server Listen() WaitGroup add 1").Add(1)

	go func() {
		time.Sleep(10 * time.Millisecond)
		// 发送 监听事件
		this.EmitGo("listen", network, address)
	}()

	this.Listener = server
	// 处理 连接
	go this.ConnectHandle()

}

func (this *TLSServer) Destroy() {
	go func() {
		this.TCPServer.Destory()
	}()
}
