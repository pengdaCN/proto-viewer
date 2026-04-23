## ADDED Requirements

### Requirement: 清空输入按钮

系统 SHALL 提供"清空输入"按钮，点击后清空输入框内容但保留其他设置。

#### Scenario: 点击清空输入按钮
- **WHEN** 用户点击"清空输入"按钮
- **THEN** 系统清空输入框内容
- **AND** 类型选择、输出 JSON、二进制数据保持不变

### Requirement: 全部清空按钮

系统 SHALL 提供"全部清空"按钮，点击后清空所有相关数据。

#### Scenario: 点击全部清空按钮
- **WHEN** 用户点击"全部清空"按钮
- **THEN** 系统清空输入框内容
- **AND** 类型选择恢复默认状态
- **AND** 输出的 JSON 内容清空
- **AND** 二进制数据显示区域清空