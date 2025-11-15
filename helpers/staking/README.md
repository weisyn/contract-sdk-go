# Staking 业务语义模块

**版本**: 1.0  
**状态**: ✅ 稳定  
**最后更新**: 2025-11-11

---

## 📋 概述

Staking 模块提供质押和委托相关的业务语义API，包括质押、解质押、委托、取消委托等功能。

---

## 🎯 核心功能

### 1. Stake - 质押

**功能**: 将代币质押给验证者

**签名**:
```go
func Stake(staker, validator framework.Address, tokenID framework.TokenID, amount framework.Amount) error
```

**示例**:
```go
err := staking.Stake(caller, validatorAddr, nil, framework.Amount(10000))
```

**输入输出组合模式**:
- `N inputs + M outputs + ContractLock`
- 将代币转移到验证者地址，并添加ContractLock锁定条件

---

### 2. Unstake - 解质押

**功能**: 解除质押，取回代币

**签名**:
```go
func Unstake(staker, validator framework.Address, tokenID framework.TokenID, amount framework.Amount) error
```

**示例**:
```go
err := staking.Unstake(caller, validatorAddr, nil, framework.Amount(5000))
```

**输入输出组合模式**:
- `ContractLock UTXO + 解锁`
- 从验证者地址转回质押者，解锁ContractLock

---

### 3. Delegate - 委托

**功能**: 将代币委托给验证者

**签名**:
```go
func Delegate(delegator, validator framework.Address, tokenID framework.TokenID, amount framework.Amount) error
```

**示例**:
```go
err := staking.Delegate(caller, validatorAddr, nil, framework.Amount(5000))
```

**输入输出组合模式**:
- `N inputs + M outputs + DelegationLock`
- 将代币转移到验证者地址，并添加DelegationLock锁定条件

---

### 4. Undelegate - 取消委托

**功能**: 取消委托，取回代币

**签名**:
```go
func Undelegate(delegator, validator framework.Address, tokenID framework.TokenID, amount framework.Amount) error
```

**示例**:
```go
err := staking.Undelegate(caller, validatorAddr, nil, framework.Amount(2000))
```

**输入输出组合模式**:
- `DelegationLock UTXO + 解锁`
- 从验证者地址转回委托者，解锁DelegationLock

---

## 💡 使用示例

### 完整示例：质押合约

```go
package main

import (
    "github.com/weisyn/contract-sdk-go/helpers/staking"
    "github.com/weisyn/contract-sdk-go/framework"
)

//export Stake
func Stake() uint32 {
    params := framework.GetContractParams()
    validatorStr := params.ParseJSON("validator")
    amount := params.ParseJSONInt("amount")
    
    validator, err := framework.ParseAddressBase58(validatorStr)
    if err != nil {
        return framework.ERROR_INVALID_PARAMS
    }
    
    caller := framework.GetCaller()
    err = staking.Stake(caller, validator, nil, framework.Amount(amount))
    if err != nil {
        return framework.ERROR_EXECUTION_FAILED
    }
    
    return framework.SUCCESS
}

//export Unstake
func Unstake() uint32 {
    params := framework.GetContractParams()
    validatorStr := params.ParseJSON("validator")
    amount := params.ParseJSONInt("amount")
    
    validator, err := framework.ParseAddressBase58(validatorStr)
    if err != nil {
        return framework.ERROR_INVALID_PARAMS
    }
    
    caller := framework.GetCaller()
    err = staking.Unstake(caller, validator, nil, framework.Amount(amount))
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

