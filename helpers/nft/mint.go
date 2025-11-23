//go:build tinygo || (js && wasm)

package nft

import (
	"github.com/weisyn/contract-sdk-go/framework"
	"github.com/weisyn/contract-sdk-go/helpers/token"
)

// Mint 合约内NFT铸造操作
//
// 🎯 **用途**：在合约代码中铸造NFT
//
// **参数**：
//   - to: 接收者地址
//   - tokenID: NFT代币ID（必须唯一）
//   - metadata: NFT元数据（可选）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - NFT铸造通过token.Mint实现，数量固定为1
//   - 权限控制和元数据格式验证是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func MintNFT() uint32 {
//	    params := framework.GetContractParams()
//	    toStr := params.ParseJSON("to")
//	    tokenIDStr := params.ParseJSON("token_id")
//	    
//	    to, err := framework.ParseAddressBase58(toStr)
//	    if err != nil {
//	        return framework.ERROR_INVALID_PARAMS
//	    }
//	    
//	    err = nft.Mint(
//	        to,
//	        framework.TokenID(tokenIDStr),
//	        []byte(params.ParseJSON("metadata")),
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Mint(to framework.Address, tokenID framework.TokenID, metadata []byte) error {
	// 1. 参数验证
	if err := validateMintParams(to, tokenID); err != nil {
		return err
	}

	// 2. 检查NFT是否已存在
	owner := OwnerOf(tokenID)
	if owner != nil {
		return framework.NewContractError(
			framework.ERROR_ALREADY_EXISTS,
			"NFT already exists",
		)
	}

	// 3. 使用token.Mint铸造NFT（数量为1）
	err := token.Mint(to, tokenID, framework.Amount(1))
	if err != nil {
		return err
	}

	// 4. 存储NFT元数据（使用StateOutput）
	if len(metadata) > 0 {
		stateID := buildMetadataStateID(tokenID)
		execHash := computeMetadataHash(stateID, metadata)
		
		success, _, errCode := framework.BeginTransaction().
			AddStateOutput(stateID, 1, execHash).
			Finalize()
		
		if !success {
			return framework.NewContractError(errCode, "failed to store metadata")
		}
	}

	// 5. 发出NFT铸造事件
	caller := framework.GetCaller()
	event := framework.NewEvent("NFTMint")
	event.AddAddressField("to", to)
	event.AddStringField("token_id", string(tokenID))
	event.AddAddressField("minter", caller)
	if len(metadata) > 0 {
		event.AddField("metadata", string(metadata))
	}
	framework.EmitEvent(event)

	return nil
}

// validateMintParams 验证铸造参数
func validateMintParams(to framework.Address, tokenID framework.TokenID) error {
	zeroAddr := framework.Address{}
	if to == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"to address cannot be zero",
		)
	}

	if tokenID == "" {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"tokenID cannot be empty",
		)
	}

	return nil
}

// buildMetadataStateID 构建元数据状态ID
func buildMetadataStateID(tokenID framework.TokenID) []byte {
	stateID := "nft_metadata:" + string(tokenID)
	return []byte(stateID)
}

// computeMetadataHash 计算元数据哈希
func computeMetadataHash(stateID []byte, metadata []byte) []byte {
	data := make([]byte, 0, len(stateID)+len(metadata))
	data = append(data, stateID...)
	data = append(data, metadata...)
	
	hash := framework.ComputeHash(data)
	return hash.ToBytes()
}

