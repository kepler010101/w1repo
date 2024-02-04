package main

import (
	"encoding/csv"
	"log"
	"net"
	"net/http"
)

type IPRange struct {
	Start   string
	End     string
	Country string
}

var ipRanges []IPRange

func downloadAndParseDB(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка при скачивании базы данных: %v", err)
	}
	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("Ошибка при чтении файла CSV: %v", err)
	}

	for _, record := range records {
		if len(record) == 3 {
			ipRange := IPRange{
				Start:   record[0],
				End:     record[1],
				Country: record[2],
			}
			ipRanges = append(ipRanges, ipRange)
		}
	}
}

func ipToUint32(ipStr string) uint32 {
	ip := net.ParseIP(ipStr)
	ip = ip.To4()
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

func isIPAllowed(ipStr string) (bool, string) {
	ipUint := ipToUint32(ipStr)
	for _, ipRange := range ipRanges {
		startUint := ipToUint32(ipRange.Start)
		endUint := ipToUint32(ipRange.End)
		if ipUint >= startUint && ipUint <= endUint {
			return true, ipRange.Country
		}
	}
	return false, "" 
