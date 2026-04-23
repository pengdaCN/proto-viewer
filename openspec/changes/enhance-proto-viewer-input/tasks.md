## 1. Frontend - 自动类型检测

- [x] 1.1 在前端实现 `detectEncoding(input string)` 函数，通过特征判断数据类型
- [x] 1.2 添加 curl-raw 检测逻辑：输入以 `curl ` 开头则识别为 curl-raw
- [x] 1.3 添加 hex 格式检测逻辑：仅含 `0-9a-fA-F` 字符且长度为偶数
- [x] 1.4 添加 base64 格式检测逻辑：符合 base64 字符集且长度为 4 的倍数
- [x] 1.5 自动检测失败时，显示提示信息（小气泡提示用户手动选择类型）

## 2. Frontend - curl-raw 解析

- [x] 2.1 实现 `parseCurlRaw(input string)` 函数
- [x] 2.2 提取 `--data-raw` 参数值，支持单引号、双引号、无引号三种格式
- [x] 2.3 添加 `--data` 参数别名支持（`-d`）
- [x] 2.4 处理 Unicode 转义序列（`\u0000` → byte(0x00)，`\n` → 0x0a 等）
- [x] 2.5 将处理后的数据转为 base64 编码（grpc-web 帧头由后端处理）
- [x] 2.6 实现 `extractPathLastSegment(curlURL: string)` 函数，从 URL 提取路径最后一段

## 3. Frontend - UI 按钮

- [x] 3.1 在输入区域添加"清空输入"按钮，点击清空输入框
- [x] 3.2 在输入区域添加"全部清空"按钮，点击清空输入、类型选择、输出 JSON 和二进制数据
- [x] 3.3 样式调整：确保按钮布局合理，易于点击

## 4. Frontend - curl-raw 路径自动填充

- [x] 4.1 添加"自动填充路径"勾选框（默认勾选）
- [x] 4.2 curl-raw 解析成功后，根据勾选状态决定是否自动填充 Proto 类型搜索框
- [x] 4.3 填充时使用 `extractPathLastSegment` 提取的路径最后一段

## 5. Frontend - API 调用

- [x] 5.1 修改 `/api/proto/decode` 调用逻辑
- [x] 5.2 前端自动识别类型后，设置对应的 encoding（hex 或 base64）
- [x] 5.3 curl-raw 解析后，以 base64 编码发送数据

## 6. Backend - 无变更

后端保持现有实现不变，不做任何修改：
- 仅支持 `hex` 和 `base64` 两种 encoding
- curl-raw 的解析和帧头处理全部在前端完成