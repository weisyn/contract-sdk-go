# Market 业务语义模块

**版本**: 1.0  
**状态**: ✅ 稳定  
**最后更新**: 2025-11-11

---

## 📋 概述

Market 模块提供市场相关的业务语义API，包括托管、分阶段释放等功能。

**注意**: 本模块仅提供原子操作（Escrow、Release），不包含组合场景（如Swap、Liquidity等）。

---

## 🎯 核心功能

### 1. Escrow - 托管

**功能**: 创建代币托管

**签名**:
```go
func Escrow(buyer, seller framework.Address, tokenID framework.TokenID, amount framework.Amount, escrowID []byte) error
```

**示例**:
```go
escrowID := []byte("escrow_123")
err := market.Escrow(buyer, seller, nil, framework.Amount(10000), escrowID)
```

**输入输出组合模式**:
- `N inputs + M outputs + ContractLock` - 将代币转移到托管地址
- `StateOutput` - 记录托管状态

---

### 2. Release - 分阶段释放

**功能**: 创建分阶段释放计划

**签名**:
```go
func Release(from, beneficiary framework.Address, tokenID framework.TokenID, totalAmount framework.Amount, vestingID []byte) error
```

**示例**:
```go
vestingID := []byte("vesting_123")
err := market.Release(caller, beneficiary, nil, framework.Amount(100000), vestingID)
```

**输入输出组合模式**:
- `N inputs + M outputs + TimeLock/HeightLock` - 将代币转移到受益人地址
- `StateOutput` - 记录释放计划状态

---

## ⚠️ 不包含的功能

以下功能属于组合场景，**不应由SDK提供**，应由应用层实现：

- ❌ Swap - 交换（需要价格计算、滑点保护等业务逻辑）
- ❌ AddLiquidity - 添加流动性（需要份额计算、比例管理等业务逻辑）
- ❌ RemoveLiquidity - 移除流动性（需要份额计算、比例管理等业务逻辑）

---

## 💡 使用示例

### 完整示例：市场合约

```go
package main

import (
    "github.com/weisyn/contract-sdk-go/helpers/market"
    "github.com/weisyn/contract-sdk-go/framework"
)

//export Escrow
func Escrow() uint32 {
    params := framework.GetContractParams()
    sellerStr := params.ParseJSON("seller")
    amount := params.ParseJSONInt("amount")
    escrowID := []byte(params.ParseJSON("escrow_id"))
    
    seller, err := framework.ParseAddressBase58(sellerStr)
    if err != nil {
        return framework.ERROR_INVALID_PARAMS
    }
    
    buyer := framework.GetCaller()
    err = market.Escrow(buyer, seller, nil, framework.Amount(amount), escrowID)
    if err != nil {
        return framework.ERROR_EXECUTION_FAILED
    }
    
    return framework.SUCCESS
}

//export Release
func Release() uint32 {
    params := framework.GetContractParams()
    beneficiaryStr := params.ParseJSON("beneficiary")
    totalAmount := params.ParseJSONInt("total_amount")
    vestingID := []byte(params.ParseJSON("vesting_id"))
    
    beneficiary, err := framework.ParseAddressBase58(beneficiaryStr)
    if err != nil {
        return framework.ERROR_INVALID_PARAMS
    }
    
    caller := framework.GetCaller()
    err = market.Release(caller, beneficiary, nil, framework.Amount(totalAmount), vestingID)
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
- [应用场景分析](../../docs/APPLICATION_SCENARIOS_ANALYSIS.md)

---

**最后更新**: 2025-11-11

