# gobase

提供一个简单易用的Golang服务端脚手架，该项目尽可能的包含所有的模块，如果不需要某些模块可自行删除

## 服务端代码一般需要包含如下结构
- http 在线profile,metric,health check。所以这个模块是必须的
- grpc (未支持)
- mq consumer worker
- cron (未支持)
- db gorm
- cache redis
- metric prometheus (未支持)
- trace jaeger (未支持)
- logger zap


## 如何使用
依赖中间件：
- mysql
- redis
- kafka
- etcd
依赖中间件简单部署：
  