package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Response struct {
	StatusCode   int
	ResponseTime time.Duration
}

type CountStatusCode struct {
	StatusCode int
	Count      int
}

func ShowDegreeProgression(time time.Duration, degree int, maxRequest float32, done float32) {
	progression := 50
	progressionCount := degree / (100 / progression)
	fmt.Printf("\r[%02d:%02d:%02d] [", int(time.Hours()), int(time.Minutes())%60, int(time.Seconds())%60)
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
		ResponseTime: rqEnd.Sub(rqStart),
	}
	<-*ch
}

func main() {

	// 非同期数
	Channel := 10

	// Request数
	RequestCount := 100

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
			ShowDegreeProgression(time.Now().Sub(requestStart), degree, float32(RequestCount), float32(i))
		}
		go Attack(&wg, &ch, client, result_ch)
	}
	wg.Wait()
	requestEnd := time.Now()
	ShowDegreeProgression(requestEnd.Sub(requestStart), 100, float32(RequestCount), float32(RequestCount))
	fmt.Println("\n")

	// Response結果を取得
	responses := make([]Response, RequestCount)
	// メジャーなステータスコードを初期値とする
	// https://www.sakurasaku-labo.jp/blogs/status-code-basic-knowledgess
	countStatusCode := map[int]int{
		200: 0, // 成功
		301: 0, // 恒久的にページが移動している
		302: 0, // 一時的にページが移動している
		400: 0, // リクエストが不正
		401: 0, // 要認証
		403: 0, // アクセス禁止
		404: 0, // アクセスができない
		500: 0, // サーバエラー
		503: 0, // サービス利用不可
	}

	// Latency
	maxLatency := time.Duration(0)
	minLatency := requestEnd.Sub(requestStart)
	meanLatency := time.Duration(0)
LOOP:
	for i := 0; ; {
		select {
		case data := <-result_ch:
			meanLatency += data.ResponseTime
			// 待機時間　max, min
			if data.ResponseTime > maxLatency {
				maxLatency = data.ResponseTime
			}
			if data.ResponseTime < minLatency {
				minLatency = data.ResponseTime
			}
			// Response の Status Code を数える
			v, ok := countStatusCode[data.StatusCode]
			if ok {
				v++
				countStatusCode[data.StatusCode] = v
			} else {
				countStatusCode[data.StatusCode] = 1
			}
			responses[i] = data
			i++
		default:
			break LOOP
		}
	}

	fmt.Printf("Latency:\n Total: %v\n Max:   %v\n Min:   %v\n Ave:   %v\n",
		requestEnd.Sub(requestStart), maxLatency, minLatency, meanLatency/time.Duration(RequestCount))

}
