## Why

当前 `/api/proto/decode` 接口在解析 Go grpc-web 返回的响应时会失败。gRPC-Web 响应包含多个帧（数据帧 + trailer 帧）拼接在一起，但现有代码仅尝试跳过固定的 5 字节，无法正确处理多个帧的情况。

## What Changes

- 修改 `ProtoRegistry.Decode` 函数，自动检测并正确解析 gRPC-Web 帧格式
- 循环遍历所有 gRPC 帧，只解析 DATA 帧（flags bit 7 = 0），跳过 TRAILER 帧
- 保持向后兼容：纯 protobuf 数据仍然可以直接解析

## Capabilities

### New Capabilities

- `grpc-web-decode`: gRPC-Web 帧自动检测与解析能力，自动去除帧头并定位正确的 protobuf 数据

### Modified Capabilities

- `proto-debugger`: 反序列化场景需要支持 gRPC-Web 格式数据（暂不修改 spec，因为这是后端实现细节）

## Impact

- 修改 `proto_loader.go` 中的 `Decode` 函数
- 可能需要添加 gRPC-Web 帧解析的辅助函数
