package trunks

import "github.com/hmarf/trunks/trunks/attack"

func Trunks(Channel int, RequestCount int, u string, h []attack.Header, o string) {
	// オラオラオラオラオラオラ！！！
	request := attack.Request{URL: u, Header: h}
	totalTime := request.Attack(Channel, RequestCount)

	// 結果表示
	resultBenchMark := request.GetResults(totalTime, RequestCount, Channel)
	resultBenchMark.SaveFile = o
	resultBenchMark.ShowResult()
}
