## Context

Proto Debugger 页面当前包含两个独立的文件上传区域：
1. **上传 Proto 文件** - 使用 `<input type="file">` 支持多文件选择
2. **上传 Proto 目录** - 使用 `webkitdirectory` 属性支持文件夹选择

已加载的 Proto 类型列表框使用 `paginationState.pageSize = 50`，列表高度由内容撑开，不固定。

Proto 类型选择器使用标准 `<select>` 下拉框，不支持搜索。

## Goals / Non-Goals

**Goals:**
- 合并上传文件和上传文件夹为一个统一的 Upload 区域，支持拖拽上传单文件和文件夹
- 已加载类型列表固定高度，默认分页大小改为 10，与上传区域视觉一致
- Proto 类型选择器升级为可搜索的下拉框，支持模糊匹配（忽略大小写）

**Non-Goals:**
- 不改变后端 API 行为
- 不修改文件解析逻辑
- 不增加新的 API 端点

## Decisions

### 1. 合并上传区域
**选择**：将"上传 Proto 文件"和"上传 Proto 目录"两个独立区域合并为一个统一的上传区域。

**理由**：
- 减少用户认知负担，一个入口同时支持文件和文件夹
- 保持 UI 简洁性，避免重复的上传交互区域

**实现方式**：
- 统一的上传区域支持 drag & drop 文件和文件夹
- 点击触发文件选择器，同时允许选择文件和文件夹（`webkitdirectory`）
- 使用单一 `<input type="file">` 元素，通过 `dataTransfer.items` 区分文件和文件夹
- 移除原来的目录上传区域 DOM 和相关 JS 逻辑

### 2. 已加载类型列表固定高度
**选择**：设置 `proto-list` 容器固定高度，超出部分滚动，分页大小默认 10 条。

**理由**：
- 高度固定使页面布局更稳定，不会因为内容多少而跳动
- 分页大小 10 与上传区域视觉比例更协调
- 用户可控制的分页大小保持可配置

**实现方式**：
- `.proto-list` 设置 `max-height: 300px; overflow-y: auto;`
- 移除内联 `display: flex` 冲突样式
- `paginationState.pageSize` 默认值从 50 改为 10
- API 分页参数同步更新

### 3. 可搜索的 Proto 类型选择器
**选择**：在原有 `/api/proto/types` 接口上添加 `search` 查询参数，后端支持按类型名模糊搜索。

**理由**：
- 前端过滤无法利用后端的索引和分页能力，当类型非常多时效率低下
- 后端搜索可以更精准地控制搜索逻辑和结果排序
- 在原有接口上扩展，避免新增 API 端点

**实现方式**：
- 接口改为 `GET /api/proto/types?page=1&pageSize=10&search=keyword`
- `search` 参数可选，支持模糊匹配（后端实现 `LIKE %keyword%`）
- 搜索参数区分大小写或不敏感由后端决定
- 前端 Combobox 输入时触发搜索请求，debounce 300ms 防抖
- 下拉列表展示搜索结果，选中后触发解码

## Risks / Trade-offs

- **风险**：合并上传区域后，原有目录上传用户可能需要适应
  - **缓解**：统一区域同时支持文件和文件夹上传，功能不减
- **风险**：Combobox 实现复杂度
  - **缓解**：使用简单实现，不依赖外部库，直接基于现有选择器逻辑扩展

## Migration Plan

1. 备份 `proto-debugger.html`
2. 修改 HTML 结构：合并上传区域，移除独立目录上传区域
3. 修改 CSS：添加 `.proto-list` 固定高度样式，调整 Combobox 样式
4. 修改 JS：
   - 统一文件上传处理逻辑
   - 更新 `paginationState.pageSize` 默认值
   - 实现 Combobox 搜索功能
5. 测试上传文件、上传文件夹、解码功能
6. 回滚：如有问题，还原备份文件