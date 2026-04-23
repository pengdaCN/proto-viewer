## Context

gRPC-Web 响应使用特殊的帧格式封装 protobuf 数据。每个帧包含：
- 1 字节 flags（bit 7 = 1 表示 trailer，bit 0 = 1 表示压缩）
- 4 字节 length（大端序）
- N 字节 payload

现有代码在 `proto_loader.go:630-638` 仅尝试跳过固定的 5 字节，无法处理：
1. 多个 DATA 帧的情况
2. DATA 帧后紧跟 TRAILER 帧的情况

### gRPC-Web 帧格式详解

```
┌─────────────────────────────────────────────────────────────────────┐
│                        gRPC-Web Response                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌────────┬─────────────┬─────────────────────────┐                 │
│  │ Flags  │   Length    │        Payload          │  ← Frame N     │
│  │  1 B   │    4 B      │         N B             │                 │
│  └────────┴─────────────┴─────────────────────────┘                 │
│                                                                      │
│  Flags 字节格式:                                                      │
│    bit 7 (0x80): trailer 标志 - 1=trailer 帧, 0=data 帧             │
│    bit 6-1: 保留                                                     │
│    bit 0 (0x01): compressed 标志 - 1=压缩, 0=未压缩                  │
│                                                                      │
│  Length: 4 字节无符号整数，大端序 (big-endian)                       │
│                                                                      │
│  Payload: protobuf 数据（或 trailer 元数据）                          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 测试数据分析

使用 `.testdata/1.hex` 进行实际分析：

```
原始 hex:
00000000030a01308000000044687270632d7374617475733a20300d0a782d726571756573742d69643a2031396462356462622d663730302d343030302d383639342d6534306634373636643530300d0a
```

**解析结果：**

```
┌─────────────────────────────────────────────────────────────────────┐
│ Frame 1: DATA 帧                                                      │
│   Flags:     0x00 (0000 0000) - bit7=0 (data), bit0=0 (未压缩)        │
│   Length:   0x00000003 = 3 bytes                                    │
│   Payload:  0a 01 30                                                 │
├─────────────────────────────────────────────────────────────────────┤
│ Frame 2: TRAILER 帧                                                  │
│   Flags:     0x80 (1000 0000) - bit7=1 (trailer)                    │
│   Length:   0x00000044 = 68 bytes                                   │
│   Payload:  "grpc-status: 0\r\nx-request-id: 19db5dbb-f700-4000-..."│
└─────────────────────────────────────────────────────────────────────┘

字节布局:
  位置   0     1-4      5-7        8     9-12     13-80
        ┌─────┬────────┬─────────┐ ┌─────┬────────┬────────────────┐
        │ 00  │ 000003 │ 0a0130  │ │ 80  │ 000044 │ grpc-status... │
        └─────┴────────┴─────────┘ └─────┴────────┴────────────────┘
             ↑        ↑          ↑        ↑        ↑
           flags    length     payload  flags    length
```

## Goals / Non-Goals

**Goals:**
- 自动检测并正确解析 gRPC-Web 帧格式
- 从多个帧中提取正确的 protobuf 数据
- 保持向后兼容：纯 protobuf 数据（无帧头）仍可直接解析

**Non-Goals:**
- 不支持解压缩（compression）
- 不处理 chunked 编码的 trailer

## Decisions

### Decision: 循环帧解析而非固定偏移

**Options Considered:**
1. 循环遍历所有帧，解析每个 DATA 帧
2. 使用正则表达式查找第一个 `0a`（ protobuf 消息开始标记）

**Chosen:** 循环遍历帧（选项 1）

**Rationale:** 更健壮，不依赖 protobuf 内部细节（如 `0a` 是 Any 或子消息的开始）。符合 gRPC-Web 协议规范。

**Alternatives:**
- 正则匹配：脆弱，可能误匹配 payload 中的 0a

### Decision: 实现 `extractGrpcWebPayload` 辅助函数

在 `proto_loader.go` 中添加专门处理 gRPC-Web 帧的函数：

**Rationale:** 分离关注点，保持 Decode 函数简洁，便于测试。

### Decision: 先尝试直接解析，失败后再尝试 gRPC-Web 解析

**Rationale:** 保持向后兼容，纯 protobuf 数据（无帧头）仍然可以直接解析。gRPC-Web 解析作为 fallback。

## 算法详细设计

### `extractGrpcWebPayload` 函数

```go
// extractGrpcWebPayload 从 gRPC-Web 帧格式数据中提取 protobuf payload
// 返回提取的 protobuf 数据（可能是多个 DATA 帧的 payload 拼接）
//
// 帧格式:
//   - 1 字节 flags: bit7=1 表示 trailer, bit0=1 表示压缩
//   - 4 字节 length: 大端序无符号整数
//   - N 字节 payload
//
// 返回值:
//   - []byte: 提取的 protobuf 数据
//   - error: 错误信息（压缩数据、无效格式等）
func extractGrpcWebPayload(data []byte) ([]byte, error) {
    var result []byte
    offset := 0

    for offset < len(data) {
        // 至少需要 5 字节 (1 flags + 4 length)
        if offset + 5 > len(data) {
            break
        }

        flags := data[offset]
        length := binary.BigEndian.Uint32(data[offset+1 : offset+5])
        payloadStart := offset + 5
        payloadEnd := payloadStart + int(length)

        // 检查数据完整性
        if payloadEnd > len(data) {
            break
        }

        // 检查压缩标志 (bit 0)
        if flags&0x01 != 0 {
            return nil, errors.New("不支持压缩的 gRPC-Web 数据")
        }

        // 检查 trailer 标志 (bit 7)
        if flags&0x80 == 0 {
            // DATA 帧：收集 payload
            result = append(result, data[payloadStart:payloadEnd]...)
        }
        // else: TRAILER 帧：跳过

        // 移动到下一帧
        offset = payloadEnd
    }

    if len(result) == 0 {
        return nil, errors.New("未找到有效的 gRPC-Web DATA 帧")
    }

    return result, nil
}
```

### 算法流程图

```
                          ┌─────────────────┐
                          │     开始         │
                          └────────┬────────┘
                                   │
                                   ▼
                    ┌──────────────────────────────┐
                    │ offset = 0, result = []      │
                    └──────────────┬───────────────┘
                                   │
                                   ▼
                    ┌──────────────────────────────┐
              ┌────▶│  offset + 5 > len(data)?     │
              │     └──────────────┬───────────────┘
              │              Yes    │    No
              │                    │                 │
              │                    ▼                 ▼
              │         ┌─────────────────┐   ┌─────────────────────┐
              │         │  结束 (break)   │   │  读取 flags         │
              │         └─────────────────┘   │  读取 length        │
              │                                └──────────┬──────────┘
              │                                           │
              │                                           ▼
              │                         ┌─────────────────────────────────┐
              │                         │  payloadStart = offset + 5     │
              │                         │  payloadEnd = payloadStart +    │
              │                         │            int(length)          │
              │                         └──────────────┬────────────────┘
              │                                           │
              │                                           ▼
              │                         ┌─────────────────────────────────┐
              │                    ┌─────│  payloadEnd > len(data)?        │
              │                    │ Yes └──────────────┬──────────────────┘
              │                    │                   │ No
              │                    │                   │
              │                    │                   ▼
              │                    │     ┌─────────────────────────────────┐
              │                    │     │  flags & 0x01 != 0?              │
              │                    │     │  (检查压缩标志)                   │
              │                    │     └──────────────┬──────────────────┘
              │                    │            Yes    │    No
              │                    │                  │        │
              │                    │                  ▼        │
              │                    │     ┌─────────────────┐     │
              │                    │     │ 返回错误:       │     │
              │                    │     │ "不支持压缩数据" │     │
              │                    │     └─────────────────┘     │
              │                    │                            │
              │                    │                            ▼
              │                    │     ┌─────────────────────────────────┐
              │                    │     │  flags & 0x80 != 0?            │
              │                    │     │  (检查 trailer 标志)            │
              │                    │     └──────────────┬──────────────────┘
              │                    │            Yes    │    No
              │                    │                  │        │
              │                    │                  ▼        │
              │                    │     ┌─────────────┐       │
              │                    │     │ 跳过 (TRAILER│       │
              │                    │     │ 帧)         │       │
              │                    │     └──────┬──────┘       │
              │                    │            │              │
              │                    │            │   ┌──────────┴──────────┐
              │                    │            │   │  收集 payload 到      │
              │                    │            │   │  result              │
              │                    │            │   └──────────┬──────────┘
              │                    │            │              │
              │                    │            └──────┬──────┘
              │                    │                   │
              │                    │                   ▼
              │                    │     ┌─────────────────────────────────┐
              │                    │     │  offset = payloadEnd           │
              │                    │     │  继续循环                       │
              │                    │     └──────────────┬────────────────┘
              │                    │                    │
              │                    └────────────────────┤
              │                                     │
              │                   ┌──────────────────┘
              │                   │
              │                   ▼
              │     ┌─────────────────────────────────┐
              │     │  len(result) == 0?              │
              │     └──────────────┬──────────────────┘
              │              Yes    │    No
              │                   │        │
              │                   ▼        ▼
              │        ┌─────────────────┐   ┌─────────────────┐
              │        │ 返回错误:        │   │ 返回 result     │
              │        │ "未找到DATA帧" │   │  (提取的 payload)│
              │        └─────────────────┘   └─────────────────┘
              │                   │                  │
              └───────────────────┴──────────────────┘
                                   │
                                   ▼
                          ┌─────────────────┐
                          │      结束        │
                          └─────────────────┘
```

### 修改后的 `Decode` 函数

```go
func (p *ProtoRegistry) Decode(data []byte, typeName string) (string, error) {
    p.mu.RLock()
    defer p.mu.RUnlock()

    msgType, err := p.types.FindMessageByName(protoreflect.FullName(typeName))
    if err != nil {
        return "", fmt.Errorf("未找到类型 %s: %w", typeName, err)
    }

    dynMsg := dynamicpb.NewMessage(msgType.Descriptor())

    // 策略 1: 直接尝试解析
    if err := proto.Unmarshal(data, dynMsg); err == nil {
        // 解析成功
        return marshalToJSON(dynMsg)
    }

    // 策略 2: 尝试 gRPC-Web 帧解析
    grpcData, grpcErr := extractGrpcWebPayload(data)
    if grpcErr != nil {
        // 两种方式都失败，返回详细错误信息
        return "", fmt.Errorf("反序列化失败: %v (原始数据); %v (gRPC-Web帧)", err, grpcErr)
    }

    dynMsg = dynamicpb.NewMessage(msgType.Descriptor())
    if err := proto.Unmarshal(grpcData, dynMsg); err != nil {
        return "", fmt.Errorf("gRPC-Web帧解析后仍失败: %v", err)
    }

    return marshalToJSON(dynMsg)
}

func marshalToJSON(msg *dynamicpb.Message) (string, error) {
    jsonBytes, err := protojson.MarshalOptions{
        Indent: "  ",
    }.Marshal(msg)
    if err != nil {
        return "", fmt.Errorf("JSON 转换失败: %v", err)
    }
    return string(jsonBytes), nil
}
```

## 边界情况处理

| 情况 | 处理方式 |
|------|----------|
| 数据长度 < 5 字节 | 直接返回原始 Unmarshal 错误 |
| flags & 0x01 != 0 (压缩) | 返回 "不支持压缩的 gRPC-Web 数据" |
| 未找到任何 DATA 帧 | 返回 "未找到有效的 gRPC-Web DATA 帧" |
| payload 长度超出数据范围 | 跳过该帧，继续解析下一帧 |
| 多个 DATA 帧 | 拼接所有 DATA 帧的 payload |

## Risks / Trade-offs

[Risk] 错误解析非 gRPC-Web 数据
→ **Mitigation**: 先尝试直接 Unmarshal，失败后才使用 gRPC-Web 解析

[Risk] TRAILER 帧后还有额外数据
→ **Mitigation**: 循环解析完所有帧，收集所有 DATA 帧的 payload

[Risk] 数据被压缩
→ **Mitigation**: 暂不支持，返回错误提示压缩数据无法处理

[Risk] 空 DATA 帧 (length=0)
→ **Mitigation**: length=0 时，payload 为空，append 到 result 后继续
