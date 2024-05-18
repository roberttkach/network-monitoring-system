package src

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestsError, requestLatency)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // Запоминаем время начала обработки запроса

	httpRequestsTotal.Inc()
	_, err := fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	if err != nil {
		httpRequestsError.Inc()
		http.Error(w, "Ошибка при форматировании ответа", http.StatusInternalServerError)
		return
	}

	elapsed := time.Since(start)              // Вычисляем время обработки запроса
	requestLatency.Observe(elapsed.Seconds()) // Добавляем значение в гистограмму задержек
}

func StartServerHTTP() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", handler)
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
		httpRequestsError.Inc() // Увеличиваем счетчик ошибок
	}
}
