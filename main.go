package main

import (
	"github.com/hmarf/trunks/trunks"
)

func main() {

	// 非同期数
	Channel := 10

	// Request数
	RequestCount := 10000

	// オラオラオラオラオラオラ！！！
	trunks.Trunks(Channel, RequestCount)
}
