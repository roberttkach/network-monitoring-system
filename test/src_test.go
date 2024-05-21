package test

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"net"
	"net/http"
	"net/http/httptest"
	"network-monitoring-system/src"
	"strings"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(src.Handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello !"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestStartServerHTTP(t *testing.T) {
	go src.StartServerHTTP()

	resp, err := http.Get("http://localhost:9090")
	if err != nil {
		t.Fatalf("Could not make GET request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}
}

func TestHandleDNSRequest(t *testing.T) {
	m := new(dns.Msg)
	m.SetQuestion("example.com.", dns.TypeA)

	a := new(dns.A)
	a.Hdr = dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 600}
	a.A = net.IPv4(127, 0, 0, 1)

	res := new(dns.Msg)
	res.SetReply(m)
	res.Answer = append(res.Answer, a)

	exp := new(dns.Msg)
	exp.SetReply(m)
	exp.Answer = append(exp.Answer, a)

	w := &fakeResponseWriter{}

	src.HandleDNSRequest(w, m)

	if w.msg.Answer[0].String() != exp.Answer[0].String() {
		t.Fatalf("Expected %s, got %s", exp.Answer[0].String(), w.msg.Answer[0].String())
	}
}

type fakeResponseWriter struct {
	msg *dns.Msg
}

func (w *fakeResponseWriter) LocalAddr() net.Addr       { return nil }
func (w *fakeResponseWriter) RemoteAddr() net.Addr      { return nil }
func (w *fakeResponseWriter) WriteMsg(m *dns.Msg) error { w.msg = m; return nil }
func (w *fakeResponseWriter) Write([]byte) (int, error) { return 0, nil }
func (w *fakeResponseWriter) Close() error              { return nil }
func (w *fakeResponseWriter) TsigStatus() error         { return nil }
func (w *fakeResponseWriter) TsigTimersOnly(bool)       {}
func (w *fakeResponseWriter) Hijack()                   {}

func TestCheckOpenPorts(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ln.Close()

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("Failed to get test server port: %v", err)
	}

	go src.CheckOpenPorts(port)
	time.Sleep(time.Second)
}

func TestHandleClientTCP(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		src.HandleClientTCP(conn)
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("test data"))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	timestr := strings.TrimSpace(string(buf[:n]))
	_, err = time.Parse(time.RFC3339, timestr)
	if err != nil {
		t.Errorf("response does not match time format: %s", timestr)
	}
}
func TestCheckErrorTCP(t *testing.T) {
	err := src.CheckErrorTCP(nil)
	if err != nil {
		t.Errorf("The function should not return an error for nil input")
	}

	err = src.CheckErrorTCP(fmt.Errorf("test error"))
	if err == nil {
		t.Errorf("The function should return an error for non-nil input")
	}
}

func TestStartServerTCP(t *testing.T) {
	go src.StartServerTCP()

	conn, err := net.Dial("tcp", "localhost:1200")
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}

func TestHandleClient(t *testing.T) {
	pc, err := net.ListenPacket("udp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer pc.Close()

	go func() {
		conn, err := net.Dial("udp", pc.LocalAddr().String())
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		_, err = conn.Write([]byte("test message"))
		if err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	src.HandleClient(pc.(*net.UDPConn))
}

func TestCheckError(t *testing.T) {
	src.CheckError(nil)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	err := errors.New("Test error")
	src.CheckError(err)
}
