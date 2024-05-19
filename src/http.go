package src

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // Запоминаем время начала обработки запроса

	HttpRequestsTotal.Inc()
	_, err := fmt.Fprintf(w, "Hello %s!", r.URL.Path[1:])
	if err != nil {
		HttpRequestsError.Inc()
		http.Error(w, "Ошибка при форматировании ответа", http.StatusInternalServerError)
		return
	}

	elapsed := time.Since(start)              // Вычисляем время обработки запроса
	RequestLatency.Observe(elapsed.Seconds()) // Добавляем значение в гистограмму задержек
}

func StartServerHTTP() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
		HttpRequestsError.Inc() // Увеличиваем счетчик ошибок
	}
}
