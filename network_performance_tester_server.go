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
	servers.StartServerDns(&wg, servers.UDP)
	servers.StartServerDns(&wg, servers.TCP)
	servers.StartServerPing(&wg)
	wg.Wait()
}
