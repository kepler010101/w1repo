package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var spamList = make(map[string]bool)

func updateSpamList() {
	for {
		resp, err := http.Get("https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt")
		if err != nil {
			log.Printf("Ошибка при загрузке спам-листа: %v", err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Ошибка при чтении тела ответа: %v", err)
			continue
		}

		lines := strings.Split(string(body), "\n")
		newSpamList := make(map[string]bool)
		for _, line := range lines {
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if fields := strings.Fields(line); len(fields) > 0 && !strings.HasSuffix(line, " 1") && !strings.HasSuffix(line, " 2") {
				newSpamList[fields[0]] = true
			}
		}

		spamList = newSpamList
		log.Println("Спам-лист успешно обновлен")

		time.Sleep(24 * time.Hour) // обновлять список каждые 24 часа
	}
}
