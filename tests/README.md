# SDK 测试

## 测试套件

### 单元测试

位于 `framework/` 目录下，测试SDK核心功能：

```bash
cd framework
go test -v
```

### 集成测试

位于 `tests/` 目录下，测试示例构建和SDK完整性：

```bash
cd tests
go test -v
```

## 测试覆盖

### 单元测试 (`framework/*_test.go`)

- ✅ Address 类型
- ✅ Event 创建和操作
- ✅ ContractError 错误处理
- ✅ 错误码常量
- ✅ ContractBase 基础方法
- ✅ 内存分配操作
- ✅ 状态读写操作

### 集成测试 (`tests/build_test.go`)

- ✅ 示例合约构建测试
- ✅ 示例目录结构验证
- ✅ SDK目录结构完整性检查
- ✅ WASM文件生成验证

## 运行所有测试

```bash
# 从SDK根目录运行
bash scripts/test-all.sh
```

## 注意事项

1. **TinyGo要求**: 构建测试需要安装TinyGo
2. **非WASM环境**: 单元测试在非WASM环境运行，宿主函数返回stub值
3. **集成测试**: 会实际编译WASM，确保示例可用

## 测试策略

- **单元测试**: 快速反馈，测试单个组件
- **集成测试**: 端到端验证，确保可交付
- **CI/CD**: 两种测试都应在CI中运行
