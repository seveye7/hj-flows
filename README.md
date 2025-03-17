# hj-flows
data flows engine...

## 功能需求

### 目的是基于公司内部flows库开发

### 实现对nsq，kafka的支持

### 对流数据的常用函数支持。参考flink的函数库
https://nightlies.apache.org/flink/flink-docs-release-1.20/zh/docs/dev/datastream/operators/overview/#window-join

1. Map 输入一个元素同时输出一个元素
``` go
```

2. FlatMap 输入一个元素同时产生零个、一个或多个元素
``` go
```

3. Filter 为每个元素执行一个布尔 function，并保留那些 function 输出值为 true 的元素。
``` go
```
4. KeyBy 按照指定的键对元素进行分组。
在逻辑上将流划分为不相交的分区。具有相同 key 的记录都分配到同一个分区。在内部， keyBy() 是通过哈希分区实现的。有多种指定 key 的方式。
5. Reduce 对元素应用 reduce function，生成一个包含更新后状态的单一元素的流。
6. Aggregate 对元素应用 aggregate function，生成一个包含更新状态的流。
7. Window 对流元素进行分组，并将每个分组的元素分配到一个 window 中。
8. FlatMap 对元素应用 flatMap function，生成一个包含零个、一个或多个元素的流。


### 实现工作流支持
