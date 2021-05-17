# gobase

提供一个简单易用的Golang服务端脚手架，该项目尽可能的包含所有的模块，如果不需要某些模块可自行删除

## 服务端代码一般需要包含如下结构
- http 在线profile,metric,health check。所以这个模块是必须的
- grpc 微服务 
- mq consumer worker 异步处理任务
- cron 定时任务 (未支持)
- db gorm
- cache redis
- metric prometheus 
- trace jaeger
- logger zap

## 配置
- 读取环境变量ENV_CLUSTER
- 默认读取 config/dev.yaml

## 代码结构说明
```
| - cmd 程序入口
    | - project/main.go 程序入口
| - doc API文档
| - config 配置文件
| - internal 内部代码包  
    | - project 
        | - server 服务注册
        | - biz 业务层
        | - dao 转换层 
        | - model 数据层
        | - pkg 内部工具包
| - pkg 可对外的工具包
```

## 如何使用
依赖中间件：
- mysql
- redis
- kafka
- etcd
- 中间件实例部署 rely/docker-compose.yaml

## Golang 必读文章
- [Uber go 规范](https://github.com/xxjwxc/uber_go_guide_cn)
- [Effective Go中文版](https://www.kancloud.cn/kancloud/effective/72199)
- [Go-advice](https://github.com/cristaloleg/go-advice/blob/49798ebacb18acfc70f240bf8609a227f8ac2622/README_ZH.md)