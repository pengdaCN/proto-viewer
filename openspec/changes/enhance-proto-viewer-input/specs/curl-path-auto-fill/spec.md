## ADDED Requirements

### Requirement: curl-raw 路径自动填充

当用户输入 curl-raw 格式时，系统 SHALL 自动从 URL 路径提取最后一段，填入 Proto 类型搜索框。

#### Scenario: 从 curl URL 提取类型名
- **WHEN** curl 命令 URL 为 `https://console.mpcvault.com/admin/mpcvault.internal.admin.v1.AdminService/CheckCredential`
- **THEN** 系统提取 `CheckCredential` 作为搜索关键词

#### Scenario: 自动填充到 Proto 类型搜索框
- **WHEN** curl-raw 解析成功后且自动填充开关已勾选
- **THEN** 系统自动将提取的路径最后一段填入 Proto 类型搜索框

### Requirement: 自动填充开关

系统 SHALL 提供勾选框控制是否自动执行路径填充功能，默认勾选。

#### Scenario: 默认勾选状态
- **WHEN** 用户首次使用或未更改设置
- **THEN** 自动填充功能默认开启

#### Scenario: 取消勾选则不自动填充
- **WHEN** 用户取消勾选"自动填充"选项
- **THEN** curl-raw 解析后不自动填充搜索框

#### Scenario: 手动输入优先级
- **WHEN** 用户在自动填充后手动修改了搜索框内容
- **THEN** 系统使用用户手动输入的值，不进行覆盖