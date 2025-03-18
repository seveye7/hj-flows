package main

import (
	"time"

	"hj-flows/flows"
)

func main() {
	// 注册流处理函数
	mgr := flows.NewStreamMgr(
	// flows.WithTopic("id_register", "register", registerStream),
	// flows.WithTopic("id_link_toutiao", "link_toutiao", linkToutiaoStream),
	)

	// 启动处理
	mgr.Start()

	defer mgr.Stop()

	time.Sleep(time.Second * 2)

	// send
	// mgr.GetStream("register").SendMessage(&Register{
	// 	Date:     "2024-01-01",
	// 	Uid:      1001,
	// 	CreateAt: time.Now().Format(time.DateTime),
	// 	DeviceId: "1234567890",
	// })

	time.Sleep(time.Hour)
}
