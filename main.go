package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Response struct {
	StatusCode   int
	ResponseTime string
}

func ShowDegreeProgression(time string, degree int, maxRequest float32, done float32) {
	progression := 50
	progressionCount := degree / (100 / progression)
	fmt.Printf("\r(%v) [", time)
	for i := 0; i < progressionCount; i++ {
		fmt.Printf("#")
	}
	for i := 0; i < (progression - progressionCount); i++ {
		if i == 0 {
			fmt.Printf(">")
		}
		fmt.Printf("-")
	}
	fmt.Printf("] %v%v", degree, "%")
}

func Attack(wg *sync.WaitGroup, ch *chan int, client *http.Client, re chan Response) {
	defer wg.Done()
	req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
	rqStart := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		<-*ch
		return
	}
	rqEnd := time.Now()
	re <- Response{
		StatusCode:   resp.StatusCode,
		ResponseTime: rqEnd.Sub(rqStart).String(),
	}
	<-*ch
}

func main() {

	// 非同期数
	Channel := 100

	// Request数
	RequestCount := 10000

	// MaxIdleConns: DefaultTransportでは100になっている。0にすると無制限
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 0
	// MaxIdleConnsPerHost: デフォルト値は2。0にするとデフォルト値が使われる
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = Channel

	client := &http.Client{}

	// 並行処理するスレッド数を決める
	ch := make(chan int, Channel)

	// Responseを溜めておく
	result_ch := make(chan Response, RequestCount)

	// Requestを投げる時間測定
	requestStart := time.Now()

	// ひたすらRequestを投げる
	stash := 10
	wg := sync.WaitGroup{}
	for i := 0; i < RequestCount; i++ {
		ch <- 1
		wg.Add(1)
		degree := int((float32(i) / float32(RequestCount)) * 100)
		degreeP := degree / 5
		if degreeP != stash {
			stash = degreeP
			ShowDegreeProgression(time.Now().Sub(requestStart).String(), degree, float32(RequestCount), float32(i))
		}
		go Attack(&wg, &ch, client, result_ch)
	}
	wg.Wait()
	requestEnd := time.Now()
	ShowDegreeProgression(requestEnd.Sub(requestStart).String(), 100, float32(RequestCount), float32(RequestCount))
	// Response結果を取得
	Responses := make([]Response, RequestCount)
	i := 0
LOOP:
	for {
		select {
		case data := <-result_ch:
			Responses[i] = data
			i++
		default:
			break LOOP
		}
	}
	fmt.Println(requestEnd.Sub(requestStart).String(), len(Responses))
}
