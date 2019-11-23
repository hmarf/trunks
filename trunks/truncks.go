package trunks

import "github.com/hmarf/trunks/trunks/attack"

func Trunks(o attack.Option) { //Channel int, RequestCount int, url string, h []attack.Header, outFile string) {
	// オラオラオラオラオラオラ！！！
	request := attack.Request{}
	totalTime := request.Attack(o)

	// 結果表示
	resultBenchMark := request.GetResults(totalTime, o.Requests, o.Concurrency)
	resultBenchMark.ShowResult(o.OutputFile)
}
