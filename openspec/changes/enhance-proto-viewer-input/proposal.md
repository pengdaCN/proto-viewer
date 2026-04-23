## Why

当前 Proto Viewer 的输入处理功能有限，只支持 hex 和 base64 两种固定格式，无法自动检测数据类型，也不支持直接解析 curl 命令格式。用户需要手动转换 curl 命令为可用格式，增加了使用复杂度。

## What Changes

1. **自动类型检测**：输入数据类型自动识别（hex、base64、curl-raw），无需手动选择
2. **curl-raw 输入支持**：新增 curl-raw 输入类型，自动提取 `--data-raw` 参数中的数据
3. **UI 清空功能**：
   - 新增"清空输入"按钮，单独清空输入框内容
   - 新增"全部清空"按钮，清空输入、类型选择、输出 JSON 和二进制数据
4. **curl 路径自动填充**：从 curl-raw 中自动提取 URL 路径的最后一段，填入 Proto 类型搜索框，并提供开关控制（默认开启）

## Capabilities

### New Capabilities

- `auto-type-detection`: 输入数据类型自动检测，根据内容特征判断是 hex、base64 还是 curl-raw
- `curl-raw-input`: 支持直接输入 curl 命令格式，自动提取 --data-raw 数据
- `ui-clear-buttons`: 清空输入和全部清空功能
- `curl-path-auto-fill`: 从 curl-raw 中提取路径填充到 Proto 类型搜索框

### Modified Capabilities

- 无

## Impact

- 前端：新增输入类型处理、UI 按钮组件
- 后端：新增 curl-raw 解析逻辑、auto-type-detection 逻辑