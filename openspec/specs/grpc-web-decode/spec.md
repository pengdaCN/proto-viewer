## Purpose

gRPC-Web 帧格式解析能力。当 `/api/proto/decode` 接收到的数据无法直接解析为 protobuf 时，系统 SHALL 自动尝试以 gRPC-Web 帧格式解析。

## ADDED Requirements

### Requirement: gRPC-Web 帧自动检测

系统 SHALL 在直接解析失败后，自动检测数据是否为 gRPC-Web 帧格式。

#### Scenario: 数据是 gRPC-Web 帧格式
- **WHEN** `/api/proto/decode` 收到 gRPC-Web 格式数据
- **THEN** 系统 SHALL 解析帧格式，跳过帧头和 trailer 帧
- **AND** 系统 SHALL 使用提取出的 protobuf 数据进行反序列化

#### Scenario: 数据是纯 protobuf 格式
- **WHEN** `/api/proto/decode` 收到纯 protobuf 数据（无帧头）
- **THEN** 系统 SHALL 直接解析，无需帧处理

#### Scenario: 数据是压缩的 gRPC-Web 帧
- **WHEN** `/api/proto/decode` 收到压缩的 gRPC-Web 数据
- **THEN** 系统 SHALL 返回错误 "不支持压缩的 gRPC-Web 数据"
