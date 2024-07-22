# go-utils

Utils for golang

install:

```bash
go get -u github.com/onesaltedseafish/go-utils
```

The functions are listed below:

- Shell 包装好的方便调用的系统命令
  - 简单好用
  - 支持传递`context.Context`进行控制
- Reader 包装好的读取各种格式的文件
  - csv 文件
  - txt (以`\t`为分隔符)
- Writer 包装好的写各种格式的文件
  - csv 文件
- Simulate 模拟的实现
  - dhcp 模拟 DHCP 协议分布 IP 地址

## linters

Use [`golangci-lint`](https://golangci-lint.run/) to lint Code.

## Log Utils

install:

```bash
go get -u github.com/onesaltedseafish/go-utils/log
```

Log 开箱即用的 logger库，提供以下功能:

- 简单好用的 Console 日志输出
- 简单好用的 Json 日志文件输出
- 日志等级过滤
- 分布式 traceid 支持

