package servers

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/miekg/dns"
)

var records = map[string]string{
	"test.service.": "192.168.0.2", // nslookup -port=5353 test.service 127.0.0.1
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			ip := records[q.Name]
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func StartServerDNS(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		log.Printf("The DNS server has started on port %d\n", portUDP_DNS)
		defer wg.Done()
		defer log.Printf("The DNS server has stopped\n")

		// attach request handler func
		dns.HandleFunc("service.", handleDnsRequest)

		// start server
		server := &dns.Server{Addr: ":" + strconv.Itoa(portUDP_DNS), Net: "udp"}
		err := server.ListenAndServe()
		defer server.Shutdown()
		if err != nil {
			log.Fatalf("Failed to start server: %s\n ", err.Error())
		}
	}()
}
