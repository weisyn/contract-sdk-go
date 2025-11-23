# Helpers - 业务语义层（统一创新整体）

**版本**: 1.0  
**状态**: ✅ 稳定  
**最后更新**: 2025-11-11

---

## 📋 概述

Helpers 层提供统一的业务语义API，**既具备传统能力，也具有ISPC创新力**。

**核心设计理念**：
|- ✅ **统一整体**：所有模块按业务领域分类，统一对外展示
|- ✅ **渐进增强**：ISPC能力增强传统功能，而非替代
|- ✅ **业务优先**：按业务领域分类，而非技术分类
|- ✅ **ISPC赋能**：ISPC能力作为内部能力，增强业务模块

**对外展示**：开发者看到的是一个统一的创新SDK，既具备传统能力，也具有ISPC创新力。

---

## 🎯 核心职责

1. **业务语义封装**：提供 Transfer、Mint、Stake 等业务操作
2. **ISPC能力增强**：通过ISPC能力增强传统功能（受控外部交互、业务执行等）
3. **基于 Framework**：所有操作基于 framework 层构建
4. **类型安全**：提供类型安全的业务 API
5. **错误处理**：业务级错误处理

---

## 📐 架构定位

```
L3: helpers/  ← 业务语义层（本层）
         ↓ 基于
L2: framework/         ← 框架层（HostABI封装）
         ↓ 调用
L1: HostABI原语        ← 由ISPC提供（17个原语）
```

### ISPC 的定位与体现

**ISPC 是底层执行范式，不是 SDK 中的一个独立层**：
- ISPC 通过 HostABI 原语体现（包括受控外部交互的原语）
- Framework 层封装 HostABI，提供类型安全的 API
- Helpers 层使用 Framework 能力，按照 ISPC 的思路组织业务逻辑
- **不需要独立的"ISPC 层"**，ISPC 的执行范式通过业务逻辑的组织方式体现

---

## 🏗️ 模块列表（统一业务领域分类）

### 1. Token 模块 ✅

**路径**: `helpers/token/`

**功能**:
|- ✅ **基础功能**：Transfer, Mint, Burn, Approve, Freeze, Airdrop（传统：状态操作）
|- 🌟 **ISPC增强**：SmartMint（利用external验证，未来扩展）

**状态**: 开发中

---

### 2. Staking 模块 ✅

**路径**: `helpers/staking/`

**功能**:
|- ✅ **基础功能**：Stake, Unstake, Delegate, Undelegate（传统：状态操作）
|- 🌟 **ISPC增强**：SmartStake（利用execution执行业务逻辑，未来扩展）

**状态**: 开发中

---

### 3. Market 模块 ✅

**路径**: `helpers/market/`

**功能**:
|- ✅ **基础功能**：Escrow, Release（传统：状态操作）
|- 🌟 **ISPC增强**：SmartEscrow（利用external验证，未来扩展）

**状态**: 开发中

**注意**: 仅提供原子操作，不包含组合场景（Swap、Liquidity等）

---

### 4. Governance 模块 ✅

**路径**: `helpers/governance/`

**功能**:
|- ✅ **基础功能**：Propose, Vote（传统：状态操作）
|- 🌟 **ISPC增强**：VoteAndCount（业务执行：投票并统计，自动判断阈值）

**状态**: 开发中

---

### 5. RWA 模块 🌟

**路径**: `helpers/rwa/`

**功能**:
|- 🌟 **充分利用ISPC**：ValidateAndTokenize（受控外部交互 + 业务执行）
|- 🌟 ValidateAsset（利用external验证资产）
|- 🌟 ValueAsset（利用external估值资产）

**特点**: 新模块充分利用ISPC能力，无需传统预言机

**状态**: 开发中

---

### 6. NFT 模块 ✅

**路径**: `helpers/nft/`

**功能**:
|- ✅ **基础功能**：Mint, Transfer, Burn, OwnerOf, BalanceOf, GetMetadata（NFT操作）
|- ✅ **基于Token模块**：NFT本质上是数量为1的特殊代币

**特点**: NFT模块是对Token模块的扩展，专门用于处理NFT业务场景

**状态**: 开发中

---

### 7. External 模块（受控外部交互）

**路径**: `helpers/external/`

**说明**: 封装 ISPC 受控外部交互的业务语义，供其他业务模块使用

**功能**:
|- 🌟 ValidateAndQuery - 验证并查询外部状态
|- 🌟 CallAPI - 调用外部 API
|- 🌟 QueryDatabase - 查询数据库

**本质定位**:
- `helpers/external/` 是对 Framework 层 HostABI 封装的**业务语义封装**
- 它调用 `framework.DeclareExternalState()` 等函数
- 与其他 helpers 模块（token, staking 等）一样，都是业务语义层的一部分
- **区别**：它封装的是 ISPC 创新的能力（受控外部交互）

**状态**: 开发中

---

### 7. Resource 模块 🚧

**路径**: `helpers/resource/`

**功能**:
|- 🚧 Deploy - 部署资源
|- 🚧 Query - 查询资源

**状态**: 待实现

---

## 🌟 ISPC创新体现

### ISPC 的定位

**ISPC 是底层执行范式，不是 SDK 中的一个独立层**：
- ISPC 通过 HostABI 原语体现（包括受控外部交互的原语）
- Framework 层封装 HostABI，提供类型安全的 API
- Helpers 层使用 Framework 能力，按照 ISPC 的思路组织业务逻辑
- **不需要独立的"ISPC 层"**

### ISPC能力增强模式

**核心理念**：ISPC 能力增强传统功能，而非替代

**增强方式**：
1. **受控外部交互**：通过 `helpers/external` 封装 ISPC 受控外部交互的业务语义
2. **业务执行**：通过业务逻辑的组织方式体现 ISPC 执行范式
3. **渐进增强**：在现有模块中增加 ISPC 增强功能（如 `VoteAndCount`）

**示例**：
|- `governance.Vote()` - 传统：记录投票状态
|- `governance.VoteAndCount()` - ISPC 增强：投票并统计，执行业务逻辑
|- `rwa.ValidateAndTokenize()` - ISPC 创新：充分利用 ISPC 能力（受控外部交互）

### 统一对外接口

**开发者体验**：
|- ✅ 按业务领域选择模块（token, staking, governance, rwa, external 等）
|- ✅ 不需要区分 ISPC 和传统功能
|- ✅ 通过功能名称了解是否支持 ISPC 增强（如 `VoteAndCount`）
|- ✅ ISPC 的执行范式通过业务逻辑的组织方式体现

---

## 💡 设计原则

### 1. 统一整体，不割裂

✅ **正确**：
|- 所有模块统一在业务领域分类下
|- ISPC能力作为内部能力，增强业务模块
|- 对外展示统一的创新SDK

❌ **错误**：
|- 区分"ISPC模块"和"传统模块"
|- 割裂ISPC和传统功能

### 2. 渐进增强

✅ **正确**：
|- 基础功能（传统能力）+ ISPC增强功能
|- 新模块充分利用ISPC能力
|- 逐步演进，不强制迁移

### 3. 不直接调用 HostABI

✅ **正确**：
```go
// 使用framework层
balance := framework.QueryUTXOBalance(address, tokenID)
```

❌ **错误**：
```go
// 直接调用HostABI原语
hostabi.UTXOLookup(...)
```

### 4. ISPC 的正确理解

✅ **正确**：
|- ISPC 是底层执行范式，通过 HostABI 体现
|- Framework 层封装 HostABI（包括受控外部交互的 HostABI 函数）
|- Helpers 层使用 Framework 能力，按照 ISPC 的思路组织业务逻辑

❌ **错误**：
|- ISPC 是 SDK 中的一个独立层
|- 需要创建 `ispc/` 目录作为独立层

---

## 🚀 使用示例

### 基础功能（传统能力）

```go
import (
	"github.com/weisyn/contract-sdk-go/helpers/token"
	"github.com/weisyn/contract-sdk-go/framework"
)

func Transfer() uint32 {
	caller := framework.GetCaller()
	recipient := framework.Address{...}
	amount := uint64(100)
	tokenID := framework.TokenID{...}
	
	err := token.Transfer(caller, recipient, amount, tokenID)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}
	
	return framework.SUCCESS
}
```

### ISPC 增强功能（受控外部交互）

```go
import (
	"github.com/weisyn/contract-sdk-go/helpers/rwa"
	"github.com/weisyn/contract-sdk-go/framework"
)

func TokenizeAsset() uint32 {
	assetID := "ASSET_001"
	documents := []byte("{...}")
	
	// 使用 ISPC 受控外部交互，验证并代币化资产
	result, err := rwa.ValidateAndTokenize(
		assetID,
		documents,
		"https://validator.example.com/api/validate",
		&framework.Evidence{...},
		"https://valuation.example.com/api/value",
		&framework.Evidence{...},
	)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}
	
	// 直接获得业务结果，无需知道 Transaction 的存在
	return framework.SUCCESS
}
```

---

## 📚 相关文档

- [Contract SDK 主 README](../README.md)
- [Framework 层 README](../framework/README.md)
- [ISPC 创新分析](../docs/ISPC_INNOVATION_ANALYSIS.md)
- [整体结构设计](../docs/STRUCTURE_DESIGN.md)

---

**最后更新**: 2025-11-11
