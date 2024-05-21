package src

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"time"
)

func RecordMetrics() {
	go func() {
		for {
			TotalRequests.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

func HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	start := time.Now()

	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 600},
			A:   net.IPv4(127, 0, 0, 1),
		})
	}
	w.WriteMsg(&msg)

	TotalRequests.Inc()
	HttpRequestsTotal.Inc()

	elapsed := time.Since(start)
	RequestLatency.Observe(elapsed.Seconds())
}

func StartServerDNS() {
	RecordMetrics()

	dns.HandleFunc("example.com.", HandleDNSRequest)

	server := &dns.Server{Addr: ":55", Net: "udp"}
	log.Printf("Starting at %s\n", server.Addr)

	err := server.ListenAndServe()
	defer server.Shutdown()

	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
		// HttpRequestError.Inc() // Увеличиваем счетчик ошибок
	}
}
