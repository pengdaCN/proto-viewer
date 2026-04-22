## Purpose

Proto 类型列表的分页和排序功能规范。

## ADDED Requirements

### Requirement: 分页的 Proto 类型列表

API SHALL 支持分页参数，返回排序后的类型列表。

#### Scenario: 获取第一页
- **WHEN** 客户端请求 `GET /api/proto/types?page=1&pageSize=50`
- **THEN** 系统 SHALL 返回第一页最多 50 条类型数据
- **AND** 响应 SHALL 包含 `types`, `total`, `page`, `pageSize`, `totalPages` 字段

#### Scenario: 获取后续页面
- **WHEN** 客户端请求 `GET /api/proto/types?page=2&pageSize=50`
- **THEN** 系统 SHALL 返回第二页的类型数据
- **AND** 响应 SHALL 反映当前页码和总页数

#### Scenario: 列表按名称排序
- **WHEN** 客户端请求任意页
- **THEN** 系统 SHALL 按类型名称字母顺序排序
- **AND** 排序 SHALL 是稳定的（相同名称保持一致顺序）

#### Scenario: 默认分页参数
- **WHEN** 客户端请求 `GET /api/proto/types` 不带分页参数
- **THEN** 系统 SHALL 默认返回第一页
- **AND** 系统 SHALL 使用默认 pageSize=50

#### Scenario: 空列表
- **WHEN** 后端没有加载任何 proto 类型
- **THEN** 系统 SHALL 返回空类型列表
- **AND** total SHALL 为 0
- **AND** totalPages SHALL 为 0