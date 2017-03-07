# Event dispatcher
美洽内部使用的消息订阅系统，与主流的消息订阅系统相比不是简单的提供了topic+channel的分发模式，而是使用： 
* King 
* Noble
* Knight
* Peasant
* CustomizedTags

4+n级结构，发送者可通过此4+n层结构做精细化的分类，同时对于消费者来说可以订阅任意层次结构的消息，如在一开始消息不多的情况下，消费者可以
订阅一个King来处理所有消息，随着业务的不断增长，消息数量越来越多, 单个消费者已经无法处理对应的消息这时可以通过订阅King+Noble的消息来将处理并行化。
这4层结构由用户自由定义，以下是一个使用举例：
* ProductID = King
* ModuleID = Noble
* TenantID = Knight
* CustomizedSharding = Peasant

Event dispatcher背后采用的是消息队列nsq,没有引用其他服务以尽最大可能得降低依赖。

## 如何使用
1. 启动nsqlookupd
```
nsqlookupd -broadcast-address=0.0.0.0
```
2. 启动nsqd并将lookup指向刚刚launch的lookupd
```
nsqd -lookupd-tcp-address="127.0.0.1:4160" -broadcast-address="127.0.0.1"
```
3. 使用repo中的producer example发送一条消息
```
go build git.meiqia.com/infrastructure/event-dispatcher/example/producer
./producer
```

4. 使用repo中的consumer就能订阅消息了
```
go build git.meiqia.com/infrastructure/event-dispatcher/example/consumer
./consumer
```

## 实现细节
producer本质是将4+n级结构生成对应的topic，生成的规则是：
```
topicPrefix = "x7x6y9-"

var buffer bytes.Buffer
buffer.WriteString(topicPrefix)
buffer.WriteString(fmt.Sprintf("%s-%s-%s-%s", e.King, e.Noble, e.Knight, e.Peasant))
for _, tag := range e.Tags {
    buffer.WriteString(fmt.Sprintf("-%s", tag))
}
```

consumer根据订阅层次情况来决定订阅哪些topic，如果对应的层次字符串为空，意味着这层已经下面的层次的所有消息全部订阅，例如：一个消费者订阅了king="abc",其余
的层次都没有定义，此时event-dispatcher将会遍历所有的topic，只要符合king="abc"的topic就会全部订阅。另外，consumer每隔一分钟去遍历一遍所有的topic,所以
如果有新的类型的消息生成，消费者将会在一分钟后才能消费到。




