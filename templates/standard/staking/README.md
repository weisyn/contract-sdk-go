# Staking 质押示例

**分类**: 质押相关示例  
**难度**: ⭐⭐ - ⭐⭐⭐  
**最后更新**: 2025-11-11

---

## 📋 概述

Staking 质押示例展示如何使用 WES Contract SDK Go 构建质押和委托相关的智能合约。

---

## 📚 示例列表

本目录包含以下质押示例：

### 1. [Basic Staking](basic-staking/) ✅

**难度**: ⭐⭐⭐ 高级  
**功能**: 使用helpers的基础质押合约

**学习点**:
- helpers/staking模块使用
- 质押和委托业务语义API
- Staking模块所有功能：Stake, Unstake, Delegate, Undelegate

**功能列表**:
- ✅ Initialize - 初始化质押合约
- ✅ Stake - 质押
- ✅ Unstake - 解质押
- ✅ Delegate - 委托
- ✅ Undelegate - 取消委托

**适用场景**:
- 🔒 网络安全质押
- 💰 收益获取
- 🏛️ 治理参与

---

### 2. [Delegation](delegation/) ✅

**难度**: ⭐⭐⭐ 高级  
**功能**: 委托质押

**功能列表**:
- ✅ Initialize - 初始化委托合约
- ✅ Delegate - 委托给验证者
- ✅ Undelegate - 取消委托
- ✅ QueryDelegation - 查询委托信息

**适用场景**:
- 🤝 委托管理
- 💰 委托收益分配
- 🏛️ 验证者委托

---

## 🎯 使用场景

- 🔒 **网络安全**：质押代币参与网络共识
- 💰 **收益获取**：通过质押获得收益
- 🏛️ **治理参与**：通过质押获得投票权
- 🤝 **委托管理**：将质押权委托给验证者

---

## 🔗 相关文档

- [Staking 模块文档](../../helpers/staking/README.md) - Staking 模块详细说明
- [示例总览](../README.md) - 所有示例索引
- [示例总览](../README.md) - 示例组织结构规划

---

**最后更新**: 2025-11-11

