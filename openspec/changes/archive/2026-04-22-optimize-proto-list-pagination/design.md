## Context

当前 `/api/proto/types` 返回全部类型列表，无分页。当 proto 类型数量较多时（如数百个），前端加载缓慢，选择困难。列表顺序依赖于 protobuf 的内部迭代顺序，不稳定。

## Goals / Non-Goals

**Goals:**
- API 支持分页参数（page, pageSize），返回分页数据
- 类型列表按名称字母排序，保证稳定顺序
- 前端支持翻页操作

**Non-Goals:**
- 不支持服务端搜索/过滤（后续可扩展）
- 不修改现有非分页逻辑的兼容性（page 参数可选）

## Decisions

### 1. API 分页参数设计

**决定**: 使用 `page` 和 `pageSize` 作为查询参数

```go
GET /api/proto/types?page=1&pageSize=50
```

响应格式:
```json
{
  "types": ["Type1", "Type2", ...],
  "total": 150,
  "page": 1,
  "pageSize": 50,
  "totalPages": 3
}
```

**理由**:
- `page`/`pageSize` 是最常见的分页模式，直观易懂
- `pageSize` 默认值 50，平衡一次加载量和请求次数
- 不传参数时返回第一页，保持向后兼容

### 2. 排序策略

**决定**: 在 `GetLoadedTypes` 返回前进行字母排序

**理由**:
- 简单直接，无需修改底层数据结构
- 确保 API 返回顺序始终稳定

### 3. 前端分页控件

**决定**: 在 proto-debugger.html 类型选择区域添加分页控件

- 显示"第 X/Y 页，每页 Z 条"
- "上一页" / "下一页" 按钮
- 禁用状态下不可点击

## Risks / Trade-offs

- [风险] 大页数时前端请求频繁 → **缓解**: 默认 pageSize=50，减少请求次数
- [风险] 旧前端未传分页参数 → **缓解**: 参数可选，不传时返回第一页