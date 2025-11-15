# Framework Internal 包

**版本**: v1.0  
**创建日期**: 2025-11-11  
**状态**: ✅ 已实施

---

## ⚠️ 重要说明

**此包仅供 helpers 层使用，外部开发者不应导入。**

根据Go语言的internal包规则，此包只能被 `github.com/weisyn/contract-sdk-go` 模块内的包导入，外部模块无法导入。

---

## 📋 包结构

```
framework/internal/
├── transaction.go    # 链式API实现
├── hostabi.go        # HostABI原语封装
├── state.go          # 状态输出相关
├── resource.go       # 资源输出相关
├── utxo.go           # UTXO操作相关
├── batch.go          # 批量操作相关
└── README.md         # 本文档
```

---

## 🔍 接口列表

### transaction.go - 链式API

- `BeginTransaction()` - 开始交易构建
- `AddAssetOutput()` - 添加资产输出
- `AddStateOutput()` - 添加状态输出
- `AddResourceOutput()` - 添加资源输出
- `AddInput()` - 添加交易输入
- `Transfer()` - 添加转账意图
- `Stake()` - 添加质押意图
- `WithFee()` - 设置费用偏好
- `Finalize()` - 完成交易构建

### hostabi.go - HostABI原语

- `TxAddInput()` - 添加交易输入
- `TxAddAssetOutput()` - 添加资产输出
- `TxAddStateOutput()` - 添加状态输出
- `TxAddResourceOutput()` - 添加资源输出

### state.go - 状态输出

- `AppendStateOutput()` - 追加状态输出
- `AppendStateOutputSimple()` - 追加状态输出（简化）

### resource.go - 资源输出

- `AppendResourceOutput()` - 追加资源输出

### utxo.go - UTXO操作

- `CreateUTXO()` - 创建UTXO（原语函数）
- `CreateAssetOutputWithLock()` - 创建带锁定的资产输出（原语函数）
- ⚠️ `ExecuteUTXOTransferEx()` - 已删除，请使用 `TransactionBuilder.Transfer()` 或 `helpers/token/Transfer`

### batch.go - 批量操作

- `BatchCreateOutputs()` - 批量创建输出
- `BatchCreateOutputsSimple()` - 批量创建输出（简化）

---

## 📖 使用示例

### helpers层使用

```go
package token

import (
    "github.com/weisyn/contract-sdk-go/framework"
    "github.com/weisyn/contract-sdk-go/framework/internal"
)

func Transfer(from, to framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
    // 使用internal包
    success, _, errCode := internal.BeginTransaction().
        Transfer(from, to, tokenID, amount).
        Finalize()
    
    if !success {
        return framework.NewContractError(errCode, "transfer failed")
    }
    
    return nil
}
```

---

## 🚫 外部开发者无法使用

```go
// 外部代码（无法编译）
package main

import (
    "github.com/weisyn/contract-sdk-go/framework/internal" // ❌ 编译错误
)

func main() {
    internal.BeginTransaction() // 无法访问
}
```

**编译错误**：
```
cannot import internal package "github.com/weisyn/contract-sdk-go/framework/internal"
```

---

## ✅ 验证

### 验证1：helpers可以使用

✅ **通过**：helpers层可以正常导入和使用internal包

### 验证2：外部开发者无法导入

✅ **通过**：Go语言的internal包机制确保外部模块无法导入

---

## 🔗 相关文档

- [Framework 层文档](../README.md) - Framework 层详细说明
- [Helpers 层文档](../../helpers/README.md) - Helpers 层详细说明
- [SDK 主 README](../../README.md) - SDK 总览

---

**最后更新**: 2025-11-11

