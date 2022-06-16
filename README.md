# geecache

geecache基本模仿[groupcache](github.com/golang/groupcache)的实现，是一个轻量级内嵌的分布式缓存系统。

该项目主要是为了Go语言熟悉以及进阶，**请勿用于生产环境**。

## 特性

- 单机缓存和基于 HTTP 的分布式缓存
- 最近最少访问(Least Recently Used, LRU) 缓存策略
- 使用 Go  singleflight机制防止缓存击穿
- 使用一致性哈希选择节点，实现负载均衡
- 使用 protobuf 优化节点间二进制通信

## 使用

```bash
go get -u github.com/nc-77/geecache
```

## 快速开始

详情见 [main.go](https://github.com/nc-77/geecache/blob/main/server/main.go) 和 [run.sh](https://github.com/nc-77/geecache/blob/main/run.sh)

## 感谢

- [gee-cache](https://github.com/geektutu/7days-golang/tree/master/gee-cache)
