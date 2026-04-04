# tao-gin

[![Go Report Card](https://goreportcard.com/badge/github.com/taouniverse/tao-gin)](https://goreportcard.com/report/github.com/taouniverse/tao-gin)
[![GoDoc](https://pkg.go.dev/badge/github.com/taouniverse/tao-gin?status.svg)](https://pkg.go.dev/github.com/taouniverse/tao-gin?tab=doc)

```
 _______             _____ _       
|__   __|           / ____(_)      
   | | __ _  ___   | |  __ _ _ __  
   | |/ _` |/ _ \  | | |_ | | '_ \ 
   | | (_| | (_) | | |__| | | | |
   |_|\__,_|\___/   \_____|_|_|_|
```

Tao Universe Gin Web 框架集成组件，提供 HTTP 服务能力，支持多实例工厂模式（可同时运行多个端口服务）。

## 安装

```bash
go get github.com/taouniverse/tao-gin
```

## 使用

### 导入

```go
import _ "github.com/taouniverse/tao-gin"
```

### 配置

单实例模式：

```yaml
gin:
  mode: debug              # gin mode: debug/release/test
  listen: 0.0.0.0         # listen address
  port: 8080              # listen port
  writer: zap             # logger writer name (must match a registered writer)
  trust_proxies:          # trusted proxy list
    - "192.168.1.0/24"
  html_pattern:           # HTML template glob pattern (optional)
  static_path:            # static file path (optional)
  pprof:                  # pprof configuration (optional)
    enable: false
    path: /debug/pprof
  run_after: []           # pre-task dependencies
```

多实例模式（同时监听多个端口）：

```yaml
gin:
  default_instance: public
  public:
    mode: release
    listen: 0.0.0.0
    port: 8080
    writer: zap
    trust_proxies:
      - "192.168.1.0/24"
    pprof:
      enable: true
      prefix: /pprof
  admin:
    mode: debug
    listen: 127.0.0.1
    port: 8081
    writer: zap
    pprof:
      enable: true
      prefix: /admin/pprof
  internal:
    mode: release
    listen: 10.0.0.0
    port: 8082
    writer: zap
```

### 获取 Engine

```go
// 获取默认实例
engine, err := gin.Engine()

// 获取指定名称的实例
publicEngine, err := gin.GetEngine("public")
adminEngine, err := gin.GetEngine("admin")

// 注册路由
engine.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "pong"})
})
```

## 工厂模式

| API | 说明 |
|-----|------|
| `gin.Factory` | `*tao.BaseFactory[*gin.Engine]` 工厂实例 |
| `gin.Engine()` | 获取默认实例的 `*gin.Engine`，返回 `(engine, error)` |
| `gin.GetEngine(name)` | 获取指定名称实例的 `*gin.Engine` |

## 配置项说明

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `schema` | `string` | `http` | 协议类型（http / https） |
| `host` | `string` | `localhost` | 主机名 |
| `port` | `int` | `8080` | 监听端口 |
| `listen` | `string` | `127.0.0.1` | 监听地址 |
| `mode` | `string` | `debug` | Gin 运行模式（debug / release / test） |
| `writer` | `string` | `tao` | 日志 writer 名称（需匹配已注册的 writer） |
| `trust_proxies` | `[]string` | `[]` | 可信代理列表（CIDR 格式） |
| `html_pattern` | `string` | `` | HTML 模板 glob 模式 |
| `static_path` | `string` | `` | 静态文件目录路径 |
| `pprof.enable` | `bool` | `true` | 是否启用 pprof |
| `pprof.prefix` | `string` | `/pprof` | pprof 路由前缀 |

## 单元测试

| 测试文件 | 说明 | 运行条件 |
|---------|------|---------|
| `config_test.go` | Config 解析与默认值验证 | 无需外部依赖 |
| `gin_test.go` | 完整集成测试（启动 Engine） | 无需外部服务 |

### 运行单元测试

```bash
# 全部测试可直接运行
go test -v ./...
```

## 依赖

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web 框架
