# External - 受控外部交互模块

**版本**: 1.0  
**状态**: ✅ 稳定  
**最后更新**: 2025-11-11

---

## 📋 概述

External 模块提供受控外部交互能力，通过ISPC的"受控声明+佐证+验证"机制，替代传统预言机。

**模块定位**：
- ✅ **业务语义封装**：封装 Framework 层 HostABI 的业务语义
- ✅ **调用 Framework**：调用 `framework.DeclareExternalState()` 等函数
- ✅ **与其他 helpers 模块平级**：都是业务语义层的一部分
- ✅ **区别**：封装的是 ISPC 创新的能力（受控外部交互）

**ISPC 创新体现**：
- ✅ **无需传统预言机**：直接通过 ISPC 受控机制获取外部数据
- ✅ **可验证的外部交互**：所有外部调用都有密码学验证的佐证
- ✅ **单次调用保证**：只有执行节点调用一次，验证节点只验证证明
- ✅ **自动生成 ZK 证明**：外部交互验证过程包含在 ZK 证明中

**架构定位**：
```
HostABI 原语（ISPC 提供）
    ↓
framework/  → HostABI 封装（DeclareExternalState 等）
    ↓
helpers/external/  → 业务语义封装（ValidateAndQuery 等）
    ↓
helpers/rwa/, helpers/governance/ 等  → 使用 external 模块
```

---

## 🎯 核心功能

### 1. ValidateAndQuery - 验证并查询外部状态

**用途**：通用的受控外部状态查询接口

**流程**：
1. 声明外部状态预期
2. 提供验证佐证
3. 运行时验证
4. 查询已验证的外部状态

### 2. CallAPI - 受控外部API调用

**用途**：调用外部API，替代传统预言机

**特点**：
- 通过API签名和响应哈希验证
- 单次调用，多点验证
- 自动生成ZK证明

### 3. QueryDatabase - 受控数据库查询

**用途**：查询数据库，替代传统预言机

**特点**：
- 通过状态哈希和默克尔证明验证
- 单次查询，多点验证
- 自动生成ZK证明

---

## 💡 使用示例

### 示例1：查询API价格

```go
import (
    "github.com/weisyn/contract-sdk-go/helpers/external"
    "github.com/weisyn/contract-sdk-go/framework"
)

func GetPrice() uint32 {
    // 调用外部API获取价格
    data, err := external.CallAPI(
        "https://api.example.com/price",
        "GET",
        map[string]interface{}{"symbol": "BTC"},
        apiSignature,    // API数字签名
        responseHash,    // 响应数据哈希
    )
    if err != nil {
        return framework.ERROR_EXECUTION_FAILED
    }
    
    // 使用data进行业务逻辑处理
    // ...
    
    return framework.SUCCESS
}
```

### 示例2：查询数据库

```go
func QueryUser() uint32 {
    // 查询数据库
    data, err := external.QueryDatabase(
        "user_db",
        "SELECT * FROM users WHERE id = ?",
        []interface{}{userID},
        stateHash,      // 数据库状态哈希
        merkleProof,    // 默克尔证明
    )
    if err != nil {
        return framework.ERROR_EXECUTION_FAILED
    }
    
    // 使用data进行业务逻辑处理
    // ...
    
    return framework.SUCCESS
}
```

---

## 🔗 相关文档

- [ISPC 创新分析](../docs/ISPC_INNOVATION_ANALYSIS.md)
- [Framework 层文档](../../framework/README.md)

---

**最后更新**: 2025-11-11

