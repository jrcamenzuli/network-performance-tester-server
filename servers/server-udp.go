package servers

import (
	"log"
	"net"
	"sync"
)

func StartServerUDP(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		log.Printf("The UDP server has started on port %d\n", portUDP)
		defer wg.Done()
		defer log.Printf("The UDP server has stopped\n")
		conn, err := net.ListenUDP("udp", &net.UDPAddr{
			Port: portUDP,
			IP:   net.ParseIP("0.0.0.0"),
		})
		if err != nil {
			panic(err)
		}

		defer conn.Close()

		message := make([]byte, 512)

		for {
			len, remote, err := conn.ReadFromUDP(message[:])
			if err != nil {
				continue
			}
			conn.WriteToUDP(message[:len], remote)
		}
	}()
}
