## 1. 合并上传区域

- [x] 1.1 备份 `embed/proto-debugger.html` 为 `proto-debugger.html.bak`
- [x] 1.2 移除"上传 Proto 目录"区域 DOM（id="directoryUploadArea" 及其父级 panel-header）
- [x] 1.3 修改"上传 Proto 文件"区域文字为"点击或拖拽 .proto 文件或文件夹到此处"
- [x] 1.4 更新 fileInput 元素，添加 `webkitdirectory` 属性支持文件夹选择
- [x] 1.5 修改 `handleFileUpload` 函数，支持区分文件和文件夹并统一处理
- [x] 1.6 移除 `directoryInput` 及其事件监听器

## 2. 已加载类型列表固定高度

- [x] 2.1 修改 `.proto-list` CSS，添加 `max-height: 300px; overflow-y: auto;`
- [x] 2.2 移除 `paginationControls` 的内联 `display: flex` 冲突样式
- [x] 2.3 修改 `paginationState.pageSize` 默认值从 50 改为 10
- [x] 2.4 更新分页信息显示文字

## 3. 可搜索的 Proto 类型选择器

### 前端实现
- [x] 3.1 将 `<select id="protoType">` 替换为自定义 Combobox 结构（input + dropdown list）
- [x] 3.2 添加 Combobox CSS 样式（下拉框、选中高亮、输入框）
- [x] 3.3 实现下拉列表展开/收起逻辑
- [x] 3.4 实现键盘导航（上下键高亮，Enter 选中，Esc 关闭）
- [x] 3.5 实现 debounce 300ms 的搜索输入，发请求到后端
- [x] 3.6 保持 `change` 事件触发解码的兼容性

### 后端实现
- [x] 3.7 修改 `GET /api/proto/types` 接口，添加 `search` 查询参数
- [x] 3.8 后端实现模糊匹配搜索（LIKE %keyword%），忽略大小写
- [x] 3.9 搜索结果分页保持与原接口一致

## 4. 测试验证

- [x] 4.1 测试上传单个 .proto 文件
- [x] 4.2 测试上传多个 .proto 文件
- [x] 4.3 测试上传文件夹（包含 .proto 文件）
- [x] 4.4 测试拖拽上传文件和文件夹
- [x] 4.5 测试已加载类型列表分页
- [x] 4.6 测试 Proto 类型选择器搜索过滤（验证后端搜索 API）
- [x] 4.7 测试解码功能