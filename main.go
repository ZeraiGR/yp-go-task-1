package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	url           = "http://srv.msk01.gigacorp.local/_stats"
	maxRetries   = 3
	sleepTime     = time.Second
	memoryLimit   = 80
	diskLimit     = 90
	bandwidthLimit = 90
	loadLimit     = 30
)

var currentRetries int

func main() {
	fetchStatistics()
}

func parseResponseBody(body io.Reader) ([7]int, error) {
	var data [7]int{}

	rawBody, err := io.ReadAll(body)

	if err != nil {
		return fetchWithMaxRetry()
	}

	rawDataArray := strings.Split(string(body), ",")
	
	if len(rawDataArray) != 7 {
		return fetchWithMaxRetry()
	}

	var isValidData = true
	for index, element := range rawDataArray {
		i, err := strconv.Atoi(element)

		if err != nil {
			isValidData = false
			break;
		}

		data[index] = i
	}

	if !isValidData {
		return fetchWithMaxRetry()
	}

	return data, nil
}

func analyzeData(data [7]int) {
	if data[0] > 30 {
		fmt.Printf("Load Average is too high: %d\n", data[0])
	}

	memoryUsage := data[2] * 100 / data[1]
	if memoryUsage > 80 {
		fmt.Printf("Memory usage too high: %d%%\n", memoryUsage)
	}

	if data[4] * 100 / data[3] > 90 {
		freeDiskSpace := (data[3] - data[4]) / (1024*1024)
		fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskSpace)
	}

	if data[6] * 100 / data[5] > 90 {
		freeBandwidth := (data[5] - data[6]) / 1000000
		fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeBandwidth)
	}
}

func fetchStatistics() {
	response, err := http.Get(url)

	if err != nil || response.StatusCode != http.StatusOK {
		return fetchWithMaxRetry()
	}

	data, err := parseResponseBody(response.Body)

	if err != nil {
		return fetchWithMaxRetry()
	}

	analyzeData(data)

	defer response.Body.Close()
}

func fetchWithMaxRetry() {
	currentRetries++

	if currentRetries >= maxRetries {
		fmt.Println("Unable to fetch server statistic")
		currentRetries = 0
	}

	time.Sleep(time.Second)
	fetchStatistics()
}