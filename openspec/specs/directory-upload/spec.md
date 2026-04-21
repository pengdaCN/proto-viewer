## Purpose

支持目录级别的 proto 文件上传与解析能力。

## ADDED Requirements

### Requirement: 目录上传
系统 SHALL 支持通过 tar 包上传包含多个 .proto 文件的目录，并自动解析依赖关系。

#### Scenario: 成功上传目录
- **WHEN** 客户端发送 POST 请求到 `/api/proto/upload-directory`，包含 tar 包
- **THEN** 系统解压 tar 包，扫描所有 .proto 文件，解析 import 依赖，编译并加载

#### Scenario: tar 包为空
- **WHEN** 客户端上传空的 tar 包或不包含 .proto 文件的 tar 包
- **THEN** 系统返回错误：无可加载的 proto 文件

#### Scenario: 包含循环依赖
- **WHEN** tar 包中的 proto 文件存在循环 import（如 A import B, B import A）
- **THEN** 系统返回错误：检测到循环依赖

#### Scenario: Google Proto 依赖
- **WHEN** proto 文件 import 标准 Google Proto（如 `google/protobuf/timestamp.proto`）
- **THEN** 系统自动从 Go module cache 查找并包含 Google Proto 文件

#### Scenario: 覆盖历史
- **WHEN** 客户端成功上传并加载 proto 目录后，再次上传新的目录
- **THEN** 系统清空之前的 proto 定义，只保留本次上传的内容

### Requirement: 自动依赖解析
系统 SHALL 支持自动分析 proto 文件的 import 语句，确保 protoc 能正确编译所有 proto 文件。

#### Scenario: 识别 root protos
- **WHEN** tar 包包含多个互相引用的 proto 文件
- **THEN** 系统识别没有被其他 proto import 的文件作为 root protos（用于优化编译参数）

#### Scenario: 相对路径解析
- **WHEN** proto 文件使用相对路径 import 其他 proto（如 `../common/types.proto`）
- **THEN** 系统根据文件位置正确解析相对路径

#### Scenario: 标准路径解析
- **WHEN** proto 文件使用标准路径 import（如 `google/protobuf/any.proto`）
- **THEN** 系统正确识别为 Google Proto 并从 GOMODCACHE 获取

#### Scenario: 自动包含依赖
- **WHEN** 编译 proto 文件时
- **THEN** 系统将所有 proto 文件或 root protos 传给 protoc，protoc 自动解析 import 依赖

### Requirement: API 路由扩展
系统 SHALL 提供新的 API 路由以支持目录上传功能。

#### Scenario: 目录上传路由
- **WHEN** 客户端发送 POST 请求到 `/api/proto/upload-directory`
- **THEN** 系统解析 tar 包并返回与单文件上传相同的响应格式（类型列表）

#### Scenario: 路由方法正确
- **WHEN** 客户端使用 GET 或其他非 POST 方法访问 `/api/proto/upload-directory`
- **THEN** 系统返回 405 Method Not Allowed
