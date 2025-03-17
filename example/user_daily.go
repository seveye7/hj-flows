package main

import (
	"os"

	"hj-flows/flows"
)

var f1 *os.File

func init() {
	f1, _ = os.OpenFile("userDaily.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o666)
}

// userDailyStream client流处理
func userDailyStream(stream *flows.Stream) {
	stream = flows.Filter(stream, func(client *UserDaily) bool {
		// 落地备份
		f1.WriteString(flows.StructToString(client))
		f1.Write([]byte("\n"))
		f1.Sync()
		return true
	})
}
