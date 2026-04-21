## Context

当前 `ProtoRegistry` 只支持单文件上传（`LoadProto`），每个文件独立编译。企业级 proto 项目通常包含多个互相引用的文件，需要完整目录结构和 import 解析。

## Goals / Non-Goals

**Goals:**
- 支持 tar 包上传包含多个 .proto 文件
- 自动分析 import 依赖关系
- 自动识别 root protos 进行编译
- 集成 Google Common Proto（从 Go module cache 获取）
- 每次加载后清空历史，只保留本次

**Non-Goals:**
- 不支持增量更新
- 不支持多用户隔离
- 不支持 Git URL 直接拉取

## Decisions

1. **tar 包只包含 .proto 文件**
   - 前端负责过滤，只打包 .proto 文件
   - 后端不需要判断文件类型

2. **依赖解析策略**
   - 解析每个 proto 文件的 import 语句
   - 构建有向图：file → imports
   - root protos = 没有其他文件 import 它的文件
   - 使用拓扑排序确定编译顺序

3. **Google Proto 获取**
   - 构建时从 Go module cache 复制 `google/protobuf/*.proto` 到 `assets/google-protobuf/`
   - 编译参数使用 `-I assets/google-protobuf`
   - 应用自带，无需运行时依赖 GOMODCACHE
   - 可以通过构建脚本自动化这个过程

4. **临时文件管理**
   - 解压到 `os.MkdirTemp("", "proto-dir-*")`
   - 函数结束后 `defer os.RemoveAll(tmpDir)`
   - 保持原有单文件上传的错误处理模式

5. **清空历史时机**
   - 每次成功调用 `LoadDirectory` 后，立即清空 `loadedDescs`
   - 确保只保留最近一次上传的 proto 定义

## Decisions

| 决策点 | 选择 | 理由 |
|--------|------|------|
| tar 格式 | tar（而非 zip） | Go 标准库原生支持 `archive/tar` |
| Google Proto | 从 GOMODCACHE 获取 | 零外部依赖，自动同步 |
| 依赖解析 | 静态分析 import 语句 | 无需运行 protoc 即可解析 |
| 编译策略 | 全量编译 root protos | 确保所有依赖都被包含 |

## Risks / Trade-offs

[风险] import 路径写法不一致
→ **缓解**: 解析 import 语句时规范化路径，识别相对路径（如 `../common/types.proto`）和标准路径（如 `google/protobuf/timestamp.proto`）

[风险] 循环依赖（理论上不应该，但可能出现）
→ **缓解**: 拓扑排序时检测环，若检测到则返回明确错误信息

[风险] tar 包过大（大型 proto 项目可能几十 MB）
→ **缓解**: 后端限制上传大小（如 50MB），前端可提示用户先检查

## Open Questions

无
