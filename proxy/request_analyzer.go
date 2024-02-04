package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"
)

func LoadBlocklist(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var blocklist []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		blocklist = append(blocklist, line)
	}
	return blocklist, scanner.Err()
}

func IsRequestBlocked(r *http.Request, blocklist []string) bool {
	userAgent := r.Header.Get("User-Agent")
	for _, word := range blocklist {
		if strings.Contains(userAgent, word) {
			return true
		}
	}
	return false
}
