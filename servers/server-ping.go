package servers

import (
	"log"
	"net"
	"sync"
)

func StartServerPing(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		log.Printf("The Ping server has started on port %d\n", portPing)
		defer wg.Done()
		defer log.Printf("The Ping server has stopped\n")
		conn, err := net.ListenUDP("udp", &net.UDPAddr{
			Port: portPing,
			IP:   net.ParseIP("0.0.0.0"),
		})
		if err != nil {
			panic(err)
		}

		defer conn.Close()

		message := make([]byte, 1)

		for {
			len, remote, err := conn.ReadFromUDP(message[:])
			if err != nil {
				continue
			}
			conn.WriteToUDP(message[:len], remote)
		}
	}()
}
