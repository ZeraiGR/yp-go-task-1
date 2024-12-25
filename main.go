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
 url          = "http://srv.msk01.gigacorp.local/_stats"
 maxAttempts  = 3
 sleepTime    = time.Second
 memoryLimit  = 80
 diskLimit    = 90
 bandwidthLimit = 90
 loadLimit    = 30
)

var failedAttempts int

func main() {
 for {
  if !fetchAndProcessStats() {
   handleFailedAttempt()
  }
  time.Sleep(sleepTime)
 }
}

// Fetch and process server statistics
func fetchAndProcessStats() bool {
 response, err := http.Get(url)
 if err != nil || response.StatusCode != http.StatusOK {
  return false
 }
 defer response.Body.Close()

 data, err := parseResponseBody(response.Body)
 if err != nil {
  return false
 }

 analyzeData(data)
 return true
}

// Parse the response body into an integer array
func parseResponseBody(body io.Reader) ([7]int, error) {
 var data [7]int

 rawBody, err := io.ReadAll(body)
 if err != nil {
  return data, err
 }

 rawDataArray := strings.Split(strings.TrimSpace(string(rawBody)), ",")
 if len(rawDataArray) != 7 {
  return data, fmt.Errorf("unexpected data length")
 }

 for index, element := range rawDataArray {
  i, err := strconv.Atoi(strings.TrimSpace(element))
  if err != nil {
   return data, err
  }
  data[index] = i
 }
 return data, nil
}

// Analyze the parsed data and print warnings if thresholds are exceeded
func analyzeData(data [7]int) {
 // Check Load Average
 if data[0] > loadLimit {
  fmt.Printf("Load Average is too high: %d\n", data[0])
 }

 // Check Memory usage
 memoryUsage := data[2] * 100 / data[1]
 if memoryUsage > memoryLimit {
  fmt.Printf("Memory usage too high: %d%%\n", memoryUsage)
 }

 // Check free disk space
 freeDiskSpace := (data[3] - data[4]) / (1024 * 1024)
 if (data[3]-data[4])*100/data[3] < (100 - diskLimit) {
  fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskSpace)
 }

 // Check network bandwidth usage
 freeBandwidth := (data[5] - data[6]) * 8 / 1000000 // Convert bytes to megabits
 if (data[5]-data[6])*100/data[5] < (100 - bandwidthLimit) {
  fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeBandwidth)
 }
}

// Handle a failed attempt and print a message if necessary
func handleFailedAttempt() {
 failedAttempts++
 if failedAttempts >= maxAttempts {
  fmt.Println("Unable to fetch server statistic")
  failedAttempts = 0 // Reset the counter
 }
}