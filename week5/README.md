# 同步转异步

作业源码地址：https://github.com/daidai53/webook/blob/week5/internal/service/sms/async/async.go

## 一、判断服务商崩溃的机制

1. 按一段时间的成功率统计：successRateCrash 中按周期统计请求和响应数，在判断是否崩溃时，取出请求和响应数计算成功率，若低于50%则判断崩溃，先返回成功并存redis，启动异步重试。如果本周期请求数小于10，则直接返回没有崩溃。
2. 自定义特殊错误码：若Send方法返回特殊错误码时，则立刻返回错误，认为是服务商的特定错误类型（比如手机号错误、模板错误之类的），直接返回失败并不纳入成功率的统计，该自定义错误码在NewAsyncSMSService时注入。

## 二、适合场景和优点

**适合的场景：**适合在运营商负载过高且能返回错误码时使用。

**优点：**默认统计成功率的周期是15秒，在负载过高时成功率快速下降，可以及时触发崩溃判断，给用户返回成功响应并能在一定时间后重试，用户只会感受到收到验证码的时间有一定延迟，对使用体验影响不大。自定义错误码可以排除运营商特定错误对崩溃的判断准确性，对于一些校验类的错误等可以和其他错误分开处理，防止误判。

**缺点：**如果异步重发全部失败，用户无感知，只能一直等待，或者用户在前端可以重发时手动重发。

**后续改进方案：**异步重发全部失败时，给前端通知，前端可以展示短信业务繁忙之类的信息。
