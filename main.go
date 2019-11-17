package main

import (
	"fmt"
	"sync"
	"time"
)

type RequestResult struct {
	StatusCode   int
	ResponseTime string
}

func Hello(wg *sync.WaitGroup, ch *chan int, re chan RequestResult) {
	// defer wg.Done()
	// fmt.Println("hello world")
	// re <- RequestResult{
	// 	StatusCode:   200,
	// 	ResponseTime: "time",
	// }
	// <-*ch

	defer wg.Done()
	rqStart := time.Now()
	resp, _ := client.Do(req)
	rqEnd := time.Now()
	result[i] = RequestResult{
		StatusCode:   resp.StatusCode,
		ResponseTime: rqEnd.Sub(rqStart).String(),
	}
}

func main() {
	// 並行処理するスレッドの個数を決める
	ch := make(chan int, 2)
	//
	result_ch := make(chan RequestResult, 2)

	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		ch <- 1
		wg.Add(1)
		go Hello(&wg, &ch, result_ch)
		fmt.Println(<-result_ch)
	}
	wg.Wait()
}

// func Attack(ch chan []RequestResult, count int, client *http.Client, req *http.Request) {
// 	result := make([]RequestResult, count)
// 	for i := 0; i < count; i++ {
// 		rqStart := time.Now()
// 		resp, _ := client.Do(req)
// 		rqEnd := time.Now()
// 		result[i] = RequestResult{
// 			StatusCode:   resp.StatusCode,
// 			ResponseTime: rqEnd.Sub(rqStart).String(),
// 		}
// 		fmt.Println(resp.StatusCode)
// 	}
// 	ch <- result
// }

// func main() {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
// 	if err != nil {
// 	}
// 	for i := 0; i < 3; i++ {
// 		txtCh := make(chan []RequestResult, 200)
// 		go Attack(txtCh, 2, client, req)
// 		fmt.Println(<-txtCh)
// 	}
// 	time.Sleep(3 * time.Second)
// }
