## 作业：为消息队列添加监控

## 提交文件

https://github.com/daidai53/webook/commit/81412bad526beec9ed484ef9b4f1075eab7c35ed

## 监控指标：

使用Counter，用来统计ConsumeClaim中各个处理点的关键信息：

1. 统计context的超时；用来观察系统的空闲状态，当有值时说明每秒处理消息量较少。
2. 统计处理消息的数量；
3. 统计反序列化失败的次数；如果有值，说明消费了非预期类型的消息，可能有代码错误。
4. 统计消息处理方法返回失败的次数，error作为label；统计消息处理过程中各种错误的数量。
5. 统计MarkMessage的次数。

## 告警设置

1. 在生产环境测试时，为反序列化失败次数设置告警，如果有值时，要排查代码问题。
2. 为各种error数量设置阈值告警，当某种error数量增多时，可能发生了问题。