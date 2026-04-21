## 1. API 路由

- [x] 1.1 在 main.go 中添加 `POST /api/proto/upload-directory` 路由
- [x] 1.2 创建 `handleDirectoryUpload` 处理函数

## 2. ProtoRegistry 扩展

- [x] 2.1 实现 `LoadDirectory(reader io.Reader) ([]string, error)` 方法
- [x] 2.2 实现 tar 包解压到临时目录
- [x] 2.3 实现扫描目录下所有 .proto 文件
- [x] 2.4 实现解析 proto 文件的 import 语句，构建依赖图
- [x] 2.5 实现识别 root protos（用于优化编译参数）
- [x] 2.6 实现循环依赖检测（仅限本地 proto 文件）
- [x] 2.7 实现编译并加载 proto

## 3. Google Proto 打包

- [x] 3.1 创建 `assets/google-protobuf/` 目录
- [x] 3.2 创建 Go 构建脚本 `scripts/copy_google_proto.go`，在构建时自动复制文件
- [x] 3.3 构建脚本支持 Windows、Linux、macOS（使用 Go 原生代码，无 shell 依赖）
- [x] 3.4 修改 `LoadDirectory` 使用 `assets/google-protobuf/` 作为 include 路径

## 4. 清理逻辑

- [x] 4.1 在 `LoadDirectory` 成功后清空历史 loadedDescs

## 5. 前端页面

- [x] 5.1 添加"上传 Proto 目录"区域，支持选择文件夹
- [x] 5.2 前端过滤只打包 .proto 文件
- [x] 5.3 实现 createTarFromFiles 函数生成 tar 包

## 6. 测试

- [x] 6.1 运行 `go build` 确认编译通过
- [x] 6.2 启动服务，使用包含多文件的 tar 包测试上传功能
