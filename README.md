# Event dispatcher
美洽内部使用的消息订阅系统，与主流的消息订阅系统相比不是简单的提供了topic+channel的分发模式，而是使用： 
* King 
* Noble
* Knight
* Peasant

4级结构，发送者可通过此4层结构做精细化的分类，同时对于消费者来说可以订阅任意层次结构的消息，如在一开始消息不多的情况下，消费者可以
订阅一个King来处理所有消息，随着业务的不断增长，消息数量越来越多, 单个消费者已经无法处理对应的消息这时可以通过订阅King+Noble的消息来将处理并行化。
这4层结构由用户自由定义，以下是一个使用举例：
* ProductID = King
* ModuleID = Noble
* TenantID = Knight
* CustomizedSharding = Peasant

Event dispatcher背后采用的是消息队列nsq,没有引用其他服务以尽最大可能得降低依赖。

生产者example:
```

```

消费者example:  
```

```



