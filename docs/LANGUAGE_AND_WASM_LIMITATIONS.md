# 语言与 WASM 环境限制 - Go/TinyGo

**版本**: v1.0.0  
**最后更新**: 2025-11-23

---

## 📋 文档定位

> 📌 **重要说明**：本文档说明 Go/TinyGo 特有的限制和注意事项。  
> 这是 Go/TinyGo 特有、且非常 SDK 本地化的内容。

**本文档目标**：
- 说明 TinyGo 支持矩阵
- 说明 Unsafe 指针注意事项
- 说明 WASM 环境限制
- 说明禁止使用的标准库

---

## ⚠️ TinyGo 支持矩阵

### 支持的标准库

- ✅ `fmt` - 格式化输出（部分功能）
- ✅ `strings` - 字符串操作（部分功能）
- ✅ `bytes` - 字节操作
- ✅ `encoding/json` - JSON 编码/解码（部分功能）
- ❌ `net` - 网络操作（不支持）
- ❌ `os` - 操作系统操作（不支持）
- ❌ `time` - 时间操作（部分支持）

### 完整支持矩阵

参考 [TinyGo 官方文档](https://tinygo.org/docs/reference/lang-support/stdlib/) 获取最新支持情况。

---

## 🔧 Unsafe 指针注意事项

### 为什么需要 Unsafe 指针

在 WASM 环境中，Go 的指针和内存管理与标准 Go 不同。SDK 需要使用 `unsafe.Pointer` 来：
- 与 HostABI 交互
- 处理内存布局
- 优化性能

### 使用示例

```go
import "unsafe"

// SDK 内部使用 unsafe.Pointer 与 HostABI 交互
// 这是 WASM/WebAssembly 开发的标准做法
```

### 安全建议

- ✅ **SDK 内部使用**：SDK 内部已经添加了适当的安全检查
- ❌ **合约代码避免使用**：合约开发者应避免直接使用 `unsafe.Pointer`
- ✅ **使用 Helpers API**：优先使用 Helpers 层的业务语义接口

---

## 🚫 WASM 环境限制

### 禁止使用的标准库

- ❌ `net` - 网络操作
- ❌ `os` - 操作系统操作
- ❌ `syscall` - 系统调用
- ❌ `runtime` - 运行时操作（部分）

### 限制的功能

- ❌ **Goroutine**：不支持并发
- ❌ **Channel**：不支持通道
- ❌ **Reflection**：反射功能受限
- ❌ **GC**：垃圾回收机制不同

---

## 💡 最佳实践

### 推荐做法

1. **使用 Helpers API**：优先使用 Helpers 层的业务语义接口
2. **避免复杂操作**：避免使用不支持的标准库
3. **参考模板**：参考 SDK 提供的模板代码

### 不推荐做法

1. ❌ **直接使用 unsafe**：避免在合约代码中直接使用 `unsafe.Pointer`
2. ❌ **使用不支持的标准库**：避免使用 `net`、`os` 等不支持的标准库
3. ❌ **复杂并发**：避免使用 Goroutine 和 Channel

---

## 📖 进一步阅读

### 核心文档

- **[开发者指南](./DEVELOPER_GUIDE.md)** - 如何使用 Go SDK 开发合约
- **[API 参考](./API_REFERENCE.md)** - 详细的 API 文档

### 外部资源

- [TinyGo 官方文档](https://tinygo.org/) - TinyGo 完整文档
- [TinyGo 标准库支持](https://tinygo.org/docs/reference/lang-support/stdlib/) - 标准库支持矩阵
- [Go unsafe.Pointer 文档](https://pkg.go.dev/unsafe#Pointer) - Go unsafe.Pointer 文档

---

**最后更新**: 2025-11-23  
**维护者**: WES Core Team

