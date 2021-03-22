package tls

import (
	"crypto/tls"
	rt_net "github.com/muzin/go_rt/net"
	"net"
	"strconv"
)

type TLSSocket struct {
	rt_net.TCPSocket

	config *tls.Config
}

func NewTLSSocket() *TLSSocket {
	s := &TLSSocket{
		TCPSocket: *rt_net.NewTCPSocket(),
	}
	s.init()
	return s
}

func (this *TLSSocket) init() {

}

// connect(port, host, options)
func (this *TLSSocket) Connect(args ...interface{}) {
	if len(args) >= 1 {
		this.SetPort(args[0].(int))
	}
	if len(args) >= 2 {
		this.SetHost(args[1].(string))
	}
	if len(args) >= 3 {
		this.config = args[2].(*tls.Config)
	}

	if this.Conn != nil {
		return
	}

	this.SetOpeningStatus()

	network := "tcp"
	address := this.GetHost() + ":" + strconv.Itoa(this.GetPort())

	// 在 SocketWaitGroup 中标记 进行中
	rt_net.GetSocketWaitGroup("tls_socket Connect() WaitGroup add 1").Add(1)

	conn, err := tls.Dial(network, address, this.config)
	if err != nil {
		//try.Throw(rt_net.SocketConnectException.NewThrow(err.Error()))
		this.Emit("error", rt_net.SocketConnectException.NewThrow(err.Error()))
		this.Emit("close", true)
		return
	}

	this.SetOpenStatus()

	this.Conn = conn
	go this.ConnectHandle()
	this.Emit("connect", this)
}

func (this *TLSSocket) Reconnect() {
	this.TCPSocket.Init()
	this.init()

	rt_net.GetSocketWaitGroup("tls_socket Reconnect() Connect Listen WaitGroup add 1").Add(1)
	go func() {
		this.TCPSocket.Connect(this.GetPort(), this.GetHost(), this.config)
		rt_net.GetSocketWaitGroup("tls_socket Reconnect() Connect Listen WaitGroup done 1").Done()
	}()
}

func (this *TLSSocket) Destroy() {
	this.TCPSocket.Destory()
}

// connect(port [, host [, options]])
func ConnectTLS(port int, host string, options *tls.Config) rt_net.Socket {
	socket := NewTLSSocket()

	rt_net.GetSocketWaitGroup("tls_socket ConnectTLS() WaitGroup add 1").Add(1)
	go func() {
		socket.Connect(port, host, options)
		rt_net.GetSocketWaitGroup("tls_socket ConnectTLS() WaitGroup done 1").Done()
	}()

	return socket
}

func newSocketForTLSServer(conn net.Conn) *TLSSocket {

	tcpSocket := rt_net.NewTCPSocket()
	tcpSocket.Conn = conn
	tcpSocket.SetOpeningStatus()
	tcpSocket.SetOpenStatus()

	s := &TLSSocket{
		TCPSocket: *tcpSocket,
	}
	s.Init()

	// 在 SocketWaitGroup 中标记 进行中
	rt_net.GetSocketWaitGroup("tls_socket newSocketForTlsServer() WaitGroup add 1").Add(1)

	return s
}
