## Purpose

Proto Debugger 前端页面功能规范。

## ADDED Requirements

### Requirement: 前端分页控件

前端 SHALL 提供分页控件，支持在类型列表中翻页。

#### Scenario: 显示分页信息
- **WHEN** 用户打开 proto-debugger 页面
- **THEN** 系统 SHALL 显示当前页码和总页数
- **AND** 系统 SHALL 显示"上一页"和"下一页"按钮

#### Scenario: 点击下一页
- **WHEN** 用户点击"下一页"按钮
- **AND** 当前不是最后一页
- **THEN** 系统 SHALL 加载并显示下一页数据
- **AND** 系统 SHALL 更新页码显示

#### Scenario: 点击上一页
- **WHEN** 用户点击"上一页"按钮
- **AND** 当前不是第一页
- **THEN** 系统 SHALL 加载并显示上一页数据
- **AND** 系统 SHALL 更新页码显示

#### Scenario: 第一页时禁用上一页
- **WHEN** 用户在第一页
- **THEN** 系统 SHALL 禁用"上一页"按钮
- **AND** 按钮 SHALL 不可点击

#### Scenario: 最后一页时禁用下一页
- **WHEN** 用户在最后一页
- **THEN** 系统 SHALL 禁用"下一页"按钮
- **AND** 按钮 SHALL 不可点击

## MODIFIED Requirements

### Requirement: Proto 类型列表展示

系统 SHALL 在页面加载时从后端获取可用的 proto 类型列表，并支持分页展示供用户选择。

#### Scenario: 页面加载时获取类型列表（第一页）
- **WHEN** 用户打开 proto-debugger 页面
- **THEN** 系统 SHALL 从 `/api/proto/types?page=1&pageSize=50` 获取类型列表
- **AND** 系统 SHALL 将第一页类型填充到下拉选择框中
- **AND** 系统 SHALL 显示分页控件

#### Scenario: 类型列表为空时
- **WHEN** 后端返回空的类型列表
- **THEN** 系统 SHALL 在页面显示提示信息"无可用的 proto 类型，请检查 proto 定义文件夹配置"