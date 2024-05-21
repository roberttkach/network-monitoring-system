package src

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	HttpRequestsTotal.Inc()
	_, err := fmt.Fprintf(w, "Hello %s!", r.URL.Path[1:])
	if err != nil {
		HttpRequestsError.Inc()
		http.Error(w, "Ошибка при форматировании ответа", http.StatusInternalServerError)
		return
	}

	elapsed := time.Since(start)
	RequestLatency.Observe(elapsed.Seconds())
}

func StartServerHTTP() {
	http.HandleFunc("/", Handler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
		HttpRequestsError.Inc()
	}
}
