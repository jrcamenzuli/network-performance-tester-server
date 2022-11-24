package main

import (
	"sync"

	"github.com/jrcamenzuli/network-performance-tester-server/servers"
)

func main() {
	var wg sync.WaitGroup
	servers.StartServerHTTP(&wg)
	servers.StartServerHTTPS(&wg)
	servers.StartServerUDP(&wg)
	servers.StartServerDNS(&wg)
	servers.StartServerPing(&wg)
	wg.Wait()
}
