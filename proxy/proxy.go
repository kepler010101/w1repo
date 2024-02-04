package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
)

var (
	servers       = []string{"http://localhost:8081", "http://localhost:8082", "http://localhost:8083"}
	currentServer int32
	blocklist     []string
	limiter       = NewSlidingWindowLimiter(3, 60) //1- запросы, 2 - секунды
)

func getNextServer() *url.URL {
	index := atomic.AddInt32(&currentServer, 1) % int32(len(servers))
	serverURL, err := url.Parse(servers[index])
	if err != nil {
		log.Fatalf("Ошибка при разборе URL сервера: %v", err)
	}
	return serverURL
}

func logRequest(r *http.Request) {
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	log.Printf("Запрос: %s %s, IP: %s, Заголовки: %v, Тело: %s",
		r.Method, r.URL.Path, r.RemoteAddr, r.Header, string(bodyBytes))
}

func handleRequest(blocklist []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// получаем IP-адрес клиента
		clientIP := strings.Split(r.RemoteAddr, ":")[0]

		// чекаем спам лист
		if _, found := spamList[clientIP]; found {
			http.Error(w, "Ваш IP-адрес заблокирован из-за подозрений в спаме", http.StatusForbidden)
			return
		}

		// чекаем по лимиту запросов
		if !limiter.Allow(clientIP) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Заблокирован за превышение скорости"))
			return
		}

		// проверяем геолокацию IP-адреса
		allowed, country := isIPAllowed(clientIP)
		if !allowed {
			// возвращаем ошибку с сообщением о блокировке
			log.Printf("Доступ для IP-адреса %s (страна: %s) заблокирован по геолокации", clientIP, country)
			http.Error(w, "Ваш IP-адрес заблокирован из-за вашей геолокации", http.StatusForbidden)
			return
		}

		logRequest(r)

		// чекаем User-Agent
		if IsRequestBlocked(r, blocklist) {
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}

		// балансируем
		targetURL := getNextServer()
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	var err error
	blocklist, err = LoadBlocklist("scanners.txt")
	if err != nil {
		log.Fatalf("Ошибка при загрузке списка блокировки: %v", err)
	}

	http.HandleFunc("/", handleRequest(blocklist))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
