## ADDED Requirements

### Requirement: HTTP 路由使用标准库
系统 SHALL 使用 Go 标准库 `net/http` 进行所有 HTTP 路由，替换第三方路由依赖。

#### Scenario: 上传 proto 文件
- **WHEN** 客户端发送 POST 请求到 `/api/proto/upload`，包含 proto 文件内容
- **THEN** 系统解析 proto 文件并返回已加载的类型名称

#### Scenario: 列出 proto 类型
- **WHEN** 客户端发送 GET 请求到 `/api/proto/types`
- **THEN** 系统返回所有已加载 proto 类型名称的 JSON 数组

#### Scenario: 删除 proto 类型
- **WHEN** 客户端发送 DELETE 请求到 `/api/proto/types/{name}`，其中 `{name}` 是路径参数
- **THEN** 系统移除该名称的 proto 类型并返回成功

#### Scenario: 解码二进制数据
- **WHEN** 客户端发送 POST 请求到 `/api/proto/decode`，包含 type、data 和 encoding 字段
- **THEN** 系统使用指定类型解码数据并返回 JSON 结果

#### Scenario: 静态文件服务
- **WHEN** 客户端请求任何不匹配 API 路由的路径
- **THEN** 系统从嵌入式文件系统或 `./` 目录提供相应的静态文件

#### Scenario: 根路径服务
- **WHEN** 客户端发送 GET 请求到 `/`
- **THEN** 系统直接提供 proto-debugger.html 作为首页

#### Scenario: 路径参数提取
- **WHEN** 处理带有路径参数的请求，如 `/api/proto/types/{name}`
- **THEN** 系统使用 `r.PathValue("name")` 提取参数值
