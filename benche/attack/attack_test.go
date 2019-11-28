package attack

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestKikouhaSuccess(t *testing.T) {

	// test server
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-Type", "text")
			fmt.Fprintf(w, "world")
			return
		},
	))
	defer ts.Close()

	RequestCount := 1000
	Concurrency := 2
	r := Request{}
	r.Client = &http.Client{
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
	}
	r.ResponseCH = make(chan Response, RequestCount)
	// 並行処理するスレッド数を決める
	ch := make(chan int, Concurrency)
	wg := sync.WaitGroup{}
	options := []Option{
		{
			Requests:    1,
			Concurrency: 1,
			URL:         ts.URL,
			Method:      "GET",
			Header:      []Header{},
			Body:        "",
			OutputFile:  "",
		},
		{
			Requests:    1,
			Concurrency: 1,
			URL:         ts.URL,
			Method:      "POST",
			Header:      []Header{},
			Body:        "",
			OutputFile:  "",
		},
	}
	fmt.Println("aa")
	for _, op := range options {
		ch <- 1
		wg.Add(1)
		req := r.createRequest(op)
		r.Kikouha(&wg, &ch, req)
	}
	wg.Wait()
}
