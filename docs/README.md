# Contract SDK Go - 文档中心

**版本**: v1.0.0  
**最后更新**: 2025-01-23

---

<div align="center">

## 📚 文档导航中心

**按角色快速找到你需要的文档**

[👨‍💻 合约开发者](#-合约开发者) • [🏗️ 架构师/贡献者](#️-架构师贡献者) • [📖 参考文档](#-参考文档)

</div>

---

## 📋 文档定位说明

### 文档层次关系

```
主 README.md（用户入口）
    │
    ├─> 30秒上手、核心能力、架构概览
    │
    └─> docs/README.md（文档中心）← 你现在在这里
            │
            ├─> 开发者指南、API 参考、业务场景
            ├─> 架构设计文档
            └─> 模块文档、示例代码
```

**主 README.md** 的定位：
- ✅ **用户友好的入口**：快速了解 SDK，30秒上手
- ✅ **核心价值展示**：突出 SDK 的核心能力和优势
- ✅ **简洁的架构说明**：分层架构的概览
- ✅ **文档导航入口**：指向文档中心（本文件）

**docs/README.md**（本文件）的定位：
- ✅ **文档索引中心**：列出所有文档及其定位
- ✅ **按角色导航**：为不同角色提供快速导航路径
- ✅ **文档分类**：按用户文档、设计文档、参考文档分类
- ✅ **使用建议**：为不同场景提供文档使用建议

---

## 👨‍💻 合约开发者

### 🚀 快速开始路径

```
1. [主 README](../README.md)
   └─> 了解 SDK 是什么，30秒上手
   
2. [Hello World 示例](../examples/basics/hello-world/)
   └─> 完成第一个合约
   
3. [开发者指南](./DEVELOPER_GUIDE.md)
   └─> 深入学习核心概念
   
4. [示例代码](../examples/README.md)
   └─> 参考实际应用示例
```

### 📖 核心文档

#### 必读文档（P0）

- ⭐ **[主 README](../README.md)** - SDK 总览和 30 秒上手
  - SDK 简介和核心价值
  - 安装和第一个合约
  - 核心功能概览

- 📖 **[开发者指南](./DEVELOPER_GUIDE.md)** - 如何使用 SDK 开发合约
  - 快速开始
  - 核心概念
  - 常见场景
  - 最佳实践

- 📚 **[API 参考](./API_REFERENCE.md)** - SDK 接口详细说明
  - Framework 层 API
  - Helpers 层 API
  - 接口参数和返回值
  - 使用示例

#### 推荐文档（P1）

- 🎯 **[业务场景实现指南](./BUSINESS_SCENARIOS.md)** - 如何使用 SDK 实现业务场景
  - 电商场景：订单、支付、库存管理
  - 制造业场景：工单、生产、质检
  - SDK 提供的便捷操作

- 🔧 **[Helpers 层文档](../helpers/README.md)** - 业务语义层详细说明
  - Token、Staking、Governance、Market、RWA、External 模块
  - 各模块 API 和使用示例

- 💡 **[示例代码](../examples/README.md)** - 示例代码索引和指南
  - 基础示例（Hello World、简单代币）
  - 业务示例（Token、NFT、RWA、Staking、Governance、Market）
  - 高级示例（DeFi、AMM、借贷）

#### 可选文档（P2）

- 🏗️ **[Framework 层文档](../framework/README.md)** - 框架层详细说明
  - HostABI 封装
  - 环境查询、事件日志
  - 交易构建（内部实现）

---

## 🏗️ 架构师/贡献者

### 架构设计文档

#### 核心设计文档（P1）

- 🏗️ **[整体结构设计](./STRUCTURE_DESIGN.md)** - SDK 整体架构设计讨论
  - 模块组织方式
  - 依赖关系说明
  - 设计决策记录

- 📊 **[应用场景分析](./APPLICATION_SCENARIOS_ANALYSIS.md)** - 应用场景与 SDK 职责边界分析
  - DeFi、NFT、游戏、治理等场景分析
  - SDK 职责边界划分
  - 设计原则说明

- 🔮 **[ISPC 创新分析](./ISPC_INNOVATION_ANALYSIS.md)** - ISPC 在 SDK 中的体现
  - 受控外部交互机制
  - 业务执行即上链
  - 与传统区块链的对比

#### 规划文档（P2）

- 📈 **[架构规划](./ARCHITECTURE_PLAN.md)** - SDK 架构规划文档
  - 未来演进方向
  - 技术债务管理

- 📐 **[场景可视化指南](./SCENARIOS_VISUAL_GUIDE.md)** - 业务场景的可视化说明
  - 流程图和架构图
  - 场景实现示例

---

## 📖 参考文档

### 规范文档

- **HostABI 规范**: [HostABI原语能力](../../docs/components/core/ispc/capabilities/hostabi-primitives.md)
- **接口映射**: [pkg/interfaces 到 Host ABI 映射](../../docs/reference/contracts/pkg-interfaces-to-hostabi-mapping.md)
- **工具链版本**: [工具链版本矩阵](../../docs/reference/contracts/toolchain-version-matrix.md)

### 开发指南

- **集成测试**: [合约集成测试指南](../../docs/tutorials/contracts/integration-testing.md)
- **WASM 环境说明**: [WASM 环境说明](../../docs/tutorials/contracts/wasm-environment.md)
- **Unsafe 指针警告**: [Unsafe 指针使用警告](./UNSAFE_POINTER_WARNINGS.md)

---

## 🎯 快速导航路径

### 新手入门路径

```
1. [主 README](../README.md)
   └─> 了解 SDK 是什么，30秒上手
   
2. [Hello World 示例](../examples/basics/hello-world/)
   └─> 完成第一个合约
   
3. [开发者指南](./DEVELOPER_GUIDE.md)
   └─> 深入学习核心概念
   
4. [示例代码](../examples/README.md)
   └─> 参考实际应用示例
```

### 功能开发路径

```
1. [核心功能](../README.md#-核心能力)
   └─> 了解 SDK 提供的功能
   
2. [业务场景实现指南](./BUSINESS_SCENARIOS.md)
   └─> 学习如何实现业务场景
   
3. [Helpers 层文档](../helpers/README.md)
   └─> 查看业务语义接口说明
   
4. [API 参考](./API_REFERENCE.md)
   └─> 查阅详细的 API 文档
```

### 深入理解路径

```
1. [SDK 架构](../README.md#-sdk-架构)
   └─> 理解分层架构
   
2. [整体结构设计](./STRUCTURE_DESIGN.md)
   └─> 了解架构设计决策
   
3. [ISPC 创新分析](./ISPC_INNOVATION_ANALYSIS.md)
   └─> 理解 ISPC 的核心创新
   
4. [HostABI 规范](../../docs/components/core/ispc/capabilities/hostabi-primitives.md)
   └─> 深入底层能力
```

---

## 📋 文档分类

### 用户文档（面向合约开发者）

| 文档 | 说明 | 优先级 |
|------|------|--------|
| ⭐ [主 README](../README.md) | SDK 总览和快速开始 | P0 |
| 📖 [开发者指南](./DEVELOPER_GUIDE.md) | 如何使用 SDK 开发合约 | P0 |
| 📚 [API 参考](./API_REFERENCE.md) | SDK 接口详细说明 | P0 |
| 🎯 [业务场景实现指南](./BUSINESS_SCENARIOS.md) | 如何实现业务场景 | P1 |
| 🔧 [Helpers 层文档](../helpers/README.md) | 业务语义层详细说明 | P1 |
| 🏗️ [Framework 层文档](../framework/README.md) | 框架层详细说明 | P2 |
| 💡 [示例代码](../examples/README.md) | 示例代码索引和指南 | P1 |

### 设计文档（面向架构师和贡献者）

| 文档 | 说明 | 优先级 |
|------|------|--------|
| 🏗️ [整体结构设计](./STRUCTURE_DESIGN.md) | 架构设计讨论 | P1 |
| 📊 [应用场景分析](./APPLICATION_SCENARIOS_ANALYSIS.md) | 场景与职责边界分析 | P1 |
| 📈 [架构规划](./ARCHITECTURE_PLAN.md) | 架构规划文档 | P2 |
| 🔮 [ISPC 创新分析](./ISPC_INNOVATION_ANALYSIS.md) | ISPC 技术深度分析 | P1 |
| 📐 [场景可视化指南](./SCENARIOS_VISUAL_GUIDE.md) | 业务场景可视化说明 | P2 |

### 参考文档（面向高级开发者）

| 文档 | 说明 | 优先级 |
|------|------|--------|
| 📘 [HostABI 规范](../../docs/components/core/ispc/capabilities/hostabi-primitives.md) | HostABI 原语能力 | P2 |
| 🔗 [接口映射](../../docs/reference/contracts/pkg-interfaces-to-hostabi-mapping.md) | pkg/interfaces 到 Host ABI 映射 | P2 |
| 🔧 [工具链版本矩阵](../../docs/reference/contracts/toolchain-version-matrix.md) | 工具链版本要求 | P2 |
| 🧪 [集成测试指南](../../docs/tutorials/contracts/integration-testing.md) | 合约集成测试 | P2 |
| 🌐 [WASM 环境说明](../../docs/tutorials/contracts/wasm-environment.md) | WASM 环境详解 | P2 |
| ⚠️ [Unsafe 指针警告](./UNSAFE_POINTER_WARNINGS.md) | Unsafe 指针使用警告 | P2 |

---

## 💡 文档使用建议

### 如果你是新手

1. **先看主 README**：了解 SDK 是什么，完成 30 秒上手
2. **运行 Hello World**：完成第一个合约，理解基本概念
3. **阅读开发者指南**：深入学习核心概念和最佳实践
4. **参考示例代码**：学习实际应用示例

### 如果你在开发功能

1. **查看核心功能**：了解 SDK 提供的功能
2. **查阅 API 参考**：查找具体的 API 使用方法
3. **参考业务场景指南**：学习如何实现业务场景
4. **查看模块文档**：深入了解特定模块

### 如果你想贡献代码

1. **阅读架构设计文档**：理解 SDK 的整体架构
2. **查看应用场景分析**：理解 SDK 的职责边界
3. **参考架构规划**：了解未来演进方向
4. **阅读 HostABI 规范**：深入理解底层能力

---

## 🔗 外部资源

- [WES 主项目](https://github.com/weisyn/weisyn) - WES 区块链主仓库
- [WES 文档中心](../../../docs/) - 完整技术文档
- [WES 系统架构](../../../docs/system/architecture/) - 系统架构详解
- [WES 主 README](../../../README.md) - WES 项目总览

---

**最后更新**: 2025-01-23  
**维护者**: WES Core Team
