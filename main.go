package main

import (
	"fmt"
	"net/http"
	"time"
)

type RequestResult struct {
	StatusCode   int
	ResponseTime string
}

func Attack(ch chan []RequestResult, count int, client *http.Client, req *http.Request) {
	result := make([]RequestResult, count)
	for i := 0; i < count; i++ {
		rqStart := time.Now()
		resp, _ := client.Do(req)
		rqEnd := time.Now()
		result[i] = RequestResult{
			StatusCode:   resp.StatusCode,
			ResponseTime: rqEnd.Sub(rqStart).String(),
		}
		fmt.Println(resp.StatusCode)
	}
	ch <- result
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
	}
	txtCh := make(chan []RequestResult, 200)
	go Attack(txtCh, 2, client, req)
	fmt.Println(<-txtCh)
	time.Sleep(3 * time.Second)

	// client := &http.Client{}
	// req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	// if err != nil {
	// 	// handle error
	// }
	// resp, err := client.Do(req)
	// if err != nil {
	// 	// handle error
	// }
	// defer resp.Body.Close()

	// if err != nil {
	// 	// handle error
	// }
	// fmt.Println(resp.StatusCode)
}
