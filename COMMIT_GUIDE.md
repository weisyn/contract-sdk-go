# 提交指南

本文档说明如何提交所有 GitHub 配置文件到仓库。

## 📦 准备提交的文件

### .github 目录（7 个文件）

- `.github/ISSUE_TEMPLATE/bug_report.md` - Bug 报告模板
- `.github/ISSUE_TEMPLATE/feature_request.md` - 功能请求模板
- `.github/pull_request_template.md` - Pull Request 模板
- `.github/workflows/ci.yml` - CI 工作流
- `.github/workflows/lint.yml` - Lint 工作流
- `.github/dependabot.yml` - 依赖更新配置（可选）
- `.github/CODEOWNERS` - 代码所有者配置（可选）
- `.github/FUNDING.yml` - 资助配置（可选）

### 根目录文档（5 个文件）

- `LICENSE` - MIT License
- `CONTRIBUTING.md` - 贡献指南
- `CODE_OF_CONDUCT.md` - 行为准则
- `SECURITY.md` - 安全策略
- `GITHUB_SETUP.md` - GitHub 设置指南（参考文档）
- `COMMIT_GUIDE.md` - 提交指南（本文件）

### README.md 更新

- 添加了 Go Report Card 和 CI 徽章

## 🚀 提交命令

```bash
# 1. 查看所有变更
git status

# 2. 添加所有新文件
git add .github/
git add LICENSE
git add CONTRIBUTING.md
git add CODE_OF_CONDUCT.md
git add SECURITY.md
git add GITHUB_SETUP.md
git add COMMIT_GUIDE.md
git add README.md

# 3. 提交
git commit -m "chore: add GitHub templates, workflows, and documentation

- Add Issue templates (bug report, feature request)
- Add Pull Request template
- Add CI/CD workflows (test, lint)
- Add CONTRIBUTING.md, CODE_OF_CONDUCT.md, SECURITY.md
- Add MIT License
- Add dependabot.yml and CODEOWNERS (optional)
- Add FUNDING.yml (optional)
- Update README badges (Go Report Card, CI)
- Update contact email to wx@wesing.xyz
- Update domain references to weisyn.com"

# 4. 推送到 GitHub
git push origin main
```

## ⚠️ 注意事项

1. **CI 工作流**：需要安装 TinyGo，已配置 TinyGo Action
2. **Lint 工作流**：如果未配置 golangci-lint，可以暂时禁用或移除 lint.yml
3. **分支名称**：工作流中使用 `main` 和 `develop`，请确认仓库的实际分支名称
4. **CODEOWNERS**：需要将 `@weisyn` 替换为实际的 GitHub 用户名或团队名
5. **合约编译**：CI 中会测试合约编译，确保模板可以成功编译为 WASM

## 📋 提交后需要手动在 GitHub 设置的内容

以下内容需要在 GitHub 网页界面手动设置：

### Repository Settings → General

- [ ] **Description**: `WES Contract SDK for Go - 用于智能合约开发的 Go 语言 SDK，提供业务语义优先的合约开发框架，支持 WASM 编译`
- [ ] **Website**: `https://github.com/weisyn/contract-sdk-go#readme` 或 `https://weisyn.com`
- [ ] **Topics**: `blockchain wes sdk go golang contract-sdk smart-contract wasm tinygo ispc business-semantic utxo framework template`

### Repository Settings → Features

- [ ] **Issues**: 启用
- [ ] **Projects**: 可选启用
- [ ] **Wiki**: 可选启用
- [ ] **Discussions**: 可选启用

### Repository Settings → Branches

- [ ] 为 `main` 分支设置保护规则（参考 GITHUB_SETUP.md）

---

**最后更新**: 2025-11-23

