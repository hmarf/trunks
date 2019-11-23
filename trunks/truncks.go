package trunks

import "github.com/hmarf/trunks/trunks/attack"

func Trunks(Channel int, RequestCount int, url string, h []attack.Header, outFile string) {
	// オラオラオラオラオラオラ！！！
	request := attack.Request{}
	option := attack.Option{
		Concurrency: Channel,
		Requests:    RequestCount,
		URL:         url,
		Header:      h,
		OutputFile:  outFile}
	totalTime := request.Attack(option)

	// 結果表示
	resultBenchMark := request.GetResults(totalTime, RequestCount, Channel)
	resultBenchMark.ShowResult(outFile)
}
