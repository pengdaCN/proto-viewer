## Context

proto-viewer 是一个用于查看 proto 数据的工具。当前缺少在线调试功能，开发者需要编写代码或使用命令行工具才能将序列化后的二进制数据转换为 JSON 格式进行查看。

本设计旨在提供一个完整的 Web 界面 + Go 后端，实现输入二进制数据选择 proto 类型后直接输出 JSON 的功能。

## Goals / Non-Goals

**Goals:**
- 提供 Web 界面，用户可通过浏览器访问
- 支持在 Web 界面上传和管理 .proto 文件
- 支持输入 hex 或 base64 编码的二进制数据（默认 hex）
- 将二进制数据反序列化为 JSON 并返回给前端展示
- 支持动态添加、删除已加载的 proto 定义

**Non-Goals:**
- 不支持在线编辑 proto 定义文件内容（只上传）
- 不支持将 JSON 序列化回二进制（仅解码）
- 不支持复杂的 proto 类型（如 map、oneof 的完整展示）

## Decisions

### 1. Web UI 实现方式：Go 嵌入式前端

**决策**: 使用 Go embed 将静态文件（HTML/CSS/JS）嵌入到二进制中，由同一个 Go 服务提供前端页面和 API。

**理由**:
- 部署简单：单一二进制文件包含所有内容
- 无需额外的前端工程化流程
- 对于工具类应用足够高效

**替代方案**:
- 前后端分离：单独部署前端工程 → 增加了部署复杂度，不适合工具类应用
- 使用现有前端框架 → 引入额外依赖，本项目无需复杂 UI

### 2. Proto 定义管理：Web 界面上传 + 运行时动态加载

**决策**: 用户通过 Web 界面上传 .proto 文件，后端在收到上传时调用 protoc 生成 descriptor set 并加载。

**理由**:
- 用户无需接触服务器，直接通过浏览器管理 proto 定义
- 支持随时添加、删除 proto 定义
- 重启服务后 proto 定义需要重新上传（无持久化需求）

**API 设计**:
```
POST   /api/proto/upload     - 上传 .proto 文件，返回可用的类型列表
DELETE /api/proto/types/:name - 删除指定 proto 类型
GET    /api/proto/types       - 获取当前已加载的类型列表
POST   /api/proto/decode      - 反序列化二进制数据为 JSON
```

**流程**:
```
┌─────────────────────────────────────────────────────────────┐
│  用户上传 .proto 文件                                          │
│       │                                                      │
│       ↓                                                      │
│  保存到临时目录 / 内存                                         │
│       │                                                      │
│       ↓                                                      │
│  调用 protoc --descriptor_set_out 生成 descriptor set         │
│       │                                                      │
│       ↓                                                      │
│  使用 protodesc + dynamicpb 加载类型                          │
│       │                                                      │
│       ↓                                                      │
│  返回可用的类型列表给前端                                       │
└─────────────────────────────────────────────────────────────┘
```

**替代方案**:
- 持久化存储到数据库 → 增加复杂度，本工具无需持久化
- 启动时从文件夹加载 → 用户需登录服务器操作，不够便捷

### 3. 二进制输入格式：默认 hex，支持 base64

**决策**: 前端默认使用 hex 模式，提供切换按钮支持 base64。

**理由**:
- hex 是调试场景下最常用的格式，便于肉眼比对
- base64 用于某些特殊场景（如跨越剪切板传输）

### 4. 反序列化实现：dynamicpb + protoreflect

**决策**: 使用 `google.golang.org/protobuf/types/dynamicpb` 和 `google.golang.org/protobuf/reflect/protodesc` 实现动态反序列化。

**理由**:
- 无需预先编译 `.pb.go` 文件
- 支持嵌套 message、enum 等复杂类型
- 纯 Go 实现，无 CGO 依赖

**依赖库**:
- `google.golang.org/protobuf` - 核心库
- `google.golang.org/protobuf/reflect/protodesc` - 从 descriptor set 加载类型
- `google.golang.org/protobuf/types/dynamicpb` - 动态创建 message 实例
- `google.golang.org/protobuf/encoding/protojson` - JSON 序列化

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| 系统没有安装 protoc | 任何 API 调用时检查，如不存在则返回明确错误提示 |
| Proto 文件语法错误 | 上传时调用 protoc 验证，错误时返回编译错误信息 |
| Proto 定义文件过多 | 初期不做限制，后续可按需添加 |
| Hex 输入无效导致解析失败 | 后端返回明确错误信息，前端友好提示 |

## Open Questions

1. Proto 文件是否需要持久化存储？（重启后是否保留）
2. 前端 UI 的具体样式风格是否需要统一到项目现有风格？
