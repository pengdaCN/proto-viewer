## ADDED Requirements

### Requirement: 输入数据类型自动检测

系统 SHALL 支持自动检测输入数据类型，根据内容特征判断是 hex、base64 还是 curl-raw 编码，无需用户手动选择。

#### Scenario: 检测 curl-raw 格式输入
- **WHEN** 用户输入以 `curl ` 开头
- **THEN** 系统自动识别为 curl-raw 格式并解析

#### Scenario: 检测 hex 格式输入
- **WHEN** 用户输入仅包含 `0-9a-fA-F` 字符且长度为偶数
- **THEN** 系统自动识别为 hex 格式并进行解码

#### Scenario: 检测 base64 格式输入
- **WHEN** 用户输入符合 base64 字符集（`A-Za-z0-9+/=`）且长度为 4 的倍数
- **THEN** 系统自动识别为 base64 格式并进行解码

#### Scenario: 自动检测失败时回退
- **WHEN** 自动检测无法确定数据类型
- **THEN** 系统返回错误提示，明确要求用户选择类型

#### Scenario: 手动选择优先级
- **WHEN** 用户明确选择了编码类型（hex/base64）
- **THEN** 系统使用用户选择的类型，不进行自动检测