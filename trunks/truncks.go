package trunks

import "github.com/hmarf/trunks/trunks/attack"

func Trunks(Channel int, RequestCount int, u string) {
	// オラオラオラオラオラオラ！！！
	request := attack.Request{URL: u}
	totalTime := request.Attack(Channel, RequestCount)

	// 結果表示
	resultBenchMark := request.GetResults(totalTime, RequestCount, Channel)
	resultBenchMark.ShowResult()
}
