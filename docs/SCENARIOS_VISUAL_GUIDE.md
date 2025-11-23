# 场景与架构可视化指南

**版本**: v1.0.0  
**最后更新**: 2025-01-23

---

## 📋 文档定位

> 📌 **重要说明**：本文档提供 SDK 相关的简明架构/场景图。  
> 详细业务流图请参考主仓库文档。

**本文档目标**：
- 提供 SDK 内部分层架构图
- 提供 SDK 与平台其他组件的交互图
- 帮助快速理解 SDK 的结构和关系

---

## 🏗️ SDK 内部分层架构

### 三层架构图

```mermaid
graph TB
    subgraph HELPERS["L3: Helpers 业务语义层"]
        direction TB
        TOKEN["Token 模块"]
        STAKING["Staking 模块"]
        GOVERNANCE["Governance 模块"]
        MARKET["Market 模块"]
        NFT["NFT 模块"]
        RWA["RWA 模块"]
        EXTERNAL["External 模块"]
    end
    
    subgraph FRAMEWORK["L2: Framework 框架层"]
        direction TB
        CONTRACT["Contract 接口"]
        CONTEXT["Context 环境查询"]
        EVENT["Event 事件日志"]
        STORAGE["Storage 存储"]
        INTERNAL["Internal 层<br/>交易构建等"]
    end
    
    subgraph HOSTABI["L1: HostABI 原语层"]
        direction TB
        H1["17 个最小原语<br/>由 ISPC 提供"]
    end
    
    HELPERS --> FRAMEWORK
    FRAMEWORK --> INTERNAL
    INTERNAL --> HOSTABI
    
    style HELPERS fill:#FFD700,color:#000
    style FRAMEWORK fill:#4CAF50,color:#fff
    style INTERNAL fill:#9E9E9E,color:#fff
    style HOSTABI fill:#2196F3,color:#fff
```

---

## 🔗 SDK 与平台组件的交互

### 完整应用架构

```mermaid
graph TB
    subgraph APP["完整应用"]
        CLIENT["Client SDK<br/>链下应用"]
        WORKBENCH["Workbench<br/>开发工具"]
        NODE["节点<br/>执行引擎"]
        SDK["Contract SDK Go<br/>合约代码"]
    end
    
    CLIENT -->|调用合约| NODE
    WORKBENCH -->|部署合约| NODE
    NODE -->|执行| SDK
    SDK -->|返回结果| NODE
    NODE -->|返回结果| CLIENT
    
    style CLIENT fill:#E3F2FD
    style WORKBENCH fill:#C8E6C9
    style NODE fill:#FFF9C4
    style SDK fill:#FFD700
```

---

## 💰 Token 转账场景流程

### 完整流程

```mermaid
sequenceDiagram
    participant Client as Client SDK
    participant Node as 节点
    participant SDK as Contract SDK
    
    Client->>Node: 构建交易（转账）
    Node->>SDK: 调用 Transfer()
    SDK->>SDK: 验证参数
    SDK->>SDK: 检查余额
    SDK->>SDK: 更新余额
    SDK->>Node: 发出事件
    SDK->>Node: 返回成功
    Node->>Client: 返回结果
    Node->>Client: 发出事件
```

---

## 📖 进一步阅读

### 核心文档

- **[SDK 内部架构](./SDK_ARCHITECTURE.md)** - SDK 内部分层架构设计
- **[应用场景分析](./APPLICATION_SCENARIOS_ANALYSIS.md)** - SDK 职责边界分析
- **[业务场景实现指南](./BUSINESS_SCENARIOS.md)** - 如何实现业务场景

### 平台文档（主仓库）

- [智能合约平台架构](../../../weisyn.git/docs/system/platforms/contracts/technical-architecture.md) - 平台技术架构
- [业务场景分析](../../../weisyn.git/docs/system/platforms/contracts/use-cases.md) - 详细业务流图

---

**最后更新**: 2025-01-23  
**维护者**: WES Core Team

