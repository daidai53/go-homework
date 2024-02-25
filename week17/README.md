# 一、支付消息兜底机制

代码地址（Commit）：https://github.com/daidai53/webook/commit/3a28145775293420a7a2b3c3b0544bd74545a6ed

## 思路

1. 在数据库中新增表，记录PaymentEvent，包含Id，PaymentEvent所需字段，是否发送等。
2. 在NativePaymentService更新支付状态的过程中，发送PaymentEvent到Kafka之前，先记录PaymentEvent到数据库并记录为未发送，发送到Kafka成功之后，再更新数据库记录为已发送。
3. 新增定时任务，从PaymentEvent表中批量扫描状态为未发送的记录并发送到Kafka，发送成功则记录为已发送，否则等待下次扫到再次尝试。

# 二、记账的幂等性

代码地址：https://github.com/daidai53/webook/commit/71964a8f2ca9b5e43443669459b2d7270de9b979

## 思路

在AccountRepository中新增存储Account缓存的接口，将记账记录缓存到Redis，key由biz和bizId构造而成，AccountService在记账前会先尝试添加缓存，如果缓存返回错误，说明该记录已存在，判断为重复记账，不再进行记账处理。

Redis的该记录设置3天的过期事件，保证该缓存在微信支付最大可能的回调时间后过期。

# 三、服务治理措施

可观测性、熔断限流降低等保护。

# 四、三者对账

## 思路

Reward服务定期取一批Reward记录，组装成消息向Payment和Account服务发起对账，由于Reward服务是整个打赏流程调用链的第一个服务，因此认为Reward上的记录是完整的。

Payment服务收到记录后，检查每一个交易号是否存在，不存在则记录并向Payment服务返回结果。如果该Payment记录结果为Payed，则Payment服务向Account服务申请检查是否有对应记账记录，如果没有记账记录或者缺失分成或客户账户的记录，则向Reward服务返回结果。Reward汇总后调用Payment服务和Account服务补齐记录并导出核查的不一致结果进行归档。

Account服务每天定时发起反向核查，检查账单在Payment中是否有记录，如果有的话是否在Reward中有记录，如果没有则通知Payment和Reward补齐记录，由于信息丢失，Payment和Reward上补齐的记录标记为异常并导出核查结果，供人工校对。