package attack_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/hmarf/trunks/benche/attack"
)

func TestKikouha(t *testing.T) {

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
	r := attack.Request{}
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
	r.ResponseCH = make(chan attack.Response, RequestCount)
	// 並行処理するスレッド数を決める
	ch := make(chan int, Concurrency)
	wg := sync.WaitGroup{}
	r.Kikouha(&wg, &ch)
}
