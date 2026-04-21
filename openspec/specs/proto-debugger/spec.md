## Purpose

Proto Debugger 前端页面功能规范。

## ADDED Requirements

### Requirement: Proto 类型列表展示

系统 SHALL 在页面加载时从后端获取可用的 proto 类型列表，并展示在下拉选择框中供用户选择。

#### Scenario: 页面加载时获取类型列表
- **WHEN** 用户打开 proto-debugger 页面
- **THEN** 系统 SHALL 从 `/api/proto/types` 获取类型列表
- **AND** 系统 SHALL 将类型列表填充到下拉选择框中

#### Scenario: 类型列表为空时
- **WHEN** 后端返回空的类型列表
- **THEN** 系统 SHALL 在页面显示提示信息"无可用的 proto 类型，请检查 proto 定义文件夹配置"

### Requirement: 二进制数据输入

系统 SHALL 允许用户输入二进制数据，支持 hex 和 base64 两种编码格式，默认使用 hex 格式。

#### Scenario: 用户输入 hex 格式数据
- **WHEN** 用户在输入框中输入 hex 编码的二进制数据
- **AND** 用户选择 hex 输入模式
- **THEN** 系统 SHALL 保持输入内容不变

#### Scenario: 用户切换到 base64 模式
- **WHEN** 用户点击切换按钮选择 base64 模式
- **THEN** 系统 SHALL 在输入框上方显示 "base64" 标签
- **AND** 系统 SHALL 将后续输入作为 base64 格式处理

### Requirement: Proto 反序列化

系统 SHALL 将用户输入的二进制数据根据选定的 proto 类型反序列化为 JSON 格式并展示。

#### Scenario: 成功反序列化
- **WHEN** 用户输入有效的二进制数据
- **AND** 用户选择正确的 proto 类型
- **AND** 用户点击"解码"按钮
- **THEN** 系统 SHALL 调用 `/api/proto/decode` 接口
- **AND** 系统 SHALL 在页面展示返回的 JSON 结果

#### Scenario: 反序列化失败 - 无效二进制
- **WHEN** 用户输入无效的 hex 数据（如包含非hex字符）
- **AND** 用户点击"解码"按钮
- **THEN** 系统 SHALL 在页面显示错误信息"无效的二进制数据格式"

#### Scenario: 反序列化失败 - 类型不匹配
- **WHEN** 用户输入有效的二进制数据
- **AND** 用户选择的 proto 类型与数据不匹配
- **AND** 用户点击"解码"按钮
- **THEN** 系统 SHALL 在页面显示错误信息"反序列化失败：数据与选定的 proto 类型不匹配"

#### Scenario: 反序列化失败 - 类型不存在
- **WHEN** 用户未选择任何 proto 类型
- **AND** 用户点击"解码"按钮
- **THEN** 系统 SHALL 在页面显示错误信息"请选择 proto 类型"

### Requirement: JSON 结果展示

系统 SHALL 友好地展示反序列化后的 JSON 结果。

#### Scenario: 展示格式化后的 JSON
- **WHEN** 反序列化成功返回 JSON 数据
- **THEN** 系统 SHALL 使用格式化的方式展示 JSON（带缩进和换行）
- **AND** 系统 SHALL 支持复制 JSON 到剪贴板的功能

#### Scenario: JSON 中包含嵌套消息
- **WHEN** 反序列化后的 JSON 包含嵌套的 message
- **THEN** 系统 SHALL 完整展示嵌套结构
- **AND** 系统 SHALL 支持折叠/展开嵌套节点（如适用）
