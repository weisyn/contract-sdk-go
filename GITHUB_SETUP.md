# GitHub 仓库设置指南

本文档说明如何完善 [contract-sdk-go](https://github.com/weisyn/contract-sdk-go) GitHub 仓库的设置。

## 📝 Repository Details（仓库详情）

### Description（描述）

```
WES Contract SDK for Go - 用于智能合约开发的 Go 语言 SDK，提供业务语义优先的合约开发框架，支持 WASM 编译
```

**英文版本**（如果支持英文描述）：
```
WES Contract SDK for Go - Go SDK for smart contract development, providing business-semantic-first contract framework with WASM compilation support
```

### Website（网站）

如果有独立的文档网站，填写文档网站地址。否则可以填写：

```
https://github.com/weisyn/contract-sdk-go#readme
```

或者主项目文档：
```
https://github.com/weisyn/weisyn
```

或者 WES 官网：
```
https://weisyn.com
```

### Topics（标签）

建议添加以下标签（用空格分隔）：

```
blockchain wes sdk go golang contract-sdk smart-contract wasm tinygo ispc business-semantic utxo framework template
```

**标签说明**：
- `blockchain` - 区块链相关
- `wes` - WES 区块链
- `sdk` - 软件开发工具包
- `go` / `golang` - Go 语言
- `contract-sdk` - 合约 SDK
- `smart-contract` - 智能合约
- `wasm` - WebAssembly
- `tinygo` - TinyGo
- `ispc` - ISPC（Intrinsic Self-Proving Computing）
- `business-semantic` - 业务语义
- `utxo` - UTXO 模型
- `framework` - 框架
- `template` - 模板

## 🔧 其他需要完善的内容

### 1. LICENSE 文件

**当前状态**：README 中提到了 MIT License，但仓库中可能缺少 LICENSE 文件。

**操作步骤**：
1. 在 GitHub 仓库页面，点击 "Add file" → "Create new file"
2. 文件名输入 `LICENSE`
3. 选择 "Choose a license template" → 选择 "MIT License"
4. 填写版权信息（如：Copyright 2025 Weisyn）
5. 提交文件

### 2. .github 目录结构

建议创建以下文件以提升仓库专业性：

#### 2.1 Issue 模板

创建 `.github/ISSUE_TEMPLATE/` 目录：

**bug_report.md**：
```markdown
---
name: Bug Report
about: 报告合约 SDK 的 bug
title: '[BUG] '
labels: bug
assignees: ''
---

## 描述
简要描述 bug

## 复现步骤
1. 
2. 
3. 

## 预期行为
描述预期行为

## 实际行为
描述实际行为

## 环境信息
- Go 版本：
- TinyGo 版本：
- SDK 版本：
- 操作系统：

## 日志/错误信息
```
粘贴错误日志
```

## 附加信息
其他相关信息
```

**feature_request.md**：
```markdown
---
name: Feature Request
about: 提出新功能建议
title: '[FEATURE] '
labels: enhancement
assignees: ''
---

## 功能描述
简要描述新功能

## 使用场景
描述使用场景

## 建议实现
描述建议的实现方式

## 附加信息
其他相关信息
```

#### 2.2 Pull Request 模板

创建 `.github/pull_request_template.md`：
```markdown
## 变更描述
简要描述本次 PR 的变更

## 变更类型
- [ ] Bug 修复
- [ ] 新功能
- [ ] 文档更新
- [ ] 代码重构
- [ ] 测试相关
- [ ] 模板更新
- [ ] 其他

## 测试
- [ ] 已添加单元测试
- [ ] 已测试合约编译
- [ ] 已通过所有测试

## 检查清单
- [ ] 代码遵循项目规范
- [ ] 已更新相关文档
- [ ] 已添加必要的注释
- [ ] 无编译错误和警告
- [ ] WASM 编译成功
```

#### 2.3 GitHub Actions Workflows

创建 `.github/workflows/` 目录：

**ci.yml**（持续集成）：
```yaml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Install TinyGo
        uses: tinygo-org/tinygo-action@v1
        with:
          version: 0.31.0
      
      - name: Run tests
        run: go test ./... -v
      
      - name: Build examples
        run: |
          cd templates/learning/hello-world
          ./build.sh
```

**lint.yml**（代码检查）：
```yaml
name: Lint

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  golangci-lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
```

### 3. CONTRIBUTING.md

创建贡献指南文件（详见 CONTRIBUTING.md）

### 4. CODE_OF_CONDUCT.md

创建行为准则文件（详见 CODE_OF_CONDUCT.md）

### 5. SECURITY.md

创建安全策略文件（详见 SECURITY.md）

### 6. COMMIT_GUIDE.md

创建提交指南文件（详见 COMMIT_GUIDE.md）

### 7. 仓库设置检查清单

#### 7.1 General Settings（常规设置）

- [ ] **Repository name**: `contract-sdk-go` ✅
- [ ] **Description**: 填写上述描述
- [ ] **Website**: 填写上述网站
- [ ] **Topics**: 添加上述标签

#### 7.2 Features（功能）

- [ ] **Issues**: 启用（用于 bug 报告和功能请求）
- [ ] **Projects**: 可选启用（用于项目管理）
- [ ] **Wiki**: 可选启用（如果使用 Wiki 文档）
- [ ] **Discussions**: 可选启用（用于社区讨论）
- [ ] **Sponsorships**: 可选启用

#### 7.3 Branch Protection（分支保护）

建议为 `main` 分支设置保护规则：

- [ ] **Require a pull request before merging**
  - [ ] Require approvals: 1
  - [ ] Dismiss stale pull request approvals when new commits are pushed
- [ ] **Require status checks to pass before merging**
  - [ ] Require branches to be up to date before merging
  - [ ] 选择 CI 检查（如：test, lint）
- [ ] **Require conversation resolution before merging**
- [ ] **Do not allow bypassing the above settings**

#### 7.4 Pages（页面）

如果有文档网站：

- [ ] 启用 GitHub Pages
- [ ] 选择文档源（如：`/docs` 目录）

#### 7.5 Actions（Actions）

- [ ] 确保 Actions 已启用
- [ ] 检查 Actions 权限设置

### 8. README 徽章

README 中已有一些徽章，可以添加更多：

```markdown
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![TinyGo](https://img.shields.io/badge/TinyGo-0.31+-blue.svg)](https://tinygo.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/weisyn/contract-sdk-go)](https://goreportcard.com/report/github.com/weisyn/contract-sdk-go)
[![CI](https://github.com/weisyn/contract-sdk-go/workflows/CI/badge.svg)](https://github.com/weisyn/contract-sdk-go/actions)
```

### 9. 仓库描述优化建议

**简短版本**（适合 GitHub 搜索）：
```
WES Contract SDK for Go - Business-semantic-first smart contract framework with WASM support
```

**详细版本**（适合 About 部分）：
```
WES Contract SDK for Go 是一个用于开发 WES 智能合约的 Go 语言 SDK。提供业务语义优先的合约开发框架，让开发者专注于业务逻辑而非底层细节。支持 TinyGo WASM 编译，零外部依赖，提供丰富的合约模板和示例。
```

## 📋 完整检查清单

### 必须完成

- [ ] Description（描述）
- [ ] Website（网站）
- [ ] Topics（标签）
- [ ] LICENSE 文件
- [ ] README.md（已有 ✅）

### 推荐完成

- [ ] .github/ISSUE_TEMPLATE/（Issue 模板）
- [ ] .github/pull_request_template.md（PR 模板）
- [ ] .github/workflows/（CI/CD 工作流）
- [ ] CONTRIBUTING.md（贡献指南）
- [ ] CODE_OF_CONDUCT.md（行为准则）
- [ ] SECURITY.md（安全策略）
- [ ] COMMIT_GUIDE.md（提交指南）

### 可选完成

- [ ] .github/FUNDING.yml（资助信息）
- [ ] .github/dependabot.yml（依赖更新）
- [ ] .github/CODEOWNERS（代码所有者）
- [ ] GitHub Pages（文档网站）
- [ ] Releases（发布版本）

---

**最后更新**: 2025-01-23

