## Why

当前 Proto Viewer 页面的上传区域和类型选择器存在可用性问题：上传区域分离为文件和文件夹两个入口增加了用户认知负担；已加载类型列表高度不固定导致页面布局不稳定；类型选择器仅支持下拉选择无法快速搜索。这些问题影响用户体验，需要统一优化。

## What Changes

- 合并上传文件与上传文件夹为统一的 Upload 区域，支持拖拽上传文件和文件夹
- 已加载的 Proto 类型列表框固定高度，默认分页大小调整为 10 条，与上传区域视觉一致
- Proto 类型选择器升级为可搜索的下拉框，支持模糊匹配（忽略大小写）

## Capabilities

### New Capabilities
- `ui-upload-consolidation`: 统一的文件上传入口，支持单文件和文件夹上传，统一 UI 和交互
- `ui-type-search`: 类型选择器支持模糊搜索，提升大量类型时的选择效率

### Modified Capabilities
- `directory-upload`: 当前支持文件夹上传（tar 包），UI 优化后需要确保交互一致性

## Impact

- 前端 UI 组件：上传区域、类型列表、类型选择器
- 可能影响 `/api/proto/upload` 和 `/api/proto/upload-directory` 的使用方式