package src

import (
	"github.com/prometheus/client_golang/prometheus"
	"net"
	"strconv"
	"time"
)

func init() {
	prometheus.MustRegister(openPorts, portChecks, TotalRequests, requestLatency) // Регистрируем новые счетчики
}

func CheckOpenPorts() {
	for port := 1; port <= 65535; port++ {
		go func(port int) {
			start := time.Now() // Запоминаем время начала обработки запроса

			_, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
			if err == nil {
				openPorts.WithLabelValues(strconv.Itoa(port)).Inc()
			}

			elapsed := time.Since(start)              // Вычисляем время обработки запроса
			requestLatency.Observe(elapsed.Seconds()) // Добавляем значение в гистограмму задержек

			portChecks.Inc()    // Увеличиваем счетчик на 1
			TotalRequests.Inc() // Увеличиваем общий счетчик на 1
		}(port)
	}
}
