# RWA - 现实世界资产代币化模块

**版本**: 1.0  
**状态**: ✅ 稳定  
**最后更新**: 2025-11-11

---

## 📋 概述

RWA（Real World Assets）模块提供现实世界资产代币化能力，支持资产验证、价值评估和代币化操作。该模块支持外部系统集成，无需传统预言机。

---

## 🎯 核心功能

### ValidateAndTokenize - 验证并代币化资产

**功能**：完整的资产验证和代币化流程

**签名**：
```go
func ValidateAndTokenize(
    assetID string,
    documents []byte,
    validatorAPI string,
    validatorEvidence *framework.Evidence,
    valuationAPI string,
    valuationEvidence *framework.Evidence,
) (*TokenizeResult, error)
```

**参数说明**：
- `assetID` - 资产ID
- `documents` - 资产文档（JSON格式）
- `validatorAPI` - 验证服务API端点
- `validatorEvidence` - 验证佐证（API签名、响应哈希等）
- `valuationAPI` - 估值服务API端点
- `valuationEvidence` - 估值佐证（API签名、响应哈希等）

**返回结果**：
```go
type TokenizeResult struct {
    TokenID         framework.TokenID  // 生成的代币ID
    Validated       bool               // 是否验证通过
    ValidationProof []byte            // 验证过程证明
    Valuation       uint64            // 资产估值
    ValuationProof  []byte            // 估值过程证明
    TxHash          framework.Hash    // 交易哈希
}
```

**示例**：
```go
import (
    "github.com/weisyn/contract-sdk-go/helpers/rwa"
    "github.com/weisyn/contract-sdk-go/framework"
)

//export TokenizeAsset
func TokenizeAsset() uint32 {
    params := framework.GetContractParams()
    assetID := params.ParseJSON("asset_id")
    documents := []byte(params.ParseJSON("documents"))
    
    // 验证并代币化资产
    result, err := rwa.ValidateAndTokenize(
        assetID,
        documents,
        "https://validator.example.com/api/validate",
        &framework.Evidence{
            APISignature: []byte("validator_signature"),
            ResponseHash: []byte("validation_hash"),
        },
        "https://valuation.example.com/api/value",
        &framework.Evidence{
            APISignature: []byte("valuation_signature"),
            ResponseHash: []byte("valuation_hash"),
        },
    )
    if err != nil {
        return framework.ERROR_EXECUTION_FAILED
    }
    
    // 发出事件
    framework.EmitEvent("AssetTokenized", []byte(`{
        "token_id":"`+string(result.TokenID)+`",
        "validated":`+fmt.Sprintf("%t", result.Validated)+`,
        "valuation":`+fmt.Sprintf("%d", result.Valuation)+`
    }`))
    
    return framework.SUCCESS
}
```

### ValidateAsset - 验证资产

**功能**：验证资产文档和合法性

**签名**：
```go
func ValidateAsset(
    assetID string,
    documents []byte,
    validatorAPI string,
    evidence *framework.Evidence,
) (bool, []byte, error)
```

**示例**：
```go
validated, proof, err := rwa.ValidateAsset(
    assetID,
    documents,
    "https://validator.example.com/api/validate",
    validatorEvidence,
)
```

### ValueAsset - 评估资产价值

**功能**：评估资产的市场价值

**签名**：
```go
func ValueAsset(
    assetID string,
    documents []byte,
    valuationAPI string,
    evidence *framework.Evidence,
) (uint64, []byte, error)
```

**示例**：
```go
valuation, proof, err := rwa.ValueAsset(
    assetID,
    documents,
    "https://valuation.example.com/api/value",
    valuationEvidence,
)
```

---

## 💡 使用场景

### 场景1：房地产代币化

```go
//export TokenizeRealEstate
func TokenizeRealEstate() uint32 {
    params := framework.GetContractParams()
    propertyID := params.ParseJSON("property_id")
    
    // 构建资产文档
    documents := []byte(`{
        "property_id": "` + propertyID + `",
        "type": "real_estate",
        "location": "...",
        "area": 100,
        "certificate": "..."
    }`)
    
    // 验证并代币化
    result, err := rwa.ValidateAndTokenize(
        propertyID,
        documents,
        "https://real-estate-validator.com/api/validate",
        validatorEvidence,
        "https://real-estate-valuation.com/api/value",
        valuationEvidence,
    )
    
    // 使用生成的代币ID
    // result.TokenID 可用于后续交易
}
```

### 场景2：艺术品代币化

```go
//export TokenizeArtwork
func TokenizeArtwork() uint32 {
    // 类似流程，使用艺术品验证和估值服务
    result, err := rwa.ValidateAndTokenize(
        artworkID,
        artworkDocuments,
        "https://art-validator.com/api/validate",
        validatorEvidence,
        "https://art-auction.com/api/value",
        valuationEvidence,
    )
}
```

---

## 🔄 工作流程

```
1. 调用 ValidateAndTokenize
   ↓
2. 调用外部验证服务（受控机制）
   ↓
3. 调用外部估值服务（受控机制）
   ↓
4. 生成代币（使用 token.Mint）
   ↓
5. 自动构建交易并上链
   ↓
6. 返回结果（TokenID、验证结果、估值等）
```

---

## ⚠️ 注意事项

1. **外部服务要求**：验证和估值服务需要提供数字签名和响应哈希作为佐证
2. **证据提供**：调用者需要提供 `Evidence` 结构，包含 API 签名和响应哈希
3. **错误处理**：外部服务调用失败时，函数会返回错误，合约应妥善处理
4. **成本考虑**：外部 API 调用会产生执行成本，应合理设计调用频率

---

## 🔗 相关文档

- [Helpers 层总览](../README.md) - 业务语义层总览
- [External 模块文档](../external/README.md) - 外部系统集成说明
- [Token 模块文档](../token/README.md) - 代币操作说明
- [Contract SDK 主 README](../../README.md) - SDK 总览

---

## 💡 设计说明

### 外部系统集成

RWA 模块通过 `helpers/external` 模块实现外部系统集成：
- 支持外部 API 调用（验证服务、估值服务）
- 支持数据库查询（资产数据库）
- 所有外部交互都有密码学验证的佐证

### 与传统预言机的区别

| 特性 | 传统预言机 | RWA 模块 |
|------|-----------|---------|
| **中心化风险** | 需要中心化预言机服务 | 直接调用外部服务，无需中介 |
| **验证方式** | 依赖预言机信任 | 密码学验证佐证 |
| **成本** | 预言机服务费用 | 直接调用，成本更低 |
| **灵活性** | 受限于预言机支持 | 可调用任意外部服务 |

---

**最后更新**: 2025-11-11
