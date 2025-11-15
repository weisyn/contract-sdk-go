# WES 智能合约 SDK 模板库

**版本**: 1.0  
**最后更新**: 2025-11-15

---

## 📋 概述

本目录包含 WES 智能合约 SDK Go 的完整模板库，覆盖从基础入门到高级组合场景的各种业务场景。所有模板都使用 `helpers` 层的业务语义API，展示如何使用 SDK 构建生产级智能合约。

**模板定位**：SDK 提供「如何构建合约」的开发模板，开发者可以复制后修改使用。

---

## 🎯 为什么选择 WES 智能合约 SDK？

- ✅ **业务语义API**：使用 `helpers` 层提供的业务语义接口，开发者只需关注业务逻辑
- ✅ **ISPC 创新**：利用 ISPC 受控外部交互机制，无需传统预言机
- ✅ **自动处理**：SDK 自动处理交易构建、事件发出、余额检查等
- ✅ **类型安全**：使用 framework 类型系统，编译期类型检查
- ✅ **完整功能**：覆盖代币、NFT、RWA、质押、治理、市场、DeFi 等场景

---

## 📐 模板组织结构

```
templates/
├── README.md                    # 本文件：模板总览和导航
│
├── learning/                   # 学习模板（教学导向）
│   ├── README.md              # 学习模板说明
│   ├── hello-world/           # Hello World 模板
│   ├── simple-token/          # 简单代币模板
│   ├── basic-nft/             # 基础 NFT 模板
│   └── starter-contract/      # 空白模板
│
└── standard/                    # 标准模板（生产级）
    ├── README.md              # 标准模板说明
    ├── token/                  # 代币相关模板
    │   ├── erc20-token/       # ERC-20兼容代币
    │   ├── governance-token/  # 治理代币
    │   ├── payment-token/     # 支付代币
    │   └── game-currency/     # 游戏货币
    ├── nft/                    # NFT 相关模板
    ├── rwa/                    # RWA（现实世界资产）模板
    │   ├── real-estate/       # 房地产代币化
    │   │   ├── commercial/    # 商业地产
    │   │   └── residential/   # 住宅房产
    │   ├── equity/            # 股权代币化
    │   ├── commodity/         # 商品代币化
    │   ├── bond/              # 债券代币化
    │   └── artwork/           # 艺术品代币化
│   └── intellectual-property/ # 知识产权代币化
│
├── nft/                        # NFT（非同质化代币）示例
│   ├── README.md              # NFT示例说明
│   ├── digital-art/           # 数字艺术
│   ├── collectibles/         # 收藏品
│   ├── gaming/                # 游戏道具
│   ├── certificates/          # 证书凭证
│   ├── identity/             # 身份证明
│   ├── tickets/              # 门票票务
│   ├── domains/              # 域名服务
│   └── music/                # 音乐作品
│
├── staking/                    # 质押相关示例
│   ├── README.md              # 质押示例说明
│   ├── basic-staking/         # 基础质押
│   └── delegation/           # 委托质押
│
├── governance/                 # 治理相关示例
│   ├── README.md              # 治理示例说明
│   ├── proposal-voting/      # 提案投票
│   └── dao/                   # DAO治理
│
├── market/                     # 市场相关示例
│   ├── README.md              # 市场示例说明
│   ├── escrow/                # 托管交易
│   └── vesting/               # 分阶段释放
│
└── advanced/                   # 高级示例（组合场景）
    ├── README.md              # 高级示例说明
    └── defi/                  # DeFi组合场景
        ├── lending/           # 借贷协议
        ├── amm/               # AMM交换
        └── liquidity-pool/    # 流动性池
```

**文档说明**：
- **README.md**：每个目录都有 README，说明该分类的示例列表和使用方法

---

## 🚀 快速开始

### 1. 选择示例

根据您的需求选择合适的示例：

- **入门学习**：从 `basics/hello-world` 开始
- **代币开发**：查看 `token/` 目录下的示例
- **NFT开发**：查看 `nft/` 目录下的示例
- **RWA应用**：查看 `rwa/` 目录下的示例
- **DeFi应用**：查看 `advanced/defi/` 目录下的示例

### 2. 编译示例

```bash
# 进入示例目录
cd templates/standard/token/erc20-token

# 编译合约
bash build.sh

# 编译完成后会生成 main.wasm 文件
```

### 3. 部署合约

```bash
# 使用 WES CLI 部署
wes contract deploy --wasm main.wasm
```

### 4. 调用合约

```bash
# 调用合约函数
wes contract call --address {contract_addr} \
  --function Transfer \
  --params '{"to":"Df2Lft7toFVfjlKKhsBtLQOQsQbQeRnTn","amount":100}'
```

---

## 📚 示例分类

### 🌟 基础示例 (`basics/`)

**难度**: ⭐ - ⭐⭐ 入门到进阶  
**适合**: 初学者

- [Hello World](basics/hello-world/) ✅ - 最简单的入门示例
- [Simple Token](basics/simple-token/) ✅ - 基础代币功能

详见：[基础示例 README](basics/README.md)

---

### 💰 代币示例 (`token/`)

**难度**: ⭐⭐ - ⭐⭐⭐ 进阶到高级  
**适合**: 代币开发

- [ERC-20 Token](token/erc20-token/) ✅ - ERC-20兼容代币
- [Governance Token](token/governance-token/) ✅ - 治理代币
- [Payment Token](token/payment-token/) ✅ - 支付代币
- [Game Currency](token/game-currency/) ✅ - 游戏货币

详见：[代币示例 README](token/README.md)

---

### 🏠 RWA 示例 (`rwa/`)

**难度**: ⭐⭐⭐ 高级  
**适合**: 现实世界资产代币化

- [Real Estate - Commercial](rwa/real-estate/commercial/) ✅ - 商业地产代币化
- [Real Estate - Residential](rwa/real-estate/residential/) ✅ - 住宅房产代币化
- [Equity](rwa/equity/) ✅ - 股权代币化
- [Commodity](rwa/commodity/) ✅ - 商品代币化
- [Bond](rwa/bond/) ✅ - 债券代币化
- [Artwork](rwa/artwork/) ✅ - 艺术品代币化
- [Intellectual Property](rwa/intellectual-property/) ✅ - 知识产权代币化

详见：[RWA 示例 README](rwa/README.md)

---

### 🖼️ NFT 示例 (`nft/`)

**难度**: ⭐⭐ - ⭐⭐⭐ 进阶到高级  
**适合**: NFT开发

- [Digital Art](nft/digital-art/) ✅ - 数字艺术NFT
- [Collectibles](nft/collectibles/) ✅ - 收藏品NFT
- [Gaming](nft/gaming/) ✅ - 游戏道具NFT
- [Certificates](nft/certificates/) ✅ - 证书凭证NFT
- [Identity](nft/identity/) ✅ - 身份证明NFT
- [Tickets](nft/tickets/) ✅ - 门票票务NFT
- [Domains](nft/domains/) ✅ - 域名服务NFT
- [Music](nft/music/) ✅ - 音乐作品NFT

详见：[NFT 示例 README](nft/README.md)

---

### 🔒 质押示例 (`staking/`)

**难度**: ⭐⭐ - ⭐⭐⭐ 进阶到高级  
**适合**: 质押和委托开发

- [Basic Staking](staking/basic-staking/) ✅ - 基础质押
- [Delegation](staking/delegation/) ✅ - 委托质押

详见：[质押示例 README](staking/README.md)

---

### 🏛️ 治理示例 (`governance/`)

**难度**: ⭐⭐ - ⭐⭐⭐ 进阶到高级  
**适合**: 治理和DAO开发

- [Proposal Voting](governance/proposal-voting/) ✅ - 提案投票
- [DAO](governance/dao/) ✅ - DAO治理

详见：[治理示例 README](governance/README.md)

---

### 💼 市场示例 (`market/`)

**难度**: ⭐⭐ 进阶  
**适合**: 市场交易开发

- [Escrow](market/escrow/) ✅ - 托管交易
- [Vesting](market/vesting/) ✅ - 分阶段释放

详见：[市场示例 README](market/README.md)

---

### 🚀 高级示例 (`advanced/`)

**难度**: ⭐⭐⭐⭐ 高级  
**适合**: 复杂组合场景

- DeFi ✅ - DeFi组合场景（已完成）
  - [借贷协议](advanced/defi/lending/) ✅ - 完整的借贷功能（存款、借款、还款、取款）
  - [AMM交换](advanced/defi/amm/) ✅ - 自动化做市商（添加流动性、移除流动性、代币交换）
  - [流动性池](advanced/defi/liquidity-pool/) ✅ - 流动性池管理（添加流动性、移除流动性、查询池信息）
- GameFi - GameFi组合场景（待实现）

详见：[高级示例 README](advanced/README.md)

---

## 📊 示例统计

### 按优先级分类

| 优先级 | 数量 | 完成度 | 说明 |
|--------|------|--------|------|
| **P0** | 5 | 100% | 核心基础示例（必须实现） |
| **P1** | 7 | 100% | 重要业务场景（高优先级） |
| **P2** | 12 | 100% | 扩展场景（中优先级） |
| **P3** | 3 | 100% | 高级组合场景（低优先级） |
| **总计** | **27** | **100%** | 所有示例已完成 |

### 按业务场景分类

| 场景 | 数量 | 说明 |
|------|------|------|
| **基础** | 2 | 入门级示例 |
| **代币** | 4 | 各种代币类型 |
| **RWA** | 7 | 现实世界资产代币化 |
| **NFT** | 8 | 各种NFT应用场景 |
| **质押** | 2 | 质押和委托 |
| **治理** | 2 | 提案投票和DAO |
| **市场** | 2 | 托管和释放 |
| **DeFi** | 3 | 借贷、AMM、流动性池 |
| **总计** | **30** | 包含子场景 |

---

## 📝 每个示例包含的内容

每个示例都包含完整的代码和文档：

1. **main.go** - 完整的合约代码
   - 详细的代码注释
   - 函数说明（参数格式、返回值、事件）
   - 工作流程说明
   - ⚠️ 注意事项

2. **go.mod** - Go 模块定义
   - 正确的模块路径
   - 本地开发 replace 指令

3. **build.sh** - 编译脚本
   - TinyGo 编译命令
   - 错误处理

4. **README.md** - 详细文档
   - 示例概述
   - 核心功能列表
   - 架构图（Mermaid）
   - 功能详解
   - 快速开始指南
   - SDK vs 应用层职责
   - 设计理念
   - 相关文档链接

---

## 🎯 实现特点

### 1. 统一的代码风格

- 所有示例都遵循相同的代码结构和注释风格
- 使用 `helpers` 层的业务语义API
- 清晰的错误处理
- 完整的事件发出

### 2. 完整的文档

- 每个示例都有详细的 README
- 包含架构图和流程图（Mermaid）
- 提供使用示例和最佳实践
- SDK vs 应用层职责清晰划分

### 3. ISPC 创新体现

- RWA 示例展示 ISPC 受控外部交互机制
- 无需传统预言机
- 自动生成 ZK 证明
- 业务执行即上链

### 4. 业务语义优先

- 使用 `helpers` 层提供的业务语义API
- 开发者只需关注业务逻辑
- SDK 自动处理交易构建、事件发出等

---

## 🔄 模板库定位

**说明**：本目录（`templates/`）是 SDK 提供的合约开发模板库，所有模板已统一迁移至此。

**历史说明**：

| 特性 | 旧位置（已迁移） | 新位置（当前） |
|------|----------------|--------------|
| **contracts/templates/** | WES 主仓库中的模板 | 已迁移到 `templates/` |
| **_sdks/contract-sdk-go/examples/** | SDK 中的示例 | 已迁移到 `templates/` |

**当前结构**：

- **templates/learning/**：学习模板（教学导向），开发者可以学习 SDK 使用方式
- **templates/standard/**：标准模板（生产级），开发者可以复制后修改使用

两者共同为开发者提供完整的开发支持。

### 参考的最佳实践

- ✅ **学习路径设计**：基础 → 进阶 → 标准
- ✅ **文档结构**：学习目标、概念解释、功能详解、使用示例、架构图
- ✅ **代码组织**：清晰的函数注释、错误处理、事件发出、参数验证
- ✅ **示例分类**：按业务场景、难度、优先级分类

---

## 🎯 学习路径建议

### 初学者路径

1. **基础入门**
   - `basics/hello-world` - 理解基本结构
   - `basics/simple-token` - 学习状态管理

2. **代币开发**
   - `token/erc20-token` - 标准代币实现
   - `token/governance-token` - 治理代币

3. **NFT开发**
   - `nft/digital-art` - 数字艺术NFT
   - `nft/collectibles` - 收藏品NFT

### 进阶路径

1. **RWA应用**
   - `rwa/real-estate/commercial` - 商业地产代币化
   - `rwa/equity` - 股权代币化

2. **质押和治理**
   - `staking/basic-staking` - 基础质押
   - `governance/proposal-voting` - 提案投票

3. **市场功能**
   - `market/escrow` - 托管交易
   - `market/vesting` - 分阶段释放

### 高级路径

1. **DeFi应用**
   - `advanced/defi/lending` - 借贷协议
   - `advanced/defi/amm` - AMM交换
   - `advanced/defi/liquidity-pool` - 流动性池

2. **组合场景**
   - 参考多个示例，组合实现复杂业务场景

---

## 🌟 ISPC 创新体现

### RWA 示例中的 ISPC 创新

所有 RWA 示例都利用 ISPC 受控外部交互机制：

- ✅ **无需传统预言机**：直接调用外部验证和估值服务
- ✅ **自动生成ZK证明**：验证和估值过程自动生成可验证性证明
- ✅ **单次调用保证**：只有执行节点调用外部服务，验证节点只验证证明

**示例**：
- `rwa/real-estate/commercial` - 使用 `rwa.ValidateAndTokenize()` 进行资产验证和代币化
- `rwa/equity` - 使用 ISPC 受控机制调用外部验证服务

### 治理示例中的 ISPC 创新

治理示例展示 ISPC 业务执行范式：

- ✅ **业务执行即上链**：投票、统计等业务逻辑执行后自动上链
- ✅ **自动生成ZK证明**：业务逻辑执行过程自动生成可验证性证明

**示例**：
- `governance/proposal-voting` - 使用 `governance.Vote()` 进行投票
- `governance/dao` - 使用 `governance.VoteAndCount()` 进行投票和统计

---

## 🔗 相关文档

### 平台文档（高层次视图）

- [智能合约平台文档](../../docs/system/platforms/contracts/README.md) - 智能合约平台的综合文档
  - [应用场景](../../docs/system/platforms/contracts/use-cases.md) - 实际应用案例（包含 SDK 示例）
  - [快速开始](../../docs/system/platforms/contracts/getting-started.md) - 开发者快速入门（包含 SDK 使用）
  - [产品设计](../../docs/system/platforms/contracts/product-design.md) - 产品特性和用户体验（包含 SDK 设计）

### SDK 文档

- [Contract SDK 主 README](../README.md) - SDK 总览和快速开始
- [Framework 文档](../framework/README.md) - Framework 层说明
- [Helpers 文档](../helpers/README.md) - Helpers 层说明

### 开发实践文档

- [资源级示例](../../contracts/examples/README.md) - WES 平台提供的测试合约示例
- [合约开发平台](../../contracts/README.md) - 模板库、工具链、系统合约

### 教程文档

- [合约学习路径](../../docs/tutorials/contracts/LEARNING_PATH.md) - 分阶段学习路径
- [合约核心概念](../../docs/tutorials/contracts/CONCEPTS.md) - 核心概念解释

---

## 💡 设计理念

### SDK 提供"积木"

SDK 提供基础能力（Transfer、Mint、Burn、Escrow、Release等），开发者可以：

- ✅ 直接使用基础功能创建应用
- ✅ 添加业务规则实现定制需求
- ✅ 组合多个功能实现复杂场景

### 应用层搭建"建筑"

应用层在 SDK 基础上实现：

- ✅ 业务规则（权限检查、条件验证等）
- ✅ 状态管理（使用状态输出存储业务状态）
- ✅ 复杂逻辑（利率计算、价格发现等）

---

## 📄 许可证

Apache-2.0 License

---

**最后更新**: 2025-11-11
