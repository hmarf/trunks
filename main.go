package main

import "github.com/hmarf/go_benchmark/attack"

func main() {

	// 非同期数
	Channel := 10

	// Request数
	RequestCount := 10000

	// オラオラオラオラオラオラ！！！
	request := attack.Request{ResponseCH: make(chan attack.Response, RequestCount)}
	totalTime := request.Attack(Channel, RequestCount)

	// 結果表示
	resultBenchMark := request.GetResults(totalTime, RequestCount)
	resultBenchMark.ShowResult()
}
