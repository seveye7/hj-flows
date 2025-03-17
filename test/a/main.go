package main

import (
	"strconv"
	"time"

	"hj-flows/flows"
	"hj-flows/utils"
)

func main() {
	mgr := flows.NewStreamMgr(flows.WithKafka(&flows.KafkaConfig{}),
		flows.WithTopic("id_client", "client", clientStream),
	)

	mgr.Start()

	defer mgr.Stop()

	time.Sleep(time.Hour)
}

// clientStream client流处理
func clientStream(stream *flows.Stream) {
	stream = flows.Filter(stream, func(client *Client) bool {
		// 落地备份
		// save(client)
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
	}).SendToStream("user_daily_1202")
	// 数据1202UserDaily转发到[状态流]缓存实时数据
	flows.Map(stream, func(client *Client) *UserDaily {
		if client.Eventid != 2004 {
			return nil
		}

		dailyData := &UserDaily{
			Date: client.Date,
			Uid:  client.Uid,

			Days: utils.Ptr(client.RegisterDay),
		}

		adEcpm, err := strconv.ParseFloat(client.Ext5, 64)
		if client.Ext5 != "" && err != nil {
			return nil
		}
		if client.Ext2 == "1" {
			dailyData.AdVideoNum = utils.Ptr(uint32(1))
			dailyData.AdVideoAmount = utils.Ptr(adEcpm)
		}
		dailyData.AdAmountTotal = utils.Ptr(adEcpm)
		//
		return dailyData
	}).SendToStream("user_daily_2004")
}
