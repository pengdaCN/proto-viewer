## Why

当前系统只支持单文件 proto 上传，无法满足企业级大型 proto 项目需求。这类项目通常包含多个互相引用的 proto 文件，需要完整的目录结构和 import 依赖解析能力。

## What Changes

- 新增 `POST /api/proto/upload-directory` 接口，接受 tar 包
- tar 包只包含 `.proto` 文件，保持原有目录结构
- 自动扫描并分析 import 依赖关系
- 自动识别 root protos（不被其他 proto import 的文件）
- 支持 Google Common Proto（通过 Go module cache 自动获取）
- 每次加载成功后自动清空历史 proto，只保留本次上传

## Capabilities

### New Capabilities
- `directory-upload`: 目录级别 proto 上传与解析
  - 支持 tar 包上传
  - 自动依赖解析
  - Google Common Proto 集成

### Modified Capabilities
- `http-routing`: 新增 `/api/proto/upload-directory` 路由

## Impact

- `main.go`: 新增目录上传路由
- `proto_loader.go`: 新增 `LoadDirectory` 方法，支持 tar 解析和依赖分析
- `go.mod`: 无新增依赖（使用标准库 + Go module cache）
