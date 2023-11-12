说明：
1. SafeOperate方法用来传入并发安全操作，期间对key和cntKey的操作上，都加了写锁。里面只能调用nolock的操作方法NLTTL、NLGET和NLSET。
2. 对过期数据的核查清除还没有实现。
3. SafeOperate方法内对并发方法的调用还没有增加保护。