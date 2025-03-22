# hj-flows
data flows engine...

## 介绍
hj-flows 是一个基于 Kafka 的流式数据处理引擎，支持多种数据处理方式.
现在基于hj-flows在公司内应用于数据和投放业务的数据处理.

## 特性
* 流之间通过消息队列进行通讯
* 部分流基于redis或者内存实现缓存
* 支持nsq和kafka
* 支持写入doris和duckdb

## 需求
1. 基于各种格式数据，在各自流处理函数中落地后，再转化为其他格式数据生成实时报表的基础表。[example/stream]


``` mermaid
flowchart TD
    A[客户端] -->|1.写入client埋点| B{client流}
    B -->|2.落地| C[client.log]
    B -->|3.生成userDaily数据| D{user_daily流}
    D -->|4.落地| E[userDaily.log]
```

2. 基于投放平台，生成投放实时报表的基础表。同时基于注册和广告行为触发向平台发送的激活等事件。[example/link]


``` mermaid
flowchart TD
    A[客户端] -->|1.点击广告| A0[头巨量平台条]
    A0 -->|2.触发曝光事件| B{link_toutiao流}
    B -->|3.落地| D1[link_toutiao.log]
    B -->|4.缓存| E2[曝光缓存]
    A[客户端] -->|5.点击登录| A1[游戏服务端]
    A1 -->|6.触发注册事件| C{register流}
    C -->|7.落地| D2[register.log]
    C -->|8.设备查询曝光| E2[曝光缓存]
    E2 -->|9.匹配上报激活| E[巨量平台]
```

3. 一定间隔时间，向平台请求数据，生成实时报表的计划每天每小时投放数据。[example/plan_data]
4. 每日凌晨，从数据文件抽取数据，清洗后生成T1报表数据。[example/report]



## 示例

``` go
// 初始化
mgr := flows.NewStreamMgr(flows.WithKafka(&flows.KafkaConfig{}),
    flows.WithTopic("id_client", "client", clientStream),
    flows.WithTopic("id_user_daily", "user_daily", userDailyStream),
)

// 启动处理
mgr.Start()

defer mgr.Stop()
```

