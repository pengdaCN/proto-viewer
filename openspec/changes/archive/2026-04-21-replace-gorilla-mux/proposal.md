## Why

gorilla/mux 是一个优秀的路由器库，但现在 Go 标准库的 `net/http` 已经原生支持了大多数 mux 的功能（路径参数 `r.PathValue`、路由匹配等）。移除第三方依赖可以减少供应链风险和维护成本。

## What Changes

- 移除 `github.com/gorilla/mux` 依赖
- 使用 `net/http` 标准库的 `http.ServeMux` 替代 `*mux.Router`
- 将路由注册从 ` mux.HandleFunc(path, handler)` 改为 `http.HandleFunc(path, handler)`
- 路径参数从 `mux.Vars(r)["name"]` 改为 `r.PathValue("name")`
- 静态文件服务改用 `http.FileServer` 和 `http.StripPrefix`
- 移除 `go.mod` 中的 gorilla/mux 依赖

## Capabilities

### New Capabilities
<!-- Capabilities being introduced. Replace <name> with kebab-case identifier (e.g., user-auth, data-export, api-rate-limiting). Each creates specs/<name>/spec.md -->

### Modified Capabilities
<!-- Existing capabilities whose REQUIREMENTS are changing (not just implementation). -->
- `http-routing`: 实现方式变更，从 gorilla/mux 改为标准库，但 API 路由行为保持不变

## Impact

- `main.go`: 路由代码重构
- `go.mod`: 移除 gorilla/mux 依赖
- 用户 API 不受影响（路径相同）
