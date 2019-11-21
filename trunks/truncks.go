package trunks

import "github.com/hmarf/trunks/trunks/attack"

func Trunks(Channel int, RequestCount int) {
	// オラオラオラオラオラオラ！！！
	request := attack.Request{}
	totalTime := request.Attack(Channel, RequestCount)

	// 結果表示
	resultBenchMark := request.GetResults(totalTime, RequestCount, Channel)
	resultBenchMark.ShowResult()
}
