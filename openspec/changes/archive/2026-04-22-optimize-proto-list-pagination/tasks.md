## 1. Backend - API 分页支持

- [x] 1.1 修改 `ProtoTypesResponse` 结构体，添加分页字段（total, page, pageSize, totalPages）
- [x] 1.2 修改 `GetLoadedTypes` 返回已排序的类型列表
- [x] 1.3 在 `handleProtoTypes` 中解析 `page` 和 `pageSize` 查询参数
- [x] 1.4 实现分页逻辑：计算起止索引，返回对应页面数据

## 2. Frontend - 分页控件

- [x] 2.1 在 proto-debugger.html 类型选择区域添加分页控件 HTML 结构
- [x] 2.2 添加 CSS 样式使分页按钮可用/禁用状态可见
- [x] 2.3 实现前端分页状态管理（当前页、每页大小）
- [x] 2.4 修改 API 请求逻辑，传递分页参数并处理分页响应
- [x] 2.5 实现"上一页"/"下一页"按钮点击处理
- [x] 2.6 处理边界情况（第一页、最后一页按钮状态）