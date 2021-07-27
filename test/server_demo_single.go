package main

import (
	"fmt"
	"github.com/muzin/go_rt/net"
	"github.com/muzin/go_rt/try"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {

	go func() {
		// 开启pprof，监听请求
		ip := "0.0.0.0:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
		}
	}()

	server := net.NewTCPServer()
	server.Listen(15010, "127.0.0.1")

	server.OnListen(func(args ...interface{}) {
		network := args[0].(string)
		address := args[1].(string)
		fmt.Printf("listen network: tls %v address: %v\n", network, address)
	})

	server.OnConnect(func(args ...interface{}) {

		socket := args[0].(net.Socket)
		fmt.Printf("connect %v %v\n", socket.LocalAddr(), socket.RemoteAddr())
		now := time.Now()

		//socket.SetStreamTimeout(100)

		socket.OnData(func(args ...interface{}) {
			databytes := args[0].([]byte)
			fmt.Printf("socket data: %v\n", string(databytes))
		})

		socket.OnClose(func(args ...interface{}) {
			since := time.Since(now)
			fmt.Printf("socket close: %v\n", since)
		})

		socket.OnError(func(args ...interface{}) {
			since := time.Since(now)
			throwable := args[0].(try.Throwable)
			fmt.Printf("socket error: %v %v\n", since, throwable)
		})

		socket.OnTimeout(func(args ...interface{}) {
			since := time.Since(now)
			fmt.Printf("socket timeout: %v \n", since)
		})

		socket.OnEnd(func(args ...interface{}) {
			since := time.Since(now)
			fmt.Printf("socket end: %v\n", since)

			//socket.Destroy()
			//socket = nil
		})

		count := 0

		for i := 0; i < 1; i++ {
			//for i := 0; i < 200; i++ {
			//time.Sleep(1 * time.Second)

			bytes := []byte("12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n12312345322345\n")

			socket.Write(bytes)

			count += len(bytes)

			//since := time.Since(now)

			//fmt.Printf("len: %v time: %v\n", count, since)
		}

		fmt.Printf("finish\n")

		time.Sleep(3 * time.Second)

		socket.Close()

	})

	server.OnClose(func(args ...interface{}) {
		fmt.Println("server close")
	})

	server.OnError(func(args ...interface{}) {
		throwable := args[0].(try.Throwable)
		fmt.Printf("server error: %v\n", throwable)
	})

	server.OnEnd(func(args ...interface{}) {
		fmt.Println("server end\n")
	})

	net.ExitAfterSocketEnd()

}
