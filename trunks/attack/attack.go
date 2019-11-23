package attack

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hmarf/trunks/trunks/report"
)

type Option struct {
	Requests    int
	Concurrency int
	URL         string
	Header      []Header
	OutputFile  string
}
type Header struct {
	Key   string
	Value string
}

// Request用
type Request struct {
	Client     *http.Client
	ResponseCH chan Response
}

// Response用
type Response struct {
	statusCode    int
	contextLength int64
	responseTime  time.Duration
}

func (r *Request) createRequest(o Option) *http.Request {
	req, err := http.NewRequest("GET", o.URL, nil)
	if err != nil {
		panic(err)
	}
	for _, h := range o.Header {
		req.Header.Set(h.Key, h.Value)
	}
	return req
}

func showDegreeProgression(t time.Duration, degree int, maxRequest float32) {
	progression := 50
	progressionCount := degree / (100 / progression)
	fmt.Printf("\r[%02d:%02d:%02d] [", int(t.Hours()), int(t.Minutes())%60, int(t.Seconds())%60)
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

func (rq *Request) Kikouha(wg *sync.WaitGroup, ch *chan int, req *http.Request) {
	defer wg.Done()
	rqStart := time.Now()
	resp, err := rq.Client.Do(req)
	if err != nil {
		<-*ch
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		<-*ch
		return
	}
	rq.ResponseCH <- Response{
		statusCode:    resp.StatusCode,
		contextLength: resp.ContentLength,
		responseTime:  time.Now().Sub(rqStart),
	}
	<-*ch
}

func (r *Request) Attack(o Option) time.Duration {

	r.Client = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          0, // DefaultTransport: 100, 0にすると無制限。
			MaxIdleConnsPerHost:   o.Requests,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 60 * time.Second,
	}
	r.ResponseCH = make(chan Response, o.Requests)

	// 並行処理するスレッド数を決める
	ch := make(chan int, o.Concurrency)

	req := r.createRequest(o)
	resp, err := r.Client.Do(req)
	if err != nil {
		fmt.Printf("\x1b[31m%v\x1b[0m\n", err)
		os.Exit(1)
	}
	resp.Body.Close()
	_, _ = ioutil.ReadAll(resp.Body)

	// Requestを投げる時間測定
	requestStart := time.Now()

	// stashはとりあえず0意外なら何でもいい
	stash := 10
	// ひたすらRequestを投げる
	wg := sync.WaitGroup{}
	for i := 0; i < o.Requests; i++ {
		ch <- 1
		wg.Add(1)
		degree := int((float32(i) / float32(o.Requests)) * 100)
		degreeP := degree / 5
		if degreeP != stash {
			stash = degreeP
			showDegreeProgression(time.Now().Sub(requestStart), degree, float32(o.Requests))
		}
		go r.Kikouha(&wg, &ch, req)
	}
	wg.Wait()
	totalTime := time.Now().Sub(requestStart)
	showDegreeProgression(totalTime, 100, float32(o.Requests))
	return totalTime
}

func (rq *Request) GetResults(totalTime time.Duration, requestCount int, channel int) report.ResultBenchMark {
	// メジャーなステータスコードを初期値とする
	// https://www.sakurasaku-labo.jp/blogs/status-code-basic-knowledgess
	_result := report.ResultBenchMark{}
	_result.StatusCode = map[int]int{
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

	_result.ConcurrencyLevel = channel
	_result.TotalRequests = requestCount
	_result.LatecyTotal = totalTime
	_result.LatecyMin = totalTime
	i := 0
LOOP:
	for ; ; i++ {
		select {
		case data := <-rq.ResponseCH:
			_result.LatecyAve += data.responseTime
			// 待機時間　max, min
			if data.responseTime > _result.LatecyMax {
				_result.LatecyMax = data.responseTime
			}
			if data.responseTime < _result.LatecyMin {
				_result.LatecyMin = data.responseTime
			}
			// Response の Status Code を数える
			v, ok := _result.StatusCode[data.statusCode]
			if ok {
				v++
				_result.StatusCode[data.statusCode] = v
			} else {
				_result.StatusCode[data.statusCode] = 1
			}
			// ContextLength
			_result.TotalDataReceived += data.contextLength
		default:
			break LOOP
		}
	}
	_result.SucceedRequests = i
	_result.FailedRequests = requestCount - i
	_result.LatecyAve /= time.Duration(requestCount)
	return _result
}
