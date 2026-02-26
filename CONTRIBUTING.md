# 贡献指南

感谢你对 Goon 项目的关注！我们欢迎所有形式的贡献。

## 如何贡献

### 报告 Bug

如果你发现了 bug，请创建一个 issue 并包含以下信息：

- Bug 的详细描述
- 复现步骤
- 预期行为
- 实际行为
- 你的环境信息（操作系统、Go 版本等）

### 提交功能请求

如果你有新功能的想法，请创建一个 issue 并描述：

- 功能的用途和价值
- 预期的使用方式
- 可能的实现方案

### 提交代码

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建一个 Pull Request

## 开发指南

### 环境设置

```bash
# 克隆仓库
git clone https://github.com/yourusername/goon.git
cd goon

# 安装依赖
go mod download

# 运行测试
make test

# 构建项目
make build
```

### 代码规范

- 遵循 Go 官方代码风格指南
- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 进行代码检查
- 为新功能添加测试
- 保持测试覆盖率在 80% 以上

### 提交信息规范

使用清晰的提交信息：

```
<type>: <subject>

<body>

<footer>
```

类型（type）：
- feat: 新功能
- fix: Bug 修复
- docs: 文档更新
- style: 代码格式调整
- refactor: 重构
- test: 测试相关
- chore: 构建/工具相关

### 测试

- 为所有新功能编写单元测试
- 确保所有测试通过：`make test`
- 检查测试覆盖率：`make test-coverage`

### 文档

- 更新相关文档
- 为公共 API 添加注释
- 更新 README.md（如果需要）

## Pull Request 流程

1. 确保你的代码通过所有测试
2. 更新相关文档
3. 在 PR 描述中清楚地说明你的更改
4. 等待代码审查
5. 根据反馈进行修改

## 行为准则

- 尊重所有贡献者
- 保持友好和专业
- 接受建设性的批评
- 关注对项目最有利的事情

## 许可证

通过贡献代码，你同意你的贡献将在 MIT 许可证下发布。
