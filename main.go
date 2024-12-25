package main

import (
 "fmt"
 "io"
 "net/http"
 "strconv"
 "strings"
 "time"
)

// Global variable to keep track of failed attempts
var failedAttempts = 0

func main() {
 for {
  handleResponse(fetchServerStatistics())
  time.Sleep(time.Second)
 }
}

// Fetch server statistics from the given URL
func fetchServerStatistics() (*http.Response, error) {
 response, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
 if err != nil {
  return nil, err
 }
 return response, nil
}

// Handle server response
func handleResponse(response *http.Response, err error) {
 if err != nil || response.StatusCode != 200 {
  handleFailedAttempt()
  return
 }
 defer response.Body.Close()

 data, err := parseResponseBody(response.Body)
 if err != nil {
  handleFailedAttempt()
  return
 }

 analyzeData(data)
}

// Attempt to parse response body into an integer array
func parseResponseBody(body io.Reader) ([7]int, error) {
 var data [7]int

 rawBody, err := io.ReadAll(body)
 if err != nil {
  return data, err
 }

 rawDataArray := strings.Split(string(rawBody), ",")
 if len(rawDataArray) != 7 {
  return data, fmt.Errorf("unexpected data length")
 }

 for index, element := range rawDataArray {
  i, err := strconv.Atoi(element)
  if err != nil {
   return data, err
  }
  data[index] = i
 }
 return data, nil
}

// Analyze and report based on parsed data
func analyzeData(data [7]int) {
 // Check Load Average
 if data[0] > 30 {
  fmt.Printf("Load Average is too high: %d\n", data[0])
 }

 // Check Memory usage
 memoryUsage := data[2] * 100 / data[1]
 if memoryUsage > 80 {
  fmt.Printf("Memory usage too high: %d%%\n", memoryUsage)
 }

 // Check free disk space
 if data[4] * 100 / data[3] > 90 {
  freeDiskSpace := (data[3] - data[4]) / (1024 * 1024)
  fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskSpace)
 }

 // Check network bandwidth
	if data[6] * 100 / data[5] > 90 {
  	freeBandwidth := (data[5] - data[6]) / 1000000
  	fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeBandwidth)
	}
}

// Handle failed attempts by tracking and resetting counters
func handleFailedAttempt() {
	failedAttempts++
	if failedAttempts >= 3 {
  failedAttempts = 0
  fmt.Println("Unable to fetch server statistic")
	}
}