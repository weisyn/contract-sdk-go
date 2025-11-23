# NFT 业务语义模块

**版本**: 1.0  
**状态**: ✅ 稳定  
**最后更新**: 2025-11-11

---

## 📋 概述

NFT 模块提供非同质化代币（NFT）相关的业务语义API，包括铸造、转移、销毁等功能。

**定位**：NFT 模块是对 Token 模块的扩展，专门用于处理 NFT 业务场景。NFT 本质上是数量为1的特殊代币。

---

## 🎯 核心功能

### 1. Mint - 铸造NFT

**功能**: 铸造新的NFT

**签名**:
```go
func Mint(to framework.Address, tokenID framework.TokenID, metadata []byte) error
```

**示例**:
```go
err := nft.Mint(to, framework.TokenID("nft_001"), []byte(`{"name":"My NFT"}`))
```

**输入输出组合模式**:
- `token.Mint()` - 铸造数量为1的代币
- `StateOutput` - 存储NFT元数据（可选）

---

### 2. Transfer - 转移NFT

**功能**: 转移NFT所有权

**签名**:
```go
func Transfer(from, to framework.Address, tokenID framework.TokenID) error
```

**示例**:
```go
err := nft.Transfer(from, to, framework.TokenID("nft_001"))
```

**输入输出组合模式**:
- `token.Transfer()` - 转移数量为1的代币

---

### 3. Burn - 销毁NFT

**功能**: 销毁NFT

**签名**:
```go
func Burn(from framework.Address, tokenID framework.TokenID) error
```

**示例**:
```go
err := nft.Burn(caller, framework.TokenID("nft_001"))
```

**输入输出组合模式**:
- `token.Burn()` - 销毁数量为1的代币

---

### 4. OwnerOf - 查询所有者

**功能**: 查询NFT的所有者地址

**签名**:
```go
func OwnerOf(tokenID framework.TokenID) *framework.Address
```

**示例**:
```go
owner := nft.OwnerOf(framework.TokenID("nft_001"))
if owner != nil {
    // 使用所有者地址
}
```

---

### 5. BalanceOf - 查询余额

**功能**: 查询地址拥有的NFT数量

**签名**:
```go
func BalanceOf(owner framework.Address) uint64
```

**示例**:
```go
count := nft.BalanceOf(owner)
```

---

### 6. GetMetadata - 获取元数据

**功能**: 查询NFT的元数据

**签名**:
```go
func GetMetadata(tokenID framework.TokenID) []byte
```

**示例**:
```go
metadata := nft.GetMetadata(framework.TokenID("nft_001"))
```

---

## 💡 使用示例

### 完整示例：NFT合约

```go
package main

import (
    "github.com/weisyn/contract-sdk-go/helpers/nft"
    "github.com/weisyn/contract-sdk-go/framework"
)

//export MintNFT
func MintNFT() uint32 {
    params := framework.GetContractParams()
    toStr := params.ParseJSON("to")
    tokenIDStr := params.ParseJSON("token_id")
    metadataStr := params.ParseJSON("metadata")
    
    to, err := framework.ParseAddressBase58(toStr)
    if err != nil {
        return framework.ERROR_INVALID_PARAMS
    }
    
    err = nft.Mint(to, framework.TokenID(tokenIDStr), []byte(metadataStr))
    if err != nil {
        return framework.ERROR_EXECUTION_FAILED
    }
    
    return framework.SUCCESS
}

//export TransferNFT
func TransferNFT() uint32 {
    params := framework.GetContractParams()
    toStr := params.ParseJSON("to")
    tokenIDStr := params.ParseJSON("token_id")
    
    to, err := framework.ParseAddressBase58(toStr)
    if err != nil {
        return framework.ERROR_INVALID_PARAMS
    }
    
    caller := framework.GetCaller()
    err = nft.Transfer(caller, to, framework.TokenID(tokenIDStr))
    if err != nil {
        return framework.ERROR_EXECUTION_FAILED
    }
    
    return framework.SUCCESS
}
```

---

## 📊 事件语义文档

NFT 模块发出的所有事件都遵循统一的语义规范。下表列出了所有事件的结构和字段含义：

| 事件名 | 字段名 | 类型 | 说明 |
|--------|--------|------|------|
| **NFTMint** | `to` | Address (Base58) | 接收者地址 |
| | `token_id` | string | NFT代币ID |
| | `minter` | Address (Base58) | 铸造者地址 |
| | `metadata` | string | NFT元数据（可选） |
| **NFTTransfer** | `from` | Address (Base58) | 发送者地址 |
| | `to` | Address (Base58) | 接收者地址 |
| | `token_id` | string | NFT代币ID |
| **NFTBurn** | `from` | Address (Base58) | 销毁者地址 |
| | `token_id` | string | NFT代币ID |

**事件格式说明**：
- 所有地址字段使用 Base58 编码
- 事件结构作为公共契约，未来只能增加字段、不能删减

---

## 🔗 相关文档

- [Contract Helpers总览](../README.md)
- [Token 模块文档](../token/README.md) - NFT基于Token模块实现
- [Framework层文档](../../framework/README.md)

---

**最后更新**: 2025-11-11

