# AMM（自动化做市商）合约示例

**分类**: Advanced DeFi 示例  
**难度**: ⭐⭐⭐⭐ 高级  
**最后更新**: 2025-11-11

---

## 📋 概述

本示例展示如何使用 WES Contract SDK Go 构建 AMM（Automated Market Maker）合约。通过本示例，您可以学习如何使用 `helpers/token` 和 `helpers/market` 模块实现完整的AMM功能，包括添加流动性、移除流动性、代币交换等。

---

## 🎯 核心功能

本示例实现了完整的AMM功能：

| 功能 | 函数 | 说明 |
|------|------|------|
| ✅ **添加流动性** | `AddLiquidity` | 向流动性池添加代币对，获得LP Token |
| ✅ **移除流动性** | `RemoveLiquidity` | 从流动性池移除代币对，销毁LP Token |
| ✅ **代币交换** | `SwapTokens` | 使用恒定乘积公式进行代币交换 |

---

## 🏗️ 架构设计

```mermaid
graph TB
    A[流动性提供者] -->|调用 AddLiquidity| B[合约函数]
    B -->|使用| C[helpers/token]
    C -->|调用| D[framework层]
    D -->|HostABI| E[WES节点]
    
    C -->|自动处理| F[交易构建]
    C -->|自动处理| G[事件发出]
    
    H[交易者] -->|调用 SwapTokens| B
    I[流动性提供者] -->|调用 RemoveLiquidity| B
    
    style C fill:#e1f5ff
    style D fill:#fff4e1
    style E fill:#ffe1f5
```

**架构说明**：
- **合约层**：开发者编写的合约函数
- **Token/Market层**：业务语义API，自动处理交易构建、事件发出
- **Framework层**：HostABI封装，提供基础原语
- **节点层**：WES节点，执行合约并上链

---

## 📚 功能详解

### 1. AddLiquidity - 添加流动性

**功能说明**：向流动性池添加代币对，获得流动性凭证代币（LP Token）。

**参数格式**：
```json
{
  "token_a_id": "TOKEN_A",
  "token_b_id": "TOKEN_B",
  "amount_a": 1000,
  "amount_b": 2000
}
```

**特点**：
- 使用恒定乘积公式（x*y=k）计算LP Token数量
- 首次添加流动性时，LP Token数量 = sqrt(amountA * amountB)
- 后续添加流动性时，LP Token数量按比例计算

**⚠️ 注意**：这是一个简化实现
- 实际应用中需要实现恒定乘积公式计算
- 流动性凭证代币数量计算
- 首次添加流动性的特殊处理

**使用示例**：
```bash
wes contract call --address {contract_addr} \
  --function AddLiquidity \
  --params '{"token_a_id":"TOKEN_A","token_b_id":"TOKEN_B","amount_a":1000,"amount_b":2000}'
```

---

### 2. RemoveLiquidity - 移除流动性

**功能说明**：从流动性池移除代币对，销毁流动性凭证代币。

**参数格式**：
```json
{
  "token_a_id": "TOKEN_A",
  "token_b_id": "TOKEN_B",
  "lp_token_amount": 100
}
```

**特点**：
- 根据LP Token数量计算应返还的代币数量
- 使用恒定乘积公式确保比例正确
- 销毁LP Token

**⚠️ 注意**：这是一个简化实现
- 实际应用中需要实现恒定乘积公式计算
- 应返还代币数量计算
- LP Token销毁

**使用示例**：
```bash
wes contract call --address {contract_addr} \
  --function RemoveLiquidity \
  --params '{"token_a_id":"TOKEN_A","token_b_id":"TOKEN_B","lp_token_amount":100}'
```

---

### 3. SwapTokens - 代币交换

**功能说明**：使用恒定乘积公式（x*y=k）进行代币交换。

**参数格式**：
```json
{
  "token_in_id": "TOKEN_A",
  "token_out_id": "TOKEN_B",
  "amount_in": 1000,
  "min_amount_out": 1800
}
```

**特点**：
- 使用恒定乘积公式（x*y=k）计算交换价格
- 滑点保护机制（确保输出数量 >= min_amount_out）
- 手续费分成（给流动性提供者）

**⚠️ 注意**：这是一个简化实现
- 实际应用中需要实现恒定乘积公式计算
- 滑点保护机制
- 手续费分成（给流动性提供者）

**使用示例**：
```bash
wes contract call --address {contract_addr} \
  --function SwapTokens \
  --params '{"token_in_id":"TOKEN_A","token_out_id":"TOKEN_B","amount_in":1000,"min_amount_out":1800}'
```

---

## 🚀 快速开始

### 1. 编译合约

```bash
cd advanced/defi/amm
bash build.sh
```

编译完成后会生成 `main.wasm` 文件。

### 2. 部署合约

```bash
# 使用 WES CLI 部署
wes contract deploy --wasm main.wasm
```

### 3. 调用合约

```bash
# 添加流动性
wes contract call --address {contract_addr} \
  --function AddLiquidity \
  --params '{"token_a_id":"TOKEN_A","token_b_id":"TOKEN_B","amount_a":1000,"amount_b":2000}'
```

---

## 📊 SDK vs 应用层职责

| 职责 | SDK 提供 | 应用层实现 |
|------|---------|-----------|
| **代币转移** | ✅ 自动处理 | - |
| **交易构建** | ✅ 自动处理 | - |
| **事件发出** | ✅ 自动处理 | - |
| **恒定乘积公式** | ❌ | ✅ 需要实现（x*y=k） |
| **滑点保护** | ❌ | ✅ 需要实现 |
| **手续费分成** | ❌ | ✅ 需要实现（给流动性提供者） |
| **流动性凭证代币管理** | ❌ | ✅ 需要实现（铸造、销毁、交易） |

---

## 💡 设计理念

### AMM的特点

- ✅ **自动化做市**：无需订单簿，自动匹配交易
- ✅ **价格发现**：使用恒定乘积公式自动发现价格
- ✅ **流动性提供**：流动性提供者获得交易手续费分成
- ✅ **去中心化**：无需中心化交易所

### SDK 提供"积木"

SDK 提供基础能力（Transfer、Mint、Burn），开发者可以：

- ✅ 直接使用基础功能创建AMM应用
- ✅ 添加业务规则实现定制需求
- ✅ 组合多个功能实现复杂场景

### 应用层搭建"建筑"

应用层在 SDK 基础上实现：

- ✅ 恒定乘积公式（x*y=k）价格计算
- ✅ 滑点保护机制
- ✅ 手续费分成（给流动性提供者）
- ✅ 流动性凭证代币管理（铸造、销毁、交易）

---

## 🔗 相关文档

- [Token 模块文档](../../../helpers/token/README.md) - Token 模块详细说明
- [Market 模块文档](../../../helpers/market/README.md) - Market 模块详细说明
- [Framework 文档](../../../framework/README.md) - Framework 层说明
- [示例总览](../../README.md) - 所有示例索引
- [示例总览](../../README.md) - 示例组织结构规划

---

**最后更新**: 2025-11-11

