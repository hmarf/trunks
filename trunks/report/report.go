package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// ベンチマークの結果を計算し収納する場所
type ResultBenchMark struct {
	SaveFile          string        // セーブするファイル名
	ConcurrencyLevel  int           `json:"concurrency_level"`   // 同時実行数
	TotalRequests     int           `json:"total_requests"`      // Requestを投げる総数
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

// ベンチマークの結果を計算し収納する場所
type ResultBenchMarkForJson struct {
	ConcurrencyLevel  int         `json:"concurrency_level"`   //
	TotalRequests     int         `json:"total_requests"`      //
	SucceedRequests   int         `json:"succeed"`             // 通信に成功したRequest
	FailedRequests    int         `json:"failed"`              // 何らかの理由で通信に失敗したRequest
	RequestsSec       int         `json:"requests_sec"`        // 一秒間にアクセスできたRequestの総数
	TotalDataReceived int64       `json:"total_data_reveived"` // ContentLengthの総数
	StatusCode        map[int]int `json:"status_code"`         // サーバーから返ってきたStatusCode
	LatecyTotal       string      `json:"latecy_total"`        // 全てのResponseが返ってくるまでの総時間
	LatecyMax         string      `json:"latecy_max"`          // Responseが来る待機時間の最も長かったもの
	LatecyMin         string      `json:"latecy_min"`          // Responseが来る待機時間の最も短かったもの
	LatecyAve         string      `json:"latecy_ave"`          // Responseが来る待機時間の平均
}

func (r *ResultBenchMark) ShowResult() {
	r.WriteResultFile()
	r.ShowResultConsole()
}

func convertStringFromDuration(r ResultBenchMark) ResultBenchMarkForJson {
	return ResultBenchMarkForJson{
		ConcurrencyLevel:  r.ConcurrencyLevel,
		TotalRequests:     r.TotalRequests,
		SucceedRequests:   r.SucceedRequests,
		FailedRequests:    r.FailedRequests,
		RequestsSec:       r.RequestsSec,
		TotalDataReceived: r.TotalDataReceived,
		StatusCode:        r.StatusCode,
		LatecyTotal:       fmt.Sprintf("%v", r.LatecyTotal),
		LatecyMax:         fmt.Sprintf("%v", r.LatecyMax),
		LatecyMin:         fmt.Sprintf("%v", r.LatecyMin),
		LatecyAve:         fmt.Sprintf("%v", r.LatecyAve),
	}
}

func (r *ResultBenchMark) WriteResultFile() {
	if r.SaveFile == "" {
		return
	}
	resultJson := convertStringFromDuration(*r)
	jsonBytes, err := json.Marshal(resultJson)
	if err != nil {
		log.Println("ベンチマークの結果の保存に失敗しました(json文字列の作成に失敗しました)")
		return
	}
	out := new(bytes.Buffer)
	json.Indent(out, jsonBytes, "", "    ")
	file, err := os.Create(r.SaveFile)
	if err != nil {
		log.Println("ベンチマークの結果の保存に失敗しました(fileを作成できませんでした。)")
		return
	}
	defer file.Close()

	file.Write(([]byte)(out.String()))
}

func (r *ResultBenchMark) ShowResultConsole() {
	fmt.Printf("\n\nConcurrency Level:   %v\n", r.ConcurrencyLevel)
	fmt.Printf("Total Requests:      %v\n", r.TotalRequests)
	fmt.Printf("Succeeded requests:  %v\n", r.SucceedRequests)
	fmt.Printf("Failed requests:     %v\n", r.FailedRequests)

	if r.LatecyTotal < time.Duration(1*time.Second) {
		fmt.Printf("Requests/sec:        %d\n", r.SucceedRequests)
	} else {
		fmt.Printf("Requests/sec:        %d\n", r.RequestsSec)
	}
	fmt.Printf("Total data received: %v bytes\n", r.TotalDataReceived)
	fmt.Printf("\nStatus code:\n")
	for key, value := range r.StatusCode {
		if value != 0 {
			fmt.Printf("   [%v] %v responses\n", key, value)
		}
	}
	fmt.Printf("\nLatency:\n   total: %v\n   max:   %v\n   min:   %v\n   ave:   %v\n",
		r.LatecyTotal, r.LatecyMax, r.LatecyMin,
		r.LatecyAve,
	)
}
