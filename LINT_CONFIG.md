# Lint 配置说明

本文档说明 WES Contract SDK for Go 的代码检查配置。

## 配置文件

- `.golangci.yml` - golangci-lint 配置（根目录）
- `framework/.golangci.yml` - framework 目录的特定配置（已废弃，统一使用根目录配置）

## 使用方式

### 本地开发

```bash
# 安装 golangci-lint（如果还没有）
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 运行 lint 检查
golangci-lint run

# 自动修复部分问题
golangci-lint run --fix
```

### CI/CD

GitHub Actions 会自动运行 lint 检查，确保代码质量。

## 配置说明

### 启用的 Linters

- `errcheck` - 检查错误处理
- `govet` - Go 官方 vet 工具
- `ineffassign` - 检查无效赋值
- `staticcheck` - 静态分析
- `unused` - 检查未使用的代码
- `gocritic` - Go 代码审查
- `gosec` - 安全检查（允许 unsafe.Pointer，WASM 需要）
- `misspell` - 拼写检查
- `unparam` - 检查未使用的参数
- `revive` - Go linter
- `prealloc` - 预分配检查
- `typecheck` - 类型检查

### 特殊配置

由于合约 SDK 需要编译为 WASM，配置中：
- 允许 `unsafe.Pointer` 的使用（gosec 排除 G103）
- 排除 framework 目录中与 WASM 相关的警告

### 与主项目保持一致

本配置基于 WES 主项目（weisyn.git）的 `.golangci.yml`，但针对合约 SDK 项目进行了适配：
- 保留了核心的代码质量标准
- 允许 WASM 环境必需的特性（unsafe.Pointer）
- 确保与主项目的代码风格一致

---

**最后更新**: 2025-01-23

