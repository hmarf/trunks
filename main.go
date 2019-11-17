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

func Attack(wg *sync.WaitGroup, ch *chan int, client *http.Client, re chan RequestResult) {
	defer wg.Done()
	req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
	rqStart := time.Now()
	resp, _ := client.Do(req)
	// defer resp.Body.Close()
	// fmt.Println(resp)
	rqEnd := time.Now()
	re <- RequestResult{
		StatusCode:   resp.StatusCode,
		ResponseTime: rqEnd.Sub(rqStart).String(),
	}
	<-*ch
}

func main() {

	// 非同期数
	Channel := 10

	// Request数
	RequestCount := 10000

	client := &http.Client{}

	// 並行処理するスレッド数を決める
	ch := make(chan int, Channel)

	// Responseを溜めておく
	result_ch := make(chan RequestResult, RequestCount)
	// _ = time.Now()
	// ひたすらRequestを投げる
	wg := sync.WaitGroup{}
	for i := 0; i < RequestCount; i++ {
		ch <- 1
		wg.Add(1)
		go Attack(&wg, &ch, client, result_ch)
	}
	wg.Wait()
	// _ = time.Now()

	// Response結果を取得
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
