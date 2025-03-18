package main

import (
	"time"

	"hj-flows/flows"
)

func main() {
	// 注册流处理函数
	mgr := flows.NewStreamMgr(
		flows.WithTopic("id_client", "client", clientStream),
		flows.WithTopic("id_user_daily", "user_daily", userDailyStream),
	)

	// 启动处理
	mgr.Start()

	defer mgr.Stop()

	time.Sleep(time.Second * 2)

	// send
	mgr.GetStream("client").SendMessage(&Client{
		Date:        "2024-01-01",
		Uid:         1001,
		Eventid:     1202,
		RegisterDay: 1,
		CreateAt:    time.Now().Format(time.DateTime),
	})

	time.Sleep(time.Hour)
}
