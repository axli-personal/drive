# 云端硬盘

提供Google Drive的基本功能, 并做一些扩展.

## 功能

* 文件和目录管理
* 分享
* 回收站
* 容量限制
* 搜索(TODO)
* 共享频道(TODO)
* 管理员接口(TODO)

## 部署

服务部署: Docker Compose.

```shell
$ ./build.sh # 构建脚本
```

## 服务

| 服务      | 功能        |
|---------|-----------|
| User    | 用户管理      |
| Drive   | 命名空间、容量管理 |
| Channel | 发送信息、发送文件 |
| Storage | 存储节点      |
| Mysql   | 数据库       |
| Redis   | 缓存、消息队列   |
| Web     | 客户端       |

## 项目结构

```text
/backend:
  /common:
    /auth: authentication
    /decorator: handler decorators
  /pkg:
    /events: shared event definition
    /types: shared request and response definition
    /errors: error definition
  /service:
    /adapters: interface implementation
    /domain: domain object
    /main: main function
    /ports: http and rpc handler
    /remote: remote service interface
    /repository: repository interface
    /service: application service
    /usecases: use case and event handler
/frontend: web client
/test: performance test
```

## 数据一致性

最终一致性: 保证消息不丢失, 至少被消费一次.

兜底措施: 增加定时任务(TODO).

## 可用性

云存储(TODO), 服务限流(TODO), 自动重启.

## 扩展性

云存储(TODO).

## 性能

压力测试, 索引(TODO), 数据库缓存(TODO), 文件缓存(TODO).
