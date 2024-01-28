# 第一次作业

拆分Interactive的repository为单独的微服务，Interactive的service通过grpc调用repository中的方法。web微服务中将本地流量调为100%，InteractiveService全部走本地调用，然后通过grpc与InteractiveRepository微服务进行交互。

代码变更地址：https://github.com/daidai53/webook/commit/08d866615492fc17b554d4517550eeacfc27392e

# 选做作业

拆分code service为微服务。

代码变更：

1. 增加测试覆盖：https://github.com/daidai53/webook/commit/039fde72e18f0698b1c0021a98fb0434132c5f6d
2. 微服务拆分：https://github.com/daidai53/webook/commit/c791ec81cf52e7727f5ffc89fd01670d2b452edc