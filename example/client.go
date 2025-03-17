package main

import (
	"os"

	"hj-flows/flows"
	"hj-flows/utils"
)

var f *os.File

func init() {
	f, _ = os.OpenFile("client.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o666)
}

// clientStream client流处理
func clientStream(stream *flows.Stream) {
	stream = flows.Filter(stream, func(client *Client) bool {
		// 落地备份
		f.WriteString(flows.StructToString(client))
		f.Write([]byte("\n"))
		f.Sync()
		return true
	})
	// 数据转换 1202
	flows.Map(stream, func(client *Client) *UserDaily {
		if client.Eventid != 1202 {
			return nil
		}

		dailyData := &UserDaily{
			Date: client.Date,
			Uid:  client.Uid,

			Days:    utils.Ptr(client.RegisterDay),
			GameNum: utils.Ptr(uint32(1)),
		}

		return dailyData
	}).SendToStream("user_daily")
}
