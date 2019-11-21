package report

import (
	"fmt"
	"time"
)

// ベンチマークの結果を計算し収納する場所
type ResultBenchMark struct {
	ConcurrencyLevel  int           `json:"concurrency_level"`   //
	TotalRequests     int           `json:"total_requests"`      //
	SucceedRequests   int           `json:"succeed"`             // 通信に成功したRequest
	FailedRequests    int           `json:"failed"`              // 何らかの理由で通信に失敗したRequest
	RequestsSec       int           `json:"requests_sec"`        // 一秒間にアクセスできたRequestの総数
	TotalDataReceived int64         `json:"total_data_reveived"` // ContentLengthの総数
	StatusCode        map[int]int   `json:"status_code"`         // サーバーから返ってきたStatusCode
	LatecyTotal       time.Duration `json:"latecy_total"`        // 全てのResponseが返ってくるまでの総時間
	LatecyMax         time.Duration `json:"latecy_max"`          // Responseが来る待機時間の最も長かったもの
	LatecyMin         time.Duration `json:"latecy_min"`          // Responseが来る待機時間の最も短かったもの
	LatecyAve         time.Duration `json:"latecy_ave"`          // Responseが来る待機時間の平均
}

func (result *ResultBenchMark) ShowResult() {
	fmt.Printf("\n\nConcurrency Level:   %v\n", result.ConcurrencyLevel)
	fmt.Printf("Total Requests:      %v\n", result.TotalRequests)
	fmt.Printf("Succeeded requests:  %v\n", result.SucceedRequests)
	fmt.Printf("Failed requests:     %v\n", result.FailedRequests)

	if result.LatecyTotal < time.Duration(1*time.Second) {
		fmt.Printf("Requests/sec:        %d\n", result.SucceedRequests)
	} else {
		fmt.Printf("Requests/sec:        %d\n", result.RequestsSec)
	}
	fmt.Printf("Total data received: %v bytes\n", result.TotalDataReceived)
	fmt.Printf("\nStatus code:\n")
	for key, value := range result.StatusCode {
		if value != 0 {
			fmt.Printf("   [%v] %v responses\n", key, value)
		}
	}
	fmt.Printf("\nLatency:\n   total: %v\n   max:   %v\n   min:   %v\n   ave:   %v\n",
		result.LatecyTotal, result.LatecyMax, result.LatecyMin,
		result.LatecyAve,
	)
}
