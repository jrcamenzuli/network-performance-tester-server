package servers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

type TransportProtocol int64

const (
	UDP TransportProtocol = iota
	TCP
)

// Usage:
// nslookup -port=53 test.service 127.0.0.1
// nslookup -port=53 -vc test.service 127.0.0.1
var records = map[string]string{
	"test.service.": "192.168.0.2",
}

func (t TransportProtocol) toString() string {
	switch t {
	case UDP:
		return "UDP"
	case TCP:
		return "TCP"
	default:
		return "UNKNOWN"
	}
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			ip := records[q.Name]
			log.Printf("A DNS question for \"%s\" is answered with \"%s\"\n", q.Name, ip)
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func StartServerDns(wg *sync.WaitGroup, t TransportProtocol) {
	wg.Add(1)
	go func(t TransportProtocol) {
		port := 53
		switch t {
		case UDP:
			port = portUDP_DNS
		case TCP:
			port = portTCP_DNS
		}
		log.Printf("The DNS over %s server has started on port %d\n", t.toString(), port)
		defer wg.Done()
		defer log.Printf("The DNS over %s server has stopped\n", t.toString())

		// attach request handler func
		dns.HandleFunc("service.", func(w dns.ResponseWriter, r *dns.Msg) {
			m := new(dns.Msg)
			m.SetReply(r)
			m.Compress = false

			switch r.Opcode {
			case dns.OpcodeQuery:
				parseQuery(m)
			}

			w.WriteMsg(m)
		})

		// start server
		server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: strings.ToLower(t.toString())}
		err := server.ListenAndServe()
		defer server.Shutdown()
		if err != nil {
			log.Fatalf("Failed to start the DNS over %s server: %s\n ", t.toString(), err.Error())
		}
	}(t)
}
