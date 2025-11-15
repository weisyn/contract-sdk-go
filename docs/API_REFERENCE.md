# WES 合约SDK API参考文档

**版本**: v1.0.0  
**状态**: ✅ 稳定  
**最后更新**: 2025-11-11

---

## 🎯 快速开始

> **⚠️ 重要提示**: Framework 层是 SDK 的内部实现层，合约开发者**应优先使用 Helpers 层的业务语义接口**。Framework 层主要用于环境查询、事件发出等基础能力。

### 推荐方式：使用 Helpers 层

```go
import (
    "github.com/weisyn/contract-sdk-go/helpers/token"
    "github.com/weisyn/contract-sdk-go/framework"
)

//export Transfer
func Transfer() uint32 {
    params := framework.GetContractParams()
    toStr := params.ParseJSON("to")
    amount := params.ParseJSONInt("amount")
    
    to, err := framework.ParseAddressBase58(toStr)
    if err != nil {
        return framework.ERROR_INVALID_PARAMS
    }
    
    caller := framework.GetCaller()
    err = token.Transfer(caller, to, nil, framework.Amount(amount))
    if err != nil {
        return framework.ERROR_EXECUTION_FAILED
    }
    
    return framework.SUCCESS
}
```

### Framework 层使用场景

Framework 层主要用于：
- 环境查询（GetCaller、GetBlockHeight 等）
- 事件和日志（EmitEvent、LogDebug）
- 参数解析（GetContractParams）

---

## 📚 核心API

> **⚠️ 注意**: Framework 层主要用于环境查询、事件发出等基础能力。**交易构建相关的 API（如 TransactionBuilder）是内部实现，合约开发者应使用 Helpers 层的业务语义接口**。

### 1. 环境查询函数

#### GetCaller()

获取调用者地址。

**签名**:
```go
func GetCaller() Address
```

**返回值**:
- `Address`: 调用者地址

**示例**:
```go
caller := framework.GetCaller()
```

---

#### GetContractAddress()

获取当前合约地址。

**签名**:
```go
func GetContractAddress() Address
```

**返回值**:
- `Address`: 合约地址

**示例**:
```go
contractAddr := framework.GetContractAddress()
```

---

#### GetTransactionID()

获取当前交易ID。

**签名**:
```go
func GetTransactionID() []byte
```

**返回值**:
- `[]byte`: 交易ID（32字节）

**示例**:
```go
txID := framework.GetTransactionID()
```

---

#### GetBlockHeight()

获取当前区块高度。

**签名**:
```go
func GetBlockHeight() uint64
```

**返回值**:
- `uint64`: 区块高度

**示例**:
```go
height := framework.GetBlockHeight()
```

---

#### GetBlockTimestamp()

获取当前区块时间戳。

**签名**:
```go
func GetBlockTimestamp() uint64
```

**返回值**:
- `uint64`: 区块时间戳

**示例**:
```go
timestamp := framework.GetBlockTimestamp()
```

---

#### QueryUTXOBalance()

查询UTXO余额（账户抽象）。

**签名**:
```go
func QueryUTXOBalance(
    owner Address,
    tokenID TokenID,
) Amount
```

**参数**:
- `owner`: 地址
- `tokenID`: 代币ID（nil 表示原生币）

**返回值**:
- `Amount`: 余额

**示例**:
```go
balance := framework.QueryUTXOBalance(addr, nil)
```

---

### 2. 参数解析

#### GetContractParams()

获取合约参数。

**签名**:
```go
func GetContractParams() *ContractParams
```

**返回值**:
- `*ContractParams`: 参数对象

**示例**:
```go
params := framework.GetContractParams()
toStr := params.ParseJSON("to")
amount := params.ParseJSONInt("amount")
```

---

#### ParseAddressBase58()

解析Base58编码的地址。

**签名**:
```go
func ParseAddressBase58(addrStr string) (Address, error)
```

**参数**:
- `addrStr`: Base58编码的地址字符串

**返回值**:
- `Address`: 地址对象
- `error`: 错误信息

**示例**:
```go
addr, err := framework.ParseAddressBase58(addrStr)
if err != nil {
    return framework.ERROR_INVALID_PARAMS
}
```

---

### 3. 事件与日志

#### EmitEvent()

发出链上事件。

**签名**:
```go
func EmitEvent(event *Event) error
```

**参数**:
- `event`: 事件对象

**示例**:
```go
event := framework.NewEvent("Transfer")
event.AddAddressField("from", from)
event.AddAddressField("to", to)
event.AddUint64Field("amount", amount)
framework.EmitEvent(event)
```

---

#### LogDebug()

记录调试日志。

**签名**:
```go
func LogDebug(message string)
```

**参数**:
- `message`: 日志消息

**示例**:
```go
framework.LogDebug("Processing transfer...")
```

---

### 4. 返回值设置

#### SetReturnData()

设置返回数据。

**签名**:
```go
func SetReturnData(data []byte)
```

**参数**:
- `data`: 返回数据

**示例**:
```go
framework.SetReturnData([]byte("result"))
```

---

#### SetReturnJSON()

设置JSON返回值。

**签名**:
```go
func SetReturnJSON(data map[string]interface{})
```

**参数**:
- `data`: JSON数据

**示例**:
```go
framework.SetReturnJSON(map[string]interface{}{
    "balance": 1000,
    "token_id": "my_token",
})
```

---

#### GetCaller()

获取调用者地址。

**签名**:
```go
func GetCaller() Address
```

**返回值**:
- `Address`: 调用者地址

**示例**:
```go
caller := framework.GetCaller()
```

---

#### GetContractAddress()

获取当前合约地址。

**签名**:
```go
func GetContractAddress() Address
```

**返回值**:
- `Address`: 合约地址

**示例**:
```go
contractAddr := framework.GetContractAddress()
```

---

### 4. 事件系统

#### EmitEvent()

发出事件。

**签名**:
```go
func EmitEvent(event *Event) error
```

**参数**:
- `event`: 事件对象

**示例**:
```go
event := framework.NewEvent("Transfer")
event.AddAddressField("from", from)
event.AddAddressField("to", to)
event.AddUint64Field("amount", amount)
framework.EmitEvent(event)
```

---

## 🔧 类型定义

### Address

地址类型（20字节）。

```go
type Address []byte

func (a Address) ToBytes() []byte
func (a Address) ToString() string
```

### TokenID

代币ID类型。

```go
type TokenID string
```

### Amount

数量类型。

```go
type Amount uint64
```

---

## ⚠️ 错误码

| 错误码 | 常量 | 说明 |
|--------|------|------|
| 0 | `SUCCESS` | 成功 |
| 1 | `ERROR_INVALID_PARAMS` | 参数无效 |
| 2 | `ERROR_INSUFFICIENT_BALANCE` | 余额不足 |
| 3 | `ERROR_EXECUTION_FAILED` | 执行失败 |
| 4 | `ERROR_PERMISSION_DENIED` | 权限不足 |

---

## 🎯 最佳实践

### 1. 优先使用 Helpers 层

**推荐**（使用 Helpers 层业务语义接口）:
```go
import "github.com/weisyn/contract-sdk-go/helpers/token"

err := token.Transfer(from, to, tokenID, amount)
if err != nil {
    return framework.ERROR_EXECUTION_FAILED
}
```

**不推荐**（直接使用 Framework 层交易构建）:
```go
// Framework 层的交易构建是内部实现，不应直接使用
// framework.BeginTransaction()...
```

### 2. 环境查询和事件

**推荐**:
```go
caller := framework.GetCaller()
height := framework.GetBlockHeight()

event := framework.NewEvent("Transfer")
event.AddAddressField("from", from)
event.AddAddressField("to", to)
event.AddUint64Field("amount", amount)
framework.EmitEvent(event)
```

### 3. 参数验证

**推荐**:
```go
params := framework.GetContractParams()
toStr := params.ParseJSON("to")
if toStr == "" {
    return framework.ERROR_INVALID_PARAMS
}

addr, err := framework.ParseAddressBase58(toStr)
if err != nil {
    return framework.ERROR_INVALID_PARAMS
}
```

---

## 📝 完整示例

查看示例合约：
- [ERC-20 代币合约](../examples/token/erc20-token/)
- [基础质押合约](../examples/staking/basic-staking/)
- [更多示例](../examples/README.md)

---

**最后更新**: 2025-11-11

> **注意**: 本文档描述的是 Framework 层的 API。**合约开发者应优先使用 Helpers 层的业务语义接口**（如 `token.Transfer()`, `staking.Stake()` 等），详见 [Helpers 层文档](../helpers/README.md)。

