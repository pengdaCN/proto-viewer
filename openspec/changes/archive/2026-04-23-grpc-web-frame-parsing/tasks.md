## 1. 实现 extractGrpcWebPayload 辅助函数

- [x] 1.1 在 `proto_loader.go` 中添加 `extractGrpcWebPayload(data []byte) ([]byte, error)` 函数
- [x] 1.2 实现循环帧解析逻辑：遍历数据，解析每个帧的 flags 和 length
- [x] 1.3 实现 DATA 帧过滤：flags bit 7 = 0 时收集 payload，bit 7 = 1 时跳过（trailer）
- [x] 1.4 实现压缩检测：flags bit 0 = 1 时返回错误 "不支持压缩的 gRPC-Web 数据"

## 2. 修改 Decode 函数

- [x] 2.1 在 `Decode` 函数中，当 `proto.Unmarshal` 失败时调用 `extractGrpcWebPayload`
- [x] 2.2 将 `extractGrpcWebPayload` 返回的数据用于第二次 `proto.Unmarshal` 尝试
- [x] 2.3 更新错误信息，准确反映尝试过的解析方式

## 3. 测试验证

- [x] 3.1 使用 `.testdata/1.hex` 测试数据验证 gRPC-Web 帧解析正确
- [x] 3.2 验证纯 protobuf 数据仍能正常解析（向后兼容）
- [x] 3.3 验证压缩数据返回正确的错误信息