package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RequestResult struct {
	StatusCode   int
	ResponseTime string
}

func Hello(wg *sync.WaitGroup, ch *chan int, client *http.Client,
	req *http.Request, re chan RequestResult) {
	fmt.Println("aa")
	time.Sleep(time.Second)
	defer wg.Done()
	rqStart := time.Now()
	resp, _ := client.Do(req)
	rqEnd := time.Now()
	re <- RequestResult{
		StatusCode:   resp.StatusCode,
		ResponseTime: rqEnd.Sub(rqStart).String(),
	}
	<-*ch
}

func main() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080", nil)

	// 並行処理するスレッド数を決める
	ch := make(chan int, 2)

	//
	result_ch := make(chan RequestResult, 2000)

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		ch <- 1
		wg.Add(1)
		go Hello(&wg, &ch, client, req, result_ch)
	}

	wg.Wait()

LOOP:
	for {
		select {
		case aa := <-result_ch:
			fmt.Println(aa)
		default:
			break LOOP
		}
	}

}
