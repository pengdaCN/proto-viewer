## 1. 后端基础设置

- [x] 1.1 添加 Go 依赖（protobuf 库、embed 支持）
- [x] 1.2 创建 Go embed 静态文件目录结构
- [x] 1.3 配置 HTTP 服务器和路由

## 2. Proto 动态加载核心

- [x] 2.1 添加 Go 依赖（protodesc、dynamicpb、protoregistry）
- [x] 2.2 实现 protoc 检查函数（检查系统是否有 protoc）
- [x] 2.3 实现 protoc 调用函数（生成 descriptor set）
- [x] 2.4 实现 proto 类型注册表（protoregistry + dynamicpb）
- [x] 2.5 实现动态加载 proto 类型函数
- [x] 2.6 实现动态卸载 proto 类型函数

## 3. 后端 API 实现

- [x] 3.1 实现 `POST /api/proto/upload` - 上传 .proto 文件
- [x] 3.2 实现 `DELETE /api/proto/types/:name` - 删除指定类型
- [x] 3.3 更新 `GET /api/proto/types` - 返回已加载的类型列表
- [x] 3.4 更新 `POST /api/proto/decode` - 使用 dynamicpb 反序列化
- [x] 3.5 添加错误处理：protoc 不存在
- [x] 3.6 添加错误处理：proto 文件语法错误
- [x] 3.7 添加错误处理：二进制解析失败

## 4. 前端界面

- [x] 4.1 更新 proto-debugger HTML 页面
- [x] 4.2 实现 .proto 文件上传组件
- [x] 4.3 实现已上传 proto 列表展示
- [x] 4.4 实现删除 proto 类型功能
- [x] 4.5 更新 proto 类型下拉选择框（实时刷新）
- [x] 4.6 保留二进制输入和解码功能
- [x] 4.7 保留 JSON 结果展示区域
- [x] 4.8 保留错误信息展示
- [x] 4.9 保留复制 JSON 到剪贴板功能

## 5. 集成与测试

- [x] 5.1 测试上传 .proto 文件并解析类型
- [x] 5.2 测试删除 proto 类型
- [x] 5.3 测试完整流程：上传 proto → 选择类型 → 输入 hex → 获取 JSON
- [x] 5.4 测试 base64 输入模式
- [x] 5.5 测试错误场景（无效数据、proto 语法错误、protoc 不存在）

**注意**: protoc 未安装时会显示警告提示，这是预期行为
