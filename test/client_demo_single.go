package main

import (
	"fmt"
	net "github.com/muzin/go_rt/net"
	"github.com/muzin/go_rt/try"
	"time"
)

func main() {

	for i := 0; i < 1; i++ {
		go batchTest(i)
	}

	time.Sleep(50 * time.Millisecond)

	net.ExitAfterSocketEnd()

}

func batchTest(i int) {

	connect := net.Connect(15010, "127.0.0.1")

	count := 0

	now := time.Now()

	//connect.SetStreamTimeout(120)

	connect.OnConnect(func(args ...interface{}) {
		fmt.Println("connect")
	})

	connect.OnError(func(args ...interface{}) {
		throwable := args[0].(try.Throwable)
		fmt.Printf("%v socket error: %v\n", i, throwable)
	})

	connect.OnClose(func(args ...interface{}) {
		fmt.Printf("%v connect close: \n", i)
	})

	connect.OnTimeout(func(args ...interface{}) {
		val := args[0].(string)
		fmt.Printf("%v connect timeout: %v\n", i, val)

		//fmt.Printf("%v connect reconnect: %v\n", i, val)
		//connect.Reconnect()
	})

	connect.OnData(func(args ...interface{}) {
		bytes := args[0].([]byte)
		count += 1
		fmt.Printf("data %v: len: %v \n", count, len(bytes))
	})

	connect.OnEnd(func(args ...interface{}) {

		since := time.Since(now)
		fmt.Printf("%v connect end: %v\n", i, since)

	})

}
