## Why

当前 proto 类型列表一次性返回所有类型，当类型较多时会导致页面加载缓慢、选择困难。需要添加分页和排序功能来提升用户体验，保证列表顺序稳定。

## What Changes

- **修改 API**: `/api/proto/types` 支持 `page` 和 `pageSize` 查询参数，返回分页数据
- **稳定排序**: 类型列表按名称字母顺序排序，保证每次返回顺序一致
- **前端分页**: proto-debugger 页面添加分页控件，支持翻页操作
- **响应格式**: API 返回分页元数据（总数、当前页、每页数量、总页数）

## Capabilities

### New Capabilities
- `proto-list-pagination`: Proto 类型列表的分页和排序功能

### Modified Capabilities
- `proto-debugger`: 更新类型列表展示需求，支持前端分页交互

## Impact

- API 层: 修改 `handleProtoTypes` 和 `GetLoadedTypes` 方法
- 前端: 修改 proto-debugger.html 的类型选择组件