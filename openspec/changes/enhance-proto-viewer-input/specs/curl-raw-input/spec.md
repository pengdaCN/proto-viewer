## ADDED Requirements

### Requirement: curl-raw 输入格式支持

系统 SHALL 支持直接输入 curl 命令格式，自动提取 `--data-raw` 参数中的数据进行处理。

#### Scenario: 解析标准 curl-raw 格式
- **WHEN** 用户输入以 `curl ` 开头的完整 curl 命令且包含 `--data-raw` 参数
- **THEN** 系统提取 `--data-raw` 的值作为实际数据

#### Scenario: 处理带引号的 --data-raw
- **WHEN** curl 命令中 `--data-raw` 使用单引号或双引号包裹
- **THEN** 系统正确提取引号内的数据内容

#### Scenario: 处理无引号的 --data-raw
- **WHEN** curl 命令中 `--data-raw` 直接跟随数据无引号
- **THEN** 系统正确提取数据至下一个空格或行尾

#### Scenario: 处理 Unicode 转义序列
- **WHEN** `--data-raw` 包含 `\u0000` 等 Unicode 转义
- **THEN** 系统将转义序列转换为实际字节进行解码

#### Scenario: 支持 --data 简写形式
- **WHEN** curl 命令使用 `--data` 而非 `--data-raw`
- **THEN** 系统同样支持提取数据

#### Scenario: curl-raw 自动跳过 grpc-web 帧头
- **WHEN** curl-raw 提取的数据以 5 字节 `\x00\x00\x00\x00\x00` 开头
- **THEN** 系统自动跳过这 5 字节帧头后对剩余数据进行 protobuf 解码