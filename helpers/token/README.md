# Token 业务语义模块

**版本**: 1.0  
**状态**: ✅ 稳定  
**最后更新**: 2025-11-11

---

## 📋 概述

Token 模块提供代币相关的业务语义API，包括转账、铸造、销毁、授权、冻结、空投等功能。

---

## 🎯 核心功能

### 1. Transfer - 转账

**功能**: 从指定地址转账到目标地址

**签名**:
```go
func Transfer(from, to framework.Address, tokenID framework.TokenID, amount framework.Amount) error
```

**示例**:
```go
err := token.Transfer(caller, recipient, nil, framework.Amount(1000))
```

---

### 2. Mint - 铸造

**功能**: 铸造新代币到指定地址

**签名**:
```go
func Mint(to framework.Address, tokenID framework.TokenID, amount framework.Amount) error
```

**示例**:
```go
err := token.Mint(recipient, framework.TokenID("my_token"), framework.Amount(1000))
```

---

### 3. Burn - 销毁

**功能**: 销毁指定地址的代币

**签名**:
```go
func Burn(from framework.Address, tokenID framework.TokenID, amount framework.Amount) error
```

**示例**:
```go
err := token.Burn(caller, framework.TokenID("my_token"), framework.Amount(500))
```

---

### 4. Approve - 授权

**功能**: 授权指定地址使用代币

**签名**:
```go
func Approve(owner, spender framework.Address, tokenID framework.TokenID, amount framework.Amount) error
```

**示例**:
```go
err := token.Approve(caller, spender, framework.TokenID("my_token"), framework.Amount(1000))
```

---

### 5. Freeze - 冻结

**功能**: 冻结指定地址的代币

**签名**:
```go
func Freeze(target framework.Address, tokenID framework.TokenID, amount framework.Amount) error
```

**示例**:
```go
err := token.Freeze(target, framework.TokenID("my_token"), framework.Amount(1000))
```

---

### 6. Airdrop - 空投

**功能**: 批量转账到多个地址（需要发送者有足够余额）

**签名**:
```go
func Airdrop(from framework.Address, recipients []AirdropRecipient, tokenID framework.TokenID) error
```

**示例**:
```go
recipients := []token.AirdropRecipient{
    {Address: addr1, Amount: framework.Amount(100)},
    {Address: addr2, Amount: framework.Amount(200)},
}
err := token.Airdrop(caller, recipients, framework.TokenID("my_token"))
```

---

### 7. BatchMint - 批量铸造

**功能**: 一次性向多个地址铸造代币（不需要发送者有余额）

**签名**:
```go
func BatchMint(recipients []MintRecipient, tokenID framework.TokenID) error
```

**示例**:
```go
recipients := []token.MintRecipient{
    {Address: addr1, Amount: framework.Amount(100)},
    {Address: addr2, Amount: framework.Amount(200)},
    {Address: addr3, Amount: framework.Amount(300)},
}
err := token.BatchMint(recipients, framework.TokenID("my_token"))
```

**注意**:
- 批量铸造会在一次交易中创建多个AssetOutput（UTXO输出）
- 适用于初始分配、空投等场景
- 与Airdrop的区别：BatchMint不需要发送者有余额，直接从合约铸造

---

## 💡 使用示例

### 完整示例：代币合约

```go
package main

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

//export Mint
func Mint() uint32 {
    params := framework.GetContractParams()
    toStr := params.ParseJSON("to")
    amount := params.ParseJSONInt("amount")
    
    to, err := framework.ParseAddressBase58(toStr)
    if err != nil {
        return framework.ERROR_INVALID_PARAMS
    }
    
    err = token.Mint(to, framework.TokenID("my_token"), framework.Amount(amount))
    if err != nil {
        return framework.ERROR_EXECUTION_FAILED
    }
    
    return framework.SUCCESS
}
```

---

## 🔗 相关文档

- [Contract Helpers总览](../README.md)
- [Framework层文档](../../framework/README.md)

---

**最后更新**: 2025-11-11

