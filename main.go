package main

import (
	"fmt"
	"log"
	"net/http"
)

func startServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ответ от сервера на порту %s", port)
	})

	log.Printf("Сервер запущен на http://localhost%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Ошибка запуска сервера на порту %s: %v", port, err)
	}
}

func main() {
	go startServer(":8081")
	go startServer(":8082")
	go startServer(":8083")

	select {}
}
