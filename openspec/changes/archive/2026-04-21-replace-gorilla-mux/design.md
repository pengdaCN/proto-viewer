## Context

当前使用 `github.com/gorilla/mux` 作为 HTTP 路由器，版本维护依赖第三方库。Go 1.22+ 标准库已支持路径参数提取等常用功能。

## Goals / Non-Goals

**Goals:**
- 用 `net/http` 标准库替代 gorilla/mux
- 保持现有 API 路由行为不变
- 移除第三方依赖

**Non-Goals:**
- 不改变 API 路径（`/api/proto/upload`, `/api/proto/types/{name}` 等）
- 不改变 HTTP 处理逻辑（只改路由层）
- 不添加新功能

## Decisions

1. **使用 `http.ServeMux` 替代 `*mux.Router`**
   - Go 标准库内置，无需额外依赖
   - 支持路径参数: `r.PathValue("name")` (Go 1.22+)
   - `http.HandleFunc` 注册路由

2. **具体代码变更**

| 原代码 | 新代码 |
|--------|--------|
| `mux.NewRouter()` | `http.NewServeMux()` |
| ` mux.HandleFunc(path, h)` | `http.HandleFunc(path, h)` |
| `mux.Vars(r)["name"]` | `r.PathValue("name")` |
| `r.PathPrefix("/")` | `http.StripPrefix("/", http.FileServer(...))` |

3. **移除依赖**
   - 从 `go.mod` 中移除 `github.com/gorilla/mux`
   - 删除 `main.go` 中的 import

## Risks / Trade-offs

[风险] Go 1.22 以下版本不支持 `PathValue`
→ **缓解**: 要求 Go 1.22+，当前项目应已使用较新版本

[风险] `http.ServeMux` 不支持 `mux.Methods()` 链式调用
→ **缓解**: 改用 `r.Method == http.MethodDelete` 在处理函数内判断

[风险] 不支持路由中间件
→ **缓解**: 当前未使用中间件，无需迁移

## Open Questions

无
