## ADDED Requirements

### Requirement: 统一的文件上传入口
系统 SHALL 提供单一的上传区域，同时支持单文件和多文件（.proto）上传，以及文件夹（包含 .proto 文件）上传。

#### Scenario: 上传单文件
- **WHEN** 用户点击上传区域并选择单个 .proto 文件
- **THEN** 系统解析该文件并加载类型

#### Scenario: 上传多个文件
- **WHEN** 用户点击上传区域并选择多个 .proto 文件
- **THEN** 系统依次解析所有文件并加载类型

#### Scenario: 上传文件夹
- **WHEN** 用户点击上传区域并选择一个包含 .proto 文件的文件夹
- **THEN** 系统递归扫描文件夹，加载所有 .proto 文件并解析依赖

#### Scenario: 拖拽上传文件
- **WHEN** 用户拖拽 .proto 文件到上传区域
- **THEN** 系统高亮上传区域，松开后解析并加载类型

#### Scenario: 拖拽上传文件夹
- **WHEN** 用户拖拽包含 .proto 文件的文件夹到上传区域
- **THEN** 系统高亮上传区域，松开后递归加载所有 .proto 文件

### Requirement: 上传状态反馈
系统 SHALL 显示上传进度和结果状态。

#### Scenario: 上传成功
- **WHEN** 用户成功上传并加载 proto 文件后
- **THEN** 系统显示成功状态消息

#### Scenario: 上传失败
- **WHEN** 用户上传的文件无法解析
- **THEN** 系统显示错误消息并保持上传区域可用