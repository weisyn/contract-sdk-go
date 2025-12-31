//go:build tinygo || (js && wasm)

package nft

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// OwnerOf 查询NFT所有者
//
// 🎯 **用途**：查询指定NFT的所有者地址
//
// **参数**：
//   - tokenID: NFT代币ID
//
// **返回**：
//   - owner: 所有者地址，nil表示NFT不存在
//
// **实现说明**：
//   通过查询 nft_owner:{tokenID} 状态获取所有者地址
func OwnerOf(tokenID framework.TokenID) *framework.Address {
	if tokenID == "" {
		return nil
	}
	
	// 构建所有权状态ID
	ownerStateID := buildOwnerStateID(tokenID)
	
	// 查询链上状态
	stateData := framework.GetStateFromChain(string(ownerStateID))
	if stateData == nil || len(stateData) == 0 {
		return nil
	}
	
	// 解析所有者地址（stateData 包含地址字节）
	// 地址长度为 framework.AddressLen (通常是 20 或 32 字节)
	if len(stateData) < framework.AddressLen {
		return nil
	}
	
	var addr framework.Address
	copy(addr[:], stateData[:framework.AddressLen])
	return &addr
}

// BalanceOf 查询地址拥有的NFT数量
//
// 🎯 **用途**：查询指定地址拥有的NFT数量
//
// **参数**：
//   - owner: 所有者地址
//
// **返回**：
//   - count: NFT数量
//
// **实现说明**：
//   在EUTXO模型中，通过查询该地址的所有UTXO，统计数量为1的tokenID数量。
func BalanceOf(owner framework.Address) uint64 {
	// 在EUTXO模型中，NFT数量通过查询UTXO余额来确定
	// 这里简化实现：返回0
	// 实际应用中，应该查询该地址的所有UTXO，统计数量为1的tokenID数量
	
	// 注意：这是一个简化实现
	// 实际应用中，应该实现完整的查询逻辑
	return 0
}

// GetMetadata 获取NFT元数据
//
// 🎯 **用途**：查询指定NFT的元数据
//
// **参数**：
//   - tokenID: NFT代币ID
//
// **返回**：
//   - metadata: 元数据，nil表示元数据不存在
//
// **实现说明**：
//   通过查询 nft_metadata:{tokenID} 状态获取元数据
func GetMetadata(tokenID framework.TokenID) []byte {
	if tokenID == "" {
		return nil
	}
	
	// 构建元数据状态ID
	stateID := buildMetadataStateID(tokenID)
	
	// 查询链上状态
	metadata := framework.GetStateFromChain(string(stateID))
	if metadata == nil || len(metadata) == 0 {
		return nil
	}
	
	return metadata
}

