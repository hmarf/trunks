package attack

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

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

// ベンチマークの結果を計算し収納する場所
type ResultBenchMark struct {
	succeedRequests   int           // 通信に成功したRequest
	failedRequests    int           // 何らかの理由で通信に失敗したRequest
	requestsSec       int           // 一秒間にアクセスできたRequestの総数
	totalDataReceived int64         // ContentLengthの総数
	statusCode        map[int]int   // サーバーから返ってきたStatusCode
	latecyTotal       time.Duration // 全てのResponseが返ってくるまでの総時間
	latecyMax         time.Duration // Responseが来る待機時間の最も長かったもの
	latecyMin         time.Duration // Responseが来る待機時間の最も短かったもの
	latecyAve         time.Duration // Responseが来る待機時間の平均
}

func (result *ResultBenchMark) ShowResult() {
	fmt.Printf("\n\nSucceeded requests:  %v\n", result.succeedRequests)
	fmt.Printf("Failed requests:     %v\n", result.failedRequests)

	if result.latecyTotal < time.Duration(1*time.Second) {
		fmt.Printf("Requests/sec:        %d\n", result.succeedRequests)
	} else {
		fmt.Printf("Requests/sec:        %d\n", result.requestsSec)
	}
	fmt.Printf("Total data received: %v bytes\n", result.totalDataReceived)
	fmt.Printf("\nStatus code:\n")
	for key, value := range result.statusCode {
		if value != 0 {
			fmt.Printf("   [%v] %v responses\n", key, value)
		}
	}
	fmt.Printf("\nLatency:\n   total: %v\n   max:   %v\n   min:   %v\n   ave:   %v\n",
		result.latecyTotal, result.latecyMax, result.latecyMin,
		result.latecyAve,
	)
}

func ShowDegreeProgression(t time.Duration, degree int, maxRequest float32) {
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

func (rq *Request) Kikouha(wg *sync.WaitGroup, ch *chan int) {
	defer wg.Done()
	req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
	rqStart := time.Now()
	resp, err := rq.Client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
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

func (r *Request) Attack(c int, requestCount int) time.Duration {
	// MaxIdleConns: DefaultTransportでは100になっている。0にすると無制限
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 0
	// MaxIdleConnsPerHost: デフォルト値は2。0にするとデフォルト値が使われる
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = c

	// 並行処理するスレッド数を決める
	ch := make(chan int, c)

	// request := Request{}
	r.Client = &http.Client{}
	r.ResponseCH = make(chan Response, requestCount)

	// Requestを投げる時間測定
	requestStart := time.Now()

	// stashはとりあえず0意外なら何でもいい
	stash := 10
	// ひたすらRequestを投げる
	wg := sync.WaitGroup{}
	for i := 0; i < requestCount; i++ {
		ch <- 1
		wg.Add(1)
		degree := int((float32(i) / float32(requestCount)) * 100)
		degreeP := degree / 5
		if degreeP != stash {
			stash = degreeP
			ShowDegreeProgression(time.Now().Sub(requestStart), degree, float32(requestCount))
		}
		go r.Kikouha(&wg, &ch)
	}
	wg.Wait()
	totalTime := time.Now().Sub(requestStart)
	ShowDegreeProgression(totalTime, 100, float32(requestCount))
	return totalTime
}

func (rq *Request) GetResults(totalTime time.Duration, requestCount int) ResultBenchMark {
	// メジャーなステータスコードを初期値とする
	// https://www.sakurasaku-labo.jp/blogs/status-code-basic-knowledgess
	_result := ResultBenchMark{}
	_result.statusCode = map[int]int{
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

	_result.latecyTotal = totalTime
	_result.latecyMin = totalTime
	i := 0
LOOP:
	for ; ; i++ {
		select {
		case data := <-rq.ResponseCH:
			_result.latecyAve += data.responseTime
			// 待機時間　max, min
			if data.responseTime > _result.latecyMax {
				_result.latecyMax = data.responseTime
			}
			if data.responseTime < _result.latecyMin {
				_result.latecyMin = data.responseTime
			}
			// Response の Status Code を数える
			v, ok := _result.statusCode[data.statusCode]
			if ok {
				v++
				_result.statusCode[data.statusCode] = v
			} else {
				_result.statusCode[data.statusCode] = 1
			}
			// ContextLength
			_result.totalDataReceived += data.contextLength
		default:
			break LOOP
		}
	}
	_result.succeedRequests = i
	_result.failedRequests = requestCount - i
	return _result
}
