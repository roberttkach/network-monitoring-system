package src

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	requests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
	)
)

func init() {
	prometheus.MustRegister(requests)
}

func handler(w http.ResponseWriter, r *http.Request) {
	requests.Inc()
	_, err := fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	if err != nil {
		http.Error(w, "Ошибка при форматировании ответа", http.StatusInternalServerError)
		return
	}
}

func StartServerHTTP() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", handler)
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
