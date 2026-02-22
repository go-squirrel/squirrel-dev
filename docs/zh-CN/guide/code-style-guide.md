# 代码规范

## 1. 概述

本文档定义了项目的代码编写规范，旨在提高代码的可读性、可维护性和一致性。

## 2. 行长度限制

### 2.1 规则

一行代码尽量不超过 **120 个字符**。

### 2.2 说明

- 120 字符是一个合理的上限，既保证了代码的可读性，又避免了频繁换行
- 现代显示器普遍较宽，120 字符可以充分利用屏幕空间
- 超过限制时，应适当换行，保持代码清晰

### 2.3 示例

✅ 推荐：

```go
// 单行不超过 120 字符
if err := service.CreateDeployment(ctx, &deployment); err != nil {
    return err
}

// 长参数列表换行
result, err := client.Do(ctx, &Request{
    URL:    url,
    Method: method,
    Body:   body,
})
```

❌ 不推荐：

```go
// 单行过长，不易阅读
if err := service.CreateDeployment(ctx, &model.Deployment{ServerID: request.ServerID, ApplicationID: request.ApplicationID}); err != nil { return err }
```

## 3. 文件行数限制

### 3.1 规则

| 级别 | 行数 | 说明 |
|------|------|------|
| 建议 | ≤ 300 行 | 理想情况下，单个文件应保持在此范围内 |
| 硬性上限 | ≤ 500 行 | 绝对不能超过此限制 |

### 3.2 说明

- **可读性**: 文件过长会降低代码的可读性和理解难度
- **职责单一**: 文件过长通常意味着承担了过多职责，应考虑拆分
- **维护性**: 小文件更容易定位问题、修改和测试
- **代码审查**: 较小的文件更便于进行代码审查

### 3.3 拆分建议

当文件接近或超过限制时，考虑以下拆分策略：

1. **按职责拆分**: 将不同职责的代码拆分到不同文件
2. **按功能拆分**: 将相关功能聚合到独立的模块
3. **抽取工具函数**: 将通用逻辑抽取到独立的工具文件
4. **使用结构体组合**: 通过组合而非继承来组织代码

### 3.4 示例

假设有一个 `handler.go` 文件，包含了多个处理函数，文件行数超过 500 行：

```
handler/
├── handler.go      # 公共接口定义（< 100 行）
├── create.go       # 创建相关处理函数（< 200 行）
├── update.go       # 更新相关处理函数（< 150 行）
├── delete.go       # 删除相关处理函数（< 100 行）
└── list.go         # 列表查询相关处理函数（< 150 行）
```

## 4. 语言规范

### 4.1 规则

代码中**全部使用英文**，包括但不限于：

- 变量名、函数名、类型名
- 注释
- 日志信息
- 错误信息
- 配置文件中的键名

### 4.2 说明

- **国际化友好**: 英文是开源社区的通用语言，使用英文便于全球开发者参与
- **一致性**: 统一语言避免混合使用造成的混乱
- **工具兼容**: 大多数开发工具、文档系统对英文支持更好
- **搜索便利**: 使用英文更便于搜索引擎检索和问题排查

### 4.3 示例

✅ 推荐：

```go
// CreateDeployment creates a new deployment for the given application.
func (s *Service) CreateDeployment(ctx context.Context, req *CreateDeploymentRequest) error {
    if req.ServerID == "" {
        return errors.New("server ID is required")
    }
    log.Info("Creating deployment", "app", req.ApplicationID)
    return s.repo.Create(ctx, req)
}
```

❌ 不推荐：

```go
// 创建部署 为给定的应用程序创建新的部署
func (s *Service) CreateDeployment(ctx context.Context, req *CreateDeploymentRequest) error {
    if req.ServerID == "" {
        return errors.New("服务器ID不能为空")  // 错误信息使用中文
    }
    log.Info("创建部署", "app", req.ApplicationID)  // 日志使用中文
    return s.repo.Create(ctx, req)
}
```

## 5. 注意事项

1. **不要为了行数而牺牲清晰度**: 如果拆分会使代码更难理解，可以适当放宽建议限制（但绝不能超过硬性上限）
2. **注释和空行计入行数**: 文件行数包含注释和空行
3. **生成代码可豁免**: 自动生成的代码（如 protobuf 生成的代码）不受此限制
4. **定期审视**: 在代码审查时关注文件行数，及时进行拆分

## 6. 工具支持

可以使用以下工具检查代码规范：

```bash
# 检查行长
golines --max-len=120 --dry-run ./...

# 统计文件行数
find . -name "*.go" -exec wc -l {} \; | sort -rn | head -20
```
