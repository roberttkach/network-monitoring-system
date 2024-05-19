package src

import (
	"net"

	"strconv"
	"time"
)

func CheckOpenPorts() {
	for port := 1; port <= 65535; port++ {
		go func(port int) {
			start := time.Now() // Запоминаем время начала обработки запроса

			_, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
			if err == nil {
				OpenPorts.WithLabelValues(strconv.Itoa(port)).Inc()
			}

			elapsed := time.Since(start)              // Вычисляем время обработки запроса
			RequestLatency.Observe(elapsed.Seconds()) // Добавляем значение в гистограмму задержек

			PortChecks.Inc()    // Увеличиваем счетчик на 1
			TotalRequests.Inc() // Увеличиваем общий счетчик на 1
		}(port)
	}
}
