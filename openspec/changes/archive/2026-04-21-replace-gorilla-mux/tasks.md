## 1. 移除 gorilla/mux 依赖

- [x] 1.1 从 main.go 中删除 `github.com/gorilla/mux` import
- [x] 1.2 从 go.mod 中移除 gorilla/mux 依赖

## 2. 重构路由代码

- [x] 2.1 将 `mux.NewRouter()` 替换为 `http.NewServeMux()`
- [x] 2.2 将 ` mux.HandleFunc` 替换为 `http.HandleFunc`
- [x] 2.3 将 `mux.Vars(r)["name"]` 替换为 `r.PathValue("name")`
- [x] 2.4 将 `r.PathPrefix("/").Handler(...)` 替换为 `http.StripPrefix("/", http.FileServer(...))`
- [x] 2.5 删除 DELETE 方法链式调用，改用 `r.Method == http.MethodDelete` 判断

## 3. 验证

- [x] 3.1 运行 `go build` 确认编译通过
- [x] 3.2 启动服务验证路由功能正常
