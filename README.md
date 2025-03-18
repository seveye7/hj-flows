# hj-flows
data flows engine...

## 介绍
hj-flows 是一个基于 Kafka 的流式数据处理引擎，支持多种数据处理方式.
现在基于hj-flows在公司内应用于数据和投放业务的数据处理.

## 特性
* 流之间通过消息队列进行通讯
* 部分流基于redis或者内存实现缓存
* 支持nsq和kafka
* 支持写回doris和其他数据库

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

## examle/main.go

``` mermaid
flowchart TD
    A[客户端] -->|1.写入client埋点| B{client流}
    B -->|2.落地| C[client.log]
    B -->|3.生成userDaily数据| D{user_daily流}
    D -->|4.落地| E[userDaily.log]
```