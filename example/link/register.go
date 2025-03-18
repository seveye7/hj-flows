package main

import (
	"os"

	"hj-flows/flows"
)

var f2 *os.File

func init() {
	f2, _ = os.OpenFile("register.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o666)
}

// registerStream register流处理
func registerStream(stream *flows.Stream) {
	stream = flows.Filter(stream, func(data *Register) bool {
		// 落地备份
		f2.WriteString(flows.StructToString(data))
		f2.Write([]byte("\n"))
		f2.Sync()

		// 1.根据设备查询曝光数据

		// 2.如果匹配则上报
		return true
	})

	_ = stream
}
