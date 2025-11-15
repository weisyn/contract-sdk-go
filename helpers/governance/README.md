# Governance 业务语义模块

**版本**: 1.0  
**状态**: ✅ 稳定  
**最后更新**: 2025-11-11

---

## 📋 概述

Governance 模块提供治理相关的业务语义API，包括创建提案、投票等功能。

---

## 🎯 核心功能

### 1. Propose - 创建提案

**功能**: 创建治理提案

**签名**:
```go
func Propose(proposer framework.Address, proposalID []byte, proposalData []byte) error
```

**示例**:
```go
proposalID := []byte("proposal_123")
proposalData := []byte("proposal content")
err := governance.Propose(caller, proposalID, proposalData)
```

**输入输出组合模式**:
- `StateOutput` - 记录提案状态

---

### 2. Vote - 投票

**功能**: 对提案进行投票

**签名**:
```go
func Vote(voter framework.Address, proposalID []byte, support bool) error
```

**示例**:
```go
err := governance.Vote(caller, proposalID, true)  // 支持
```

**输入输出组合模式**:
- `StateOutput` - 记录投票状态

---

### 3. VoteAndCount - 投票并统计

**功能**：投票并自动统计，判断是否通过阈值

**签名**：
```go
func VoteAndCount(
    voter framework.Address,
    proposalID []byte,
    support bool,
    threshold uint64,
) (*VoteAndCountResult, error)
```

**返回结果**：
```go
type VoteAndCountResult struct {
    TotalVotes   uint64 // 总票数
    SupportVotes uint64 // 支持票数
    OpposeVotes  uint64 // 反对票数
    Passed       bool   // 是否通过（基于阈值判断）
    Threshold    uint64 // 通过阈值
}
```

**示例**：
```go
result, err := governance.VoteAndCount(caller, proposalID, true, 1000)
if err != nil {
    return framework.ERROR_EXECUTION_FAILED
}

// result.Passed 表示是否通过阈值
```

---

## 💡 使用示例

### 完整示例：治理合约

```go
package main

import (
    "github.com/weisyn/contract-sdk-go/helpers/governance"
    "github.com/weisyn/contract-sdk-go/framework"
)

//export Propose
func Propose() uint32 {
    params := framework.GetContractParams()
    proposalID := []byte(params.ParseJSON("proposal_id"))
    proposalData := []byte(params.ParseJSON("proposal_data"))
    
    caller := framework.GetCaller()
    err := governance.Propose(caller, proposalID, proposalData)
    if err != nil {
        return framework.ERROR_EXECUTION_FAILED
    }
    
    return framework.SUCCESS
}

//export Vote
func Vote() uint32 {
    params := framework.GetContractParams()
    proposalID := []byte(params.ParseJSON("proposal_id"))
    support := params.ParseJSONBool("support")
    
    caller := framework.GetCaller()
    err := governance.Vote(caller, proposalID, support)
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

**文档状态**: 开发中  
**下一步**: 完善实现细节

