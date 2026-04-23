## Context

Proto Viewer 当前支持 hex 和 base64 两种固定格式输入，需手动选择类型。用户希望直接粘贴 curl 命令（curl-raw）进行解码，并希望输入数据类型能够自动检测。

当前 backend 实现（main.go）：
- `DecodeRequest` 包含 `data`, `type`, `encoding` 三个字段
- 支持 `hex` 和 `base64` 两种编码
- `handleProtoDecode` 根据 encoding 选择对应的解码方式

grpc-web 协议数据格式（来自 .testdata/curl-raw.txt）：
- 数据以 5 字节的 grpc-web 帧头开始（`\x00\x00\x00\x00\x00`）
- 帧头后是实际 protobuf 数据

## Goals / Non-Goals

**Goals:**
- 自动检测输入数据类型（hex vs base64 vs curl-raw），减少用户操作
- 支持 curl-raw 输入格式，自动提取 `--data-raw` 数据
- 提取 curl URL 路径最后一段作为 Proto 类型搜索关键词
- 提供 UI 清空功能（清空输入 / 全部清空）

**Non-Goals:**
- 不修改现有的 proto 解析和加载逻辑
- 不改变 DecodeResponse 的 JSON 格式

## Decisions

### 1. 数据类型自动检测逻辑

**方案**: 根据输入内容特征判断类型

```
if input starts with "curl ":
    return "curl-raw"
else if input is valid hex:
    return "hex"
else if input is valid base64:
    return "base64"
else:
    return error
```

**特征判断**:
- curl-raw: 以 `curl ` 开头
- hex: 仅包含 `0-9a-fA-F`，且长度为偶数
- base64: 符合 base64 字符集（`A-Za-z0-9+/=`），长度是 4 的倍数

**备选方案**: 尝试解码法（先用 hex 解码，失败则用 base64）- 风险较高，因为 hex 字符串可能是有效的 base64 数据

### 2. curl-raw 解析逻辑

**解析步骤**:
1. 检测输入是否以 `curl ` 开头
2. 提取 `--data-raw` 参数值（支持单引号、双引号、无引号三种格式）
3. 处理 `\u0000` 等 Unicode 转义序列（grpc-web 帧头为 5 字节 `\x00`）
4. 对提取的数据进行解码（自动检测 hex/base64 或使用默认 grpc-web 二进制）

**正则表达式**:
```
--data-raw\s+(?:'([^']*)'|"([^"]*)"|(\S+))
```

### 3. curl 路径提取逻辑

从 curl URL 中提取路径最后一段：
```
URL: https://console.mpcvault.com/admin/mpcvault.internal.admin.v1.AdminService/CheckCredential
提取结果: CheckCredential
```

### 4. 后端 API 变更

后端保持现有 `hex` 和 `base64` 两种 encoding 类型，不做变更。

curl-raw 解析和自动检测全部在前端完成：
- 前端识别到 curl-raw 后，提取 `--data-raw` 数据
- 处理 Unicode 转义序列（`\u0000` → byte）
- 转为 base64 编码后传递给后端
- 后端收到 `encoding: "base64"`，直接解码

### 5. grpc-web 二进制数据处理

前端提取 curl-raw 数据后：
1. 处理 `\uXXXX` Unicode 转义序列，还原为原始字节
2. 跳过 grpc-web 帧头（前 5 字节为 `\x00\x00\x00\x00\x00` 或类似）
3. 将剩余数据转为 base64 传递给后端

**注意**：后端不感知 curl-raw，所有数据经过前端处理后以 base64 形式发送。

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| auto-detection 无法判断时 | 回退到 hex 或返回明确错误 |
| curl 命令格式不标准 | 支持常见变体（--data, --data-raw, -d） |
| 路径提取失败 | 提供 fallback，不阻塞主流程 |

## API 变更详情

### POST /api/proto/decode

**Request Body**:
```json
{
    "data": "string",
    "type": "string",
    "encoding": "hex|base64"
}
```

**Encoding 说明**:
- `hex`: 输入是十六进制字符串
- `base64`: 输入是 base64 编码字符串

**Response**:
```json
{
    "json": "string"
}
```

### 前端变更

1. 类型选择器保持现有的 hex/base64 两个选项
2. 输入框新增 "清空输入" 和 "全部清空" 按钮
3. 前端实现自动检测逻辑：
   - 以 `curl ` 开头 → 识别为 curl-raw，提取数据并跳过帧头，转为 base64
   - 仅含 `0-9a-fA-F` 且偶数位 → hex
   - 符合 base64 字符集 → base64
4. 自动检测失败时，显示提示信息（小气泡）
5. curl-raw 输入时，新增勾选框控制是否自动填充 Proto 类型搜索框
6. curl-raw 解析成功后，自动填充搜索框（如果勾选），同时将 encoding 设为 "base64"

## Open Questions

1. curl-raw 数据是否一定包含 grpc-web 帧头？是否需要处理无帧头情况？
2. 自动检测失败时的提示方案待定（可考虑小气泡提示）