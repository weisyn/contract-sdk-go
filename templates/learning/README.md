# WES 合约学习模板库（SDK templates/learning）

---

## 📌 版本信息

- **版本**：1.0
- **状态**：stable
- **最后更新**：2025-11-15
- **最后审核**：2025-11-15
- **所有者**：合约 SDK 组
- **适用范围**：WES 智能合约 SDK 学习模板

---

## 🎯 目录定位

**路径**：`_sdks/contract-sdk-go/templates/learning/`

**目的**：提供**教学导向**的合约模板，帮助开发者学习 WES 合约开发的基础概念和最佳实践。

**与其他目录的关系**：

- `templates/learning/`：**学习模板**，侧重教学，代码中有大量注释，可复制后修改
- `templates/standard/`：**标准模板**，生产级骨架，可复制后定制
- `contracts/examples/`：**资源级示例**，行为固定，测试完备，用于验证平台能力
- `examples/`（仓库根）：**场景级示例**，组合使用模型、合约等多种资源

---

## 📐 目录结构

```
learning/
├── README.md                    # 本文档
├── simple-token/                # 代币学习模板
│   ├── README.md
│   ├── main.go
│   ├── go.mod
│   └── build.sh
├── hello-world/                 # Hello World 学习模板
│   ├── README.md
│   ├── main.go
│   ├── go.mod
│   └── build.sh
├── basic-nft/                   # NFT 学习模板
│   └── ...
└── starter-contract/            # 入门合约模板
    └── ...
```

---

## 🎓 学习模板的特点

### 与资源级示例的区别

| 维度 | `templates/learning/` | `contracts/examples/` |
|------|----------------------|----------------------|
| **定位** | 教学模板 | 资源级示例 |
| **代码风格** | 大量注释，解释概念 | 简洁，行为固定 |
| **可修改性** | ✅ 鼓励修改和实验 | ❌ 行为固定，不鼓励修改 |
| **测试用例** | 可能有基础测试 | 完整的标准测试用例 |
| **用途** | 学习如何开发合约 | 验证平台能力 |

### 使用建议

1. **初学者**：从 `learning/` 开始，理解基础概念
2. **进阶开发者**：参考 `standard/` 模板，了解生产级实现
3. **测试验证**：使用 `contracts/examples/` 验证平台能力

---

## 📚 模板列表

### simple-token

**路径**：`learning/simple-token/`

**功能**：代币合约学习模板，展示基础转账、余额查询等功能。

**学习重点**：
- 代币的基本概念
- UTXO 模型的使用
- 事件发出机制
- 状态查询方法

### hello-world

**路径**：`learning/hello-world/`

**功能**：Hello World 学习模板，展示最基本的合约结构。

**学习重点**：
- 合约的基本结构
- 导出函数定义
- 事件发出
- 返回值处理

### basic-nft

**路径**：`learning/basic-nft/`

**功能**：NFT 学习模板，展示 NFT 铸造、转移等功能。

**学习重点**：
- NFT 的基本概念
- 唯一性管理
- 元数据存储
- 所有权转移

### starter-contract

**路径**：`learning/starter-contract/`

**功能**：入门合约模板，提供空白但结构完整的合约骨架。

**学习重点**：
- 合约框架结构
- 自定义功能开发
- 最佳实践应用

---

## 🔗 相关文档

- `templates/README.md`：SDK 模板库总览
- `templates/standard/README.md`：标准模板库说明
- `contracts/examples/README.md`：资源级合约示例库

---
