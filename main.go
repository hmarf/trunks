package main

import (
	"net"
	"net/http"
	"time"

	"github.com/hmarf/go_benchmark/attack"
)

func main() {

	// 非同期数
	Channel := 10

	// Request数
	RequestCount := 10000

	// オラオラオラオラオラオラ！！！
	request := attack.Request{
		Client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          0, // DefaultTransport: 100, 0にすると無制限。
				MaxIdleConnsPerHost:   RequestCount,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 60 * time.Second,
		},
		ResponseCH: make(chan attack.Response, RequestCount),
	}
	totalTime := request.Attack(Channel, RequestCount)

	// 結果表示
	resultBenchMark := request.GetResults(totalTime, RequestCount)
	resultBenchMark.ShowResult()
}
