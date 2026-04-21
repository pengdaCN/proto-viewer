## Why

当前缺少一个便捷的 proto 数据查看工具。开发者在调试 proto 序列化的数据时，需要编写代码或使用命令行工具才能看到 JSON 格式的内容，不够直观高效。

## What Changes

- 新增 Web 界面，用户可直接输入序列化后的二进制数据并查看对应的 JSON 输出
- Go 后端提供 proto 类型列表查询接口
- Go 后端提供 proto 反序列化接口，支持将二进制数据转换为 JSON
- 默认输入模式为 hex 编码，支持切换为 base64

## Capabilities

### New Capabilities

- `proto-debugger`: Web 在线调试工具，支持选择 proto 类型、输入二进制数据并实时查看 JSON 输出

### Modified Capabilities

- (无)

## Impact

- 新增 Web 界面：proto 在线调试页面
- 新增后端支持：提供 proto 类型列表和反序列化能力（支持 Web 上传 .proto 文件动态加载）
