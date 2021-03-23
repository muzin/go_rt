package net

import (
	"sync"
	"testing"
)

func TestNewServer(t *testing.T) {
	t.Run("TestNewServer", func(t *testing.T) {

		var wg sync.WaitGroup
		wg.Add(1)

		server := NewTCPServer()

		server.OnListen(func(args ...interface{}) {
			network := args[0].(string)
			address := args[1].(string)
			t.Logf("listen network: %v address: %v\n", network, address)

			wg.Done()
		})

		server.Listen(15000, "127.0.0.1")

		wg.Wait()
		//time.Sleep(10 * time.Second)

	})
}
