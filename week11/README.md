# 作业一

## 设计思路

在原有获取分布式锁的逻辑上，加上对负载的判断，认为低于30为低负载，高于80是高负载，其余是中负载。

1. 当本地存有分布式锁为空，尝试获取分布式锁时：
   + 若为低负载，则和原来实现一样，去申请获取分布式锁，申请成功后正常执行任务、自动续约、调度后续任务。
   + 若为高负载，则此次任务调度结束，不去获取分布式锁，本周期本节点不去参与任务执行。
   + 若为中负载，按“(load/100)*100%”的概率决定是否参与申请分布式锁，load越高，参与申请锁的概率越低。
2. 当调度任务并判断已经持有分布式锁时：
   + 若为低负载，则本次调度正常执行。
   + 若为高负载，则本地调度不进行任务执行，直接解锁分布式锁。
   + 若为中负载，按“(load/100)*100%”的概率决定是否立即释放分布式锁，load越高，释放锁的概率越高。

![image-20240120230947266](https://daidai-1300215655.cos.ap-shanghai.myqcloud.com/image/image-20240120230947266.png)

## 代码地址

https://github.com/daidai53/webook/blob/1c036f9eb2d19a7f2e272bc31248f8b65369a4e8/internal/job/ranking_job.go#L67

## 极端情况分析

+ 由于各节点在中低负载情况下，都有可能会参与分布式锁的抢夺，因此抢到锁的节点不一定是全局负载最低的节点。低负载节点参与抢夺的概率是100%，中负载节点参与抢夺的概率是根据负载而定的，最高也只有100%-30%=70%，因此分布式锁被低负载节点拿到的概率较大。如果中负载节点拿到分布式锁，且该节点负载确实不高时，影响不大，如果负载较高时，后续调度过程中就有很大概率将锁释放掉，让锁回到低负载节点手里。
+ 如果选中的节点宕机，和原有实现一样，无法进行续约，在租约过期后，锁会被其他正常的节点获取到并执行任务，影响不大。
+ 该设计基于原有课上用到的分布式锁框架，因此这样设计，如果自己重新设计考虑负载的分布式锁框架的话，可能会采用其他机制，例如本调度周期内，由一个节点作为中心，统计所有参与节点汇报上来的负载，从中决策出最低负载的节点。

# 作业二

代码地址：

https://github.com/daidai53/webook/blob/1c036f9eb2d19a7f2e272bc31248f8b65369a4e8/internal/repository/dao/job.go#L36
