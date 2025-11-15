# Contract SDK 整体结构设计讨论

**创建日期**: 2025-11-11  
**状态**: ✅ 已完成  
**最后更新**: 2025-11-11  
**目标**: 从整体架构角度，系统性地设计SDK结构，体现WES的ISPC创新

> **注意**: 本文档记录了 SDK 的整体结构设计讨论和最终结论。

---

## 🎯 当前整体结构

```
contract-sdk-go/
├── framework/          # Framework层：HostABI封装
│   ├── internal/      # 内部实现（交易构建等）
│   └── ...
├── helpers/           # Helpers层：业务语义封装
│   ├── token/        # Token业务模块
│   ├── staking/      # Staking业务模块
│   ├── market/       # Market业务模块
│   ├── governance/   # Governance业务模块
│   ├── rwa/          # RWA业务模块
│   └── external/     # ⚠️ 问题：ISPC能力与其他业务模块平级
├── examples/          # 示例代码
├── docs/             # 文档（讨论文档应在此）
└── README.md         # 主README
```

---

## 🔍 问题分析

### 核心疑问：ISPC内部能力层到底是什么？

**关键问题**：
1. ISPC是底层执行范式，还是SDK中的一个独立层？
2. 受控外部交互能力应该通过HostABI体现，还是独立的"ISPC层"？
3. `helpers/external/` 的本质是什么？它与其他业务模块的区别是什么？

### ISPC的本质定位

**ISPC的本质**：
- ✅ **ISPC是WES的执行范式**：定义"可验证计算"范式，不是SDK中的一个层
- ✅ **HostABI是ISPC提供的原语**：包括受控外部交互的原语（`host_declare_external_state`等）
- ✅ **Framework层封装HostABI**：`framework.DeclareExternalState()`等函数封装了HostABI原语
- ✅ **Helpers层使用Framework能力**：`helpers/external/`调用`framework.DeclareExternalState()`

**当前结构分析**：
```
framework/
├── host_functions.go    # HostABI函数声明（包括受控外部交互）
├── hostabi.go          # HostABI封装（framework.DeclareExternalState等）
└── ...

helpers/
├── external/           # ⚠️ 问题：这是对Framework的封装，还是独立的"ISPC层"？
│   └── service.go      # 调用framework.DeclareExternalState()
└── rwa/                # 使用external
```

**核心问题**：
- `helpers/external/` 实际上是对Framework层HostABI封装的**业务语义封装**
- 它调用的是`framework.DeclareExternalState()`等函数
- **不应该有独立的"ISPC层"**，ISPC是底层执行范式，通过HostABI体现
- `helpers/external/` 应该与其他helpers模块一样，是对Framework的业务语义封装

---

## 💡 正确的结构设计

### ISPC的本质：底层执行范式，通过HostABI体现

**ISPC在SDK中的体现**：
1. **ISPC是底层执行范式**：定义"可验证计算"范式，不是SDK中的一个层
2. **HostABI是ISPC提供的原语**：包括受控外部交互的原语（`host_declare_external_state`等）
3. **Framework层封装HostABI**：提供类型安全的API（`framework.DeclareExternalState()`等）
4. **Helpers层使用Framework能力**：按照ISPC的思路组织业务逻辑

### 正确的结构设计

```
contract-sdk-go/
├── framework/              # L1: Framework层（HostABI封装）
│   ├── host_functions.go  # HostABI函数声明（包括受控外部交互）
│   ├── hostabi.go         # HostABI封装（包括DeclareExternalState等）
│   ├── types.go           # 类型定义（包括ExternalStateClaim, Evidence）
│   └── internal/          # 内部实现（交易构建、状态管理等）
│
├── helpers/               # L2: Helpers层（业务语义封装）
│   ├── token/            # Token业务模块
│   ├── staking/          # Staking业务模块
│   ├── market/           # Market业务模块
│   ├── governance/       # Governance业务模块
│   ├── rwa/              # RWA业务模块（使用external）
│   └── external/         # 受控外部交互业务语义封装（使用framework）
│       └── service.go    # 调用framework.DeclareExternalState()等
│
├── examples/              # 示例代码
├── docs/                  # 文档
└── README.md
```

### 关键设计点

1. **ISPC不是SDK中的一个层**
   - ISPC是底层执行范式，通过HostABI体现
   - Framework层封装HostABI，包括受控外部交互的HostABI函数
   - Helpers层使用Framework能力，按照ISPC的思路组织业务逻辑

2. **清晰的层次关系**
   ```
   HostABI原语（ISPC提供）
        ↓
   framework/  → HostABI封装（包括受控外部交互）
        ↓
   helpers/    → 业务语义封装（使用framework，按照ISPC思路组织）
        ↓
   examples/   → 示例代码（使用helpers）
   ```

3. **external的本质**
   - `helpers/external/` 是对Framework层HostABI封装的**业务语义封装**
   - 它调用`framework.DeclareExternalState()`等函数
   - 与其他helpers模块（token, staking等）一样，都是业务语义封装
   - **区别**：它封装的是ISPC创新的能力（受控外部交互），但仍然是业务语义层的一部分

4. **统一对外接口**
   - 合约开发者只使用 `helpers/` 层
   - `helpers/external/` 与其他helpers模块平级，都是业务语义封装
   - Framework层作为基础能力，供helpers使用

---

## 🏗️ 详细结构设计

### 1. Framework层（L1）

**定位**：HostABI封装，提供基础能力

**结构**：
```
framework/
├── host_functions.go    # HostABI函数声明（包括受控外部交互）
│   ├── hostDeclareExternalState()
│   ├── hostProvideEvidence()
│   └── hostQueryControlledState()
├── hostabi.go          # HostABI封装（类型安全API）
│   ├── DeclareExternalState()
│   ├── ProvideEvidence()
│   └── QueryControlledState()
├── types.go            # 类型定义
│   ├── ExternalStateClaim
│   └── Evidence
├── contract_base.go    # 合约基础结构
├── errors.go           # 错误定义
└── internal/           # 内部实现
    ├── transaction.go  # 交易构建
    ├── state.go        # 状态管理
    ├── utxo.go         # UTXO操作
    └── ...
```

**职责**：
- 封装HostABI原语（包括受控外部交互的HostABI函数）
- 提供类型安全的API
- 隐藏底层实现细节

**ISPC体现**：
- Framework层封装了ISPC提供的HostABI原语
- 包括受控外部交互的HostABI函数（`host_declare_external_state`等）
- ISPC的执行范式通过HostABI体现

---

### 2. Helpers层（L2）- 业务语义层

**定位**：业务语义封装，统一对外接口

**结构**：
```
helpers/
├── token/              # Token业务模块
│   ├── transfer.go    # 使用framework.QueryUTXOBalance()等
│   ├── mint.go         # 使用framework.CreateAssetOutput()等
│   └── ...
├── staking/           # Staking业务模块
├── market/            # Market业务模块
├── governance/        # Governance业务模块
│   ├── propose.go     # 使用framework.AddStateOutput()等
│   ├── vote.go        # 使用framework.AddStateOutput()等
│   └── vote_and_count.go  # 使用framework，按照ISPC思路组织业务逻辑
├── rwa/               # RWA业务模块
│   ├── validate.go    # 使用helpers/external
│   ├── value.go       # 使用helpers/external
│   └── tokenize.go    # 使用helpers/external和helpers/token
└── external/          # 受控外部交互业务语义封装
    └── service.go     # 调用framework.DeclareExternalState()等
```

**职责**：
- 提供业务语义接口（Transfer, Mint, Stake等）
- 使用framework能力
- 按照ISPC的思路组织业务逻辑（如受控外部交互、业务执行）

**external模块的本质**：
- `helpers/external/` 是对Framework层HostABI封装的**业务语义封装**
- 它调用`framework.DeclareExternalState()`等函数
- 与其他helpers模块（token, staking等）一样，都是业务语义封装
- **区别**：它封装的是ISPC创新的能力（受控外部交互），但仍然是业务语义层的一部分

**ISPC体现**：
- Helpers层使用Framework的HostABI封装
- 按照ISPC的思路组织业务逻辑（如受控外部交互、业务执行）
- ISPC的执行范式通过业务逻辑的组织方式体现

**设计原则**：
- 统一整体：所有模块按业务领域分类
- 渐进增强：基础功能 + ISPC增强功能
- ISPC赋能：通过Framework的HostABI封装，按照ISPC思路组织业务逻辑

---

## 📊 依赖关系

### 依赖图

```
examples/
    ↓ 使用
helpers/ (业务语义层)
    ↓ 使用
framework/ (HostABI封装)
    ↓ 使用
HostABI原语 (ISPC提供)
```

### 关键原则

1. **单向依赖**
   - `helpers/` → `framework/`
   - `framework/` → HostABI原语（ISPC提供）

2. **ISPC的体现方式**
   - ISPC是底层执行范式，通过HostABI体现
   - Framework层封装HostABI（包括受控外部交互的HostABI函数）
   - Helpers层使用Framework能力，按照ISPC的思路组织业务逻辑

3. **统一整体**
   - 合约开发者只看到 `helpers/` 层
   - 不需要知道ISPC和传统的区别
   - ISPC的执行范式通过业务逻辑的组织方式体现

---

## 🎯 实施建议

### 保持现有结构，明确ISPC定位

**不需要创建独立的"ISPC层"**，因为：
1. ISPC是底层执行范式，通过HostABI体现
2. Framework层已经封装了HostABI（包括受控外部交互的HostABI函数）
3. `helpers/external/` 是对Framework的业务语义封装，与其他helpers模块一样

### 需要做的调整

1. **明确ISPC定位**
   - ISPC是底层执行范式，不是SDK中的一个层
   - Framework层封装HostABI（ISPC提供的原语）
   - Helpers层使用Framework能力，按照ISPC的思路组织业务逻辑

2. **明确external模块的定位**
   - `helpers/external/` 是对Framework层HostABI封装的业务语义封装
   - 它调用`framework.DeclareExternalState()`等函数
   - 与其他helpers模块（token, staking等）一样，都是业务语义层的一部分
   - **区别**：它封装的是ISPC创新的能力（受控外部交互）

3. **更新文档**
   - 更新主README，说明ISPC的定位
   - 更新helpers/README，说明external模块的定位
   - 强调ISPC通过HostABI体现，通过业务逻辑的组织方式体现

---

## 📋 对比分析

| 维度 | 错误理解 | 正确理解 |
|------|---------|---------|
| **ISPC定位** | SDK中的一个独立层 | 底层执行范式，通过HostABI体现 |
| **external模块** | ISPC内部能力层 | Framework的业务语义封装 |
| **层次关系** | framework → ispc → helpers | framework → helpers |
| **ISPC体现** | 独立的"ISPC层" | 通过HostABI和业务逻辑组织方式体现 |

---

## 🎯 关键设计原则

1. **ISPC是底层执行范式，不是SDK中的一个层**
   - ISPC通过HostABI体现
   - Framework层封装HostABI（包括受控外部交互的HostABI函数）
   - Helpers层使用Framework能力，按照ISPC的思路组织业务逻辑

2. **统一整体**
   - 合约开发者只使用helpers层
   - 不需要知道ISPC和传统的区别
   - ISPC的执行范式通过业务逻辑的组织方式体现

3. **external模块的本质**
   - `helpers/external/` 是对Framework层HostABI封装的业务语义封装
   - 与其他helpers模块一样，都是业务语义层的一部分
   - **区别**：它封装的是ISPC创新的能力（受控外部交互）

4. **清晰层次**
   - HostABI原语（ISPC提供） → framework → helpers → examples
   - 单向依赖，职责分明
   - ISPC的执行范式贯穿整个层次

---

## 🔄 不需要迁移

### 保持现有结构

**不需要创建独立的"ISPC层"**，因为：
1. ISPC是底层执行范式，通过HostABI体现
2. Framework层已经封装了HostABI（包括受控外部交互的HostABI函数）
3. `helpers/external/` 是对Framework的业务语义封装，与其他helpers模块一样

### 只需要更新文档

1. **更新主README**
   - 说明ISPC的定位（底层执行范式，通过HostABI体现）
   - 说明整体结构（framework → helpers → examples）

2. **更新helpers/README**
   - 说明external模块的定位（Framework的业务语义封装）
   - 强调ISPC通过HostABI体现，通过业务逻辑的组织方式体现

---

## 💡 结论

### ISPC在SDK中的正确体现

1. **ISPC是底层执行范式**
   - 不是SDK中的一个独立层
   - 通过HostABI原语体现（包括受控外部交互的原语）

2. **Framework层封装HostABI**
   - 封装了ISPC提供的HostABI原语
   - 包括受控外部交互的HostABI函数（`DeclareExternalState`等）

3. **Helpers层使用Framework能力**
   - 按照ISPC的思路组织业务逻辑
   - `helpers/external/` 是对Framework的业务语义封装
   - 与其他helpers模块一样，都是业务语义层的一部分

4. **不需要独立的"ISPC层"**
   - ISPC的执行范式通过HostABI和业务逻辑的组织方式体现
   - 不需要硬凑一个"ISPC层"
   - 保持现有结构，明确ISPC定位即可

---

**最后更新**: 2025-11-11

