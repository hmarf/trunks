package attack

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hmarf/trunks/benche/report"
	"golang.org/x/net/http2"
)

type Option struct {
	Requests    int
	Concurrency int
	URL         string
	Method      string
	Header      []Header
	Body        string
	OutputFile  string
	Http2       bool
}
type Header struct {
	Key   string
	Value string
}

// Request用
type Request struct {
	Client          *http.Client
	ResponseSuccess chan Response
	ResponseFail    chan int
}

// Response用
type Response struct {
	statusCode    int
	contextLength int64
	responseTime  time.Duration
}

func (r *Request) createRequest(o Option) *http.Request {
	req, err := http.NewRequest(o.Method, o.URL, nil)
	if !(o.Body == "") {
		req.Body = ioutil.NopCloser(bytes.NewReader([]byte(o.Body)))
	}
	req.ContentLength = int64(len([]byte(o.Body)))
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
		rq.ResponseFail <- 1
		<-*ch
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		rq.ResponseFail <- 2
		<-*ch
		return
	}
	rq.ResponseSuccess <- Response{
		statusCode:    resp.StatusCode,
		contextLength: resp.ContentLength,
		responseTime:  time.Now().Sub(rqStart),
	}
	<-*ch
}

func (r *Request) Attack(o Option) time.Duration {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:          0, // DefaultTransport: 100, 0にすると無制限。
		MaxIdleConnsPerHost:   o.Requests,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if o.Http2 {
		if err := http2.ConfigureTransport(transport); err != nil {
			log.Panicf("Failed to configure h2 transport: %s", err)
		}
	}

	r.Client = &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}

	r.ResponseSuccess = make(chan Response, o.Requests)
	r.ResponseFail = make(chan int, o.Requests)

	// 並行処理するスレッド数を決める
	ch := make(chan int, o.Concurrency)

	// とりあえず一回RequestしてみてConnectできるかのテスト
	resp, err := r.Client.Do(r.createRequest(o))
	if err != nil {
		fmt.Printf("\x1b[31m%v\x1b[0m\n", err)
		os.Exit(1)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("\x1b[31m%v\x1b[0m\n", err)
		os.Exit(1)
	}
	resp.Body.Close()

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
		go r.Kikouha(&wg, &ch, r.createRequest(o))
	}
	wg.Wait()
	totalTime := time.Now().Sub(requestStart)
	r.Client.CloseIdleConnections()
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
	success := 0
	failed := 0
LOOP:
	for ; ; success++ {
		select {
		case data := <-rq.ResponseSuccess:
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
		case _ = <-rq.ResponseFail:
			failed++
		default:
			break LOOP
		}
	}
	_result.SucceedRequests = success
	_result.FailedRequests = failed
	_result.LatecyAve /= time.Duration(requestCount)
	return _result
}
