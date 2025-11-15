# 域名服务NFT合约示例

**分类**: NFT 示例  
**难度**: ⭐⭐⭐ 高级  
**最后更新**: 2025-11-11

---

## 📋 概述

本示例展示如何使用 WES Contract SDK Go 构建域名服务NFT合约。通过本示例，您可以学习如何使用 `helpers/token` 模块创建和管理域名NFT，实现去中心化域名服务。

---

## 🎯 核心功能

本示例实现了完整的域名服务NFT功能：

| 功能 | 函数 | 说明 |
|------|------|------|
| ✅ **注册域名** | `MintDomain` | 注册唯一的域名NFT |
| ✅ **转移域名** | `TransferDomain` | 转移域名所有权 |
| ✅ **查询域名** | `QueryDomain` | 查询域名信息和所有者 |

---

## 🏗️ 架构设计

```mermaid
graph TB
    A[域名注册者] -->|调用 MintDomain| B[合约函数]
    B -->|使用| C[helpers/token]
    C -->|调用| D[framework层]
    D -->|HostABI| E[WES节点]
    
    C -->|自动处理| F[交易构建]
    C -->|自动处理| G[事件发出]
    
    H[域名持有者] -->|调用 TransferDomain| B
    I[查询者] -->|调用 QueryDomain| B
    
    style C fill:#e1f5ff
    style D fill:#fff4e1
    style E fill:#ffe1f5
```

**架构说明**：
- **合约层**：开发者编写的合约函数
- **Token层**：业务语义API，自动处理交易构建、事件发出
- **Framework层**：HostABI封装，提供基础原语
- **节点层**：WES节点，执行合约并上链

---

## 📚 功能详解

### 1. MintDomain - 注册域名

**功能说明**：使用 `token.Mint()` 注册唯一的域名NFT。

**参数格式**：
```json
{
  "to": "Df2Lft7toFVfjlKKhsBtLQOQsQbQeRnTn",
  "token_id": "DOMAIN_001",
  "domain_name": "example.wes",
  "registrar": "WES Registry"
}
```

**特点**：
- 每个域名都有唯一的tokenID
- 域名包含元数据（名称、注册机构等）
- 域名不可分割，转移时数量为1

**使用示例**：
```bash
wes contract call --address {contract_addr} \
  --function MintDomain \
  --params '{"to":"Df2Lft7toFVfjlKKhsBtLQOQsQbQeRnTn","token_id":"DOMAIN_001","domain_name":"example.wes","registrar":"WES Registry"}'
```

---

### 2. TransferDomain - 转移域名

**功能说明**：使用 `token.Transfer()` 转移域名所有权。

**参数格式**：
```json
{
  "to": "Cf1Kes6snEUeykiJJgrAtKPNPrAzPdPmSn",
  "token_id": "DOMAIN_001"
}
```

**使用示例**：
```bash
wes contract call --address {contract_addr} \
  --function TransferDomain \
  --params '{"to":"Cf1Kes6snEUeykiJJgrAtKPNPrAzPdPmSn","token_id":"DOMAIN_001"}'
```

---

### 3. QueryDomain - 查询域名

**功能说明**：查询域名的详细信息和所有者。

**参数格式**：
```json
{
  "token_id": "DOMAIN_001"
}
```

**使用示例**：
```bash
wes contract call --address {contract_addr} \
  --function QueryDomain \
  --params '{"token_id":"DOMAIN_001"}'
```

---

## 🚀 快速开始

### 1. 编译合约

```bash
cd nft/domains
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
# 注册域名
wes contract call --address {contract_addr} \
  --function MintDomain \
  --params '{"to":"Df2Lft7toFVfjlKKhsBtLQOQsQbQeRnTn","token_id":"DOMAIN_001","domain_name":"example.wes","registrar":"WES Registry"}'
```

---

## 📊 SDK vs 应用层职责

| 职责 | SDK 提供 | 应用层实现 |
|------|---------|-----------|
| **NFT铸造** | ✅ 自动处理 | - |
| **NFT转移** | ✅ 自动处理 | - |
| **交易构建** | ✅ 自动处理 | - |
| **事件发出** | ✅ 自动处理 | - |
| **域名验证** | ❌ | ✅ 需要实现（验证域名格式） |
| **域名续费** | ❌ | ✅ 需要实现（续费机制） |

---

## 💡 设计理念

### 域名服务NFT的应用场景

- ✅ **去中心化域名**：去中心化域名服务
- ✅ **数字身份**：域名作为数字身份标识
- ✅ **品牌保护**：企业品牌域名保护
- ✅ **个人域名**：个人域名服务

### SDK 提供"积木"

SDK 提供基础能力（Mint、Transfer），开发者可以：

- ✅ 直接使用基础功能创建域名服务NFT应用
- ✅ 添加业务规则实现定制需求
- ✅ 组合多个功能实现复杂场景

---

## 🔗 相关文档

- [Token 模块文档](../../helpers/token/README.md) - Token 模块详细说明
- [Framework 文档](../../framework/README.md) - Framework 层说明
- [示例总览](../README.md) - 所有示例索引
- [示例总览](../README.md) - 示例组织结构规划

---

**最后更新**: 2025-11-11
