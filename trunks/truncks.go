package trunks

import (
	"fmt"

	"github.com/hmarf/trunks/trunks/attack"
)

func Trunks(o attack.Option) {
	// オラオラオラオラオラオラ！！！
	fmt.Println()
	request := attack.Request{}
	totalTime := request.Attack(o)

	// 結果表示
	resultBenchMark := request.GetResults(totalTime, o.Requests, o.Concurrency)
	resultBenchMark.ShowResult(o.OutputFile)
}
