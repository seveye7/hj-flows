package main

import (
	"os"

	"hj-flows/flows"
)

var f *os.File

func init() {
	f, _ = os.OpenFile("register.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o666)
}

// linkToutiaoStream 曝光流处理
func linkToutiaoStream(stream *flows.Stream) {
	stream = flows.Filter(stream, func(data *LinkToutiao) bool {
		// 落地备份
		f.WriteString(flows.StructToString(data))
		f.Write([]byte("\n"))
		f.Sync()

		// 保存曝光数据
		return true
	})

	_ = stream
}
