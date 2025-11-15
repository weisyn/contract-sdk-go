# Framework 层 - HostABI 封装

**版本**: 1.0  
**状态**: ✅ 稳定  
**最后更新**: 2025-11-11

---

## 📋 概述

Framework 层是 Contract SDK 的核心框架，提供 HostABI 原语的 Go 语言封装和基础能力。它为上层业务语义层（helpers）提供类型安全的 API 和统一的错误处理。

**定位**：Framework 层是 SDK 的内部实现层，合约开发者通常**不需要直接使用**，应优先使用 `helpers` 层的业务语义接口。

---

## 🎯 核心功能

### 环境查询

获取执行上下文和区块信息：

```go
import "github.com/weisyn/contract-sdk-go/framework"

// 执行上下文
caller := framework.GetCaller()              // 调用者地址
contractAddr := framework.GetContractAddress() // 合约地址
txID := framework.GetTransactionID()         // 交易ID

// 区块信息
height := framework.GetBlockHeight()         // 区块高度
timestamp := framework.GetBlockTimestamp()   // 区块时间戳
blockHash := framework.GetBlockHash(height)  // 区块哈希
chainID := framework.GetChainID()           // 链ID

// 余额查询（账户抽象）
balance := framework.QueryUTXOBalance(address, tokenID)
```

### 事件与日志

```go
import "github.com/weisyn/contract-sdk-go/framework"

// 发出链上事件
framework.EmitEvent("Transfer", []byte(`{"from":"...","to":"...","amount":100}`))

// 记录调试日志
framework.LogDebug("Processing transfer...")
```

### 参数解析

```go
import "github.com/weisyn/contract-sdk-go/framework"

// 获取合约参数
params := framework.GetContractParams()

// 解析JSON参数
toStr := params.ParseJSON("to")
amount := params.ParseJSONInt("amount")
support := params.ParseJSONBool("support")

// 解析地址
to, err := framework.ParseAddressBase58(toStr)
if err != nil {
    return framework.ERROR_INVALID_PARAMS
}
```

### 返回值设置

```go
import "github.com/weisyn/contract-sdk-go/framework"

// 设置返回数据
framework.SetReturnData([]byte("result"))

// 设置JSON返回值
framework.SetReturnJSON(map[string]interface{}{
    "balance": 1000,
    "token_id": "my_token",
})
```

---

## 📐 架构定位

```
┌─────────────────────────────────────────┐
│  业务语义层（合约开发者使用）            │
│  helpers/                               │
│  ├─ token.Transfer()                    │
│  ├─ staking.Stake()                     │
│  └─ governance.Vote()                   │
└─────────────────────────────────────────┘
         ↓ 基于
┌─────────────────────────────────────────┐
│  框架层（本层，SDK内部使用）             │
│  framework/                             │
│  ├─ HostABI 封装                       │
│  ├─ 交易构建（内部）                    │
│  └─ 状态管理（内部）                    │
└─────────────────────────────────────────┘
         ↓ 调用
┌─────────────────────────────────────────┐
│  WES 协议层（底层能力）                 │
│  - HostABI 原语                         │
│  - EUTXO 交易模型                      │
└─────────────────────────────────────────┘
```

---

## 💡 使用建议

### ✅ 推荐：使用 Helpers 层

**合约开发者应优先使用 `helpers` 层的业务语义接口**：

```go
import "github.com/weisyn/contract-sdk-go/helpers/token"

// 推荐：使用业务语义接口
err := token.Transfer(from, to, tokenID, amount)
```

### ⚠️ 谨慎：直接使用 Framework 层

**仅在以下情况直接使用 Framework 层**：
- 需要环境查询（GetCaller、GetBlockHeight 等）
- 需要事件和日志（EmitEvent、LogDebug）
- 需要参数解析（GetContractParams）
- 需要自定义底层操作（不推荐，应优先考虑使用 helpers）

```go
import "github.com/weisyn/contract-sdk-go/framework"

// 可以：环境查询
caller := framework.GetCaller()

// 可以：事件和日志
framework.EmitEvent("CustomEvent", eventData)

// 不推荐：直接构建交易（应使用 helpers）
// framework.BeginTransaction()...
```

---

## 📚 核心类型

### 基础类型

```go
type Address [20]byte        // 地址类型
type Hash [32]byte          // 哈希类型
type TokenID []byte         // 代币ID
type Amount uint64          // 金额类型
```

### 错误码

```go
const (
    SUCCESS                  = 0  // 成功
    ERROR_INVALID_PARAMS     = 1  // 参数错误
    ERROR_INSUFFICIENT_BALANCE = 2  // 余额不足
    ERROR_UNAUTHORIZED       = 3  // 未授权
    ERROR_NOT_FOUND          = 4  // 未找到
    ERROR_ALREADY_EXISTS     = 5  // 已存在
    ERROR_EXECUTION_FAILED   = 6  // 执行失败
    ERROR_INVALID_STATE      = 7  // 无效状态
    ERROR_TIMEOUT            = 8  // 超时
    ERROR_NOT_IMPLEMENTED    = 9  // 未实现
    ERROR_PERMISSION_DENIED  = 10 // 权限拒绝
)
```

---

## 🔧 使用示例

### 示例1：查询余额

```go
package main

import "github.com/weisyn/contract-sdk-go/framework"

//export GetBalance
func GetBalance() uint32 {
    params := framework.GetContractParams()
    addrStr := params.ParseJSON("address")
    
    addr, err := framework.ParseAddressBase58(addrStr)
    if err != nil {
        return framework.ERROR_INVALID_PARAMS
    }
    
    balance := framework.QueryUTXOBalance(addr, nil)
    framework.SetReturnData([]byte(fmt.Sprintf("%d", balance)))
    
    return framework.SUCCESS
}
```

### 示例2：发出事件

```go
package main

import "github.com/weisyn/contract-sdk-go/framework"

//export CustomAction
func CustomAction() uint32 {
    caller := framework.GetCaller()
    
    // 发出事件
    eventData := []byte(`{"caller":"` + string(caller) + `","action":"custom"}`)
    framework.EmitEvent("CustomAction", eventData)
    
    return framework.SUCCESS
}
```

---

## 🔗 相关文档

- [Contract SDK 主 README](../README.md) - SDK 总览
- [Helpers 层文档](../helpers/README.md) - 业务语义层（推荐使用）
- [HostABI 原语能力](../../../docs/components/core/ispc/capabilities/hostabi-primitives.md) - 底层原语说明

---

## ⚠️ 注意事项

1. **优先使用 Helpers 层**：Framework 层是 SDK 的内部实现，合约开发者应优先使用 `helpers` 层的业务语义接口
2. **类型安全**：使用 Framework 提供的类型（Address、Amount、TokenID 等），避免使用原始类型
3. **错误处理**：使用统一的错误码，便于错误处理和调试
4. **事件和日志**：合理使用事件和日志，避免过度使用影响性能
5. **TinyGo WASM 环境限制**：
   - ❌ **不支持标准库 `encoding/json`**：TinyGo WASM环境不支持完整的`encoding/json`包
   - ✅ **使用SDK提供的JSON工具**：使用`ContractParams.ParseJSON()`等方法
   - ⚠️ **限制**：仅支持基本字段提取，不支持完整JSON解析
   - 📚 **更多信息**：参考 [WASM 环境说明](../../docs/tutorials/contracts/wasm-environment.md)

---

**最后更新**: 2025-11-11
