package tls

import (
	"crypto/tls"
	rt_net "github.com/muzin/go_rt/net"
	"github.com/muzin/go_rt/try"
	"net"
	"strconv"
)

type TLSSocket struct {
	rt_net.TCPSocket

	tls.Config
}

func NewTLSSocket() *TLSSocket {
	s := &TLSSocket{
		TCPSocket: *rt_net.NewTCPSocket(),
	}
	s.Init()
	return s
}

func (this *TLSSocket) Init() {

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
		this.Config = args[2].(tls.Config)
	}

	if this.Conn != nil {
		return
	}

	this.SetOpeningStatus()

	network := "tcp"
	address := this.GetHost() + ":" + strconv.Itoa(this.GetPort())

	conn, err := tls.Dial(network, address, &this.Config)
	if err != nil {
		try.Throw(rt_net.SocketConnectException.NewThrow(err.Error()))
	}

	this.SetOpenStatus()

	this.Conn = conn
	go this.ConnectHandle()
}

func (this *TLSSocket) Reconnect() {
	this.TCPSocket.Init()
	this.Init()
	this.TCPSocket.Connect(this.GetPort(), this.GetHost(), this.Config)
}

func (this *TLSSocket) Destroy() {
	go func() {
		this.TCPSocket.Destory()
	}()
}

// connect(port [, host [, options]])
func ConnectTLS(port int, host string, options tls.Config) rt_net.Socket {
	socket := NewTLSSocket()
	socket.Connect(port, host, options)
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
	return s
}
