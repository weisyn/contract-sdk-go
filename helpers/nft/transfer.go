//go:build tinygo || (js && wasm)

package nft

import (
	"github.com/weisyn/contract-sdk-go/framework"
	"github.com/weisyn/contract-sdk-go/helpers/token"
)

// Transfer 合约内NFT转移操作
//
// 🎯 **用途**：在合约代码中转移NFT所有权
//
// **参数**：
//   - from: 发送者地址
//   - to: 接收者地址
//   - tokenID: NFT代币ID
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - NFT转移通过token.Transfer实现
//   - 权限控制是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func TransferNFT() uint32 {
//	    params := framework.GetContractParams()
//	    toStr := params.ParseJSON("to")
//	    tokenIDStr := params.ParseJSON("token_id")
//	    
//	    to, err := framework.ParseAddressBase58(toStr)
//	    if err != nil {
//	        return framework.ERROR_INVALID_PARAMS
//	    }
//	    
//	    caller := framework.GetCaller()
//	    err = nft.Transfer(caller, to, framework.TokenID(tokenIDStr))
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Transfer(from, to framework.Address, tokenID framework.TokenID) error {
	// 1. 参数验证
	if err := validateTransferParams(from, to, tokenID); err != nil {
		return err
	}

	// 2. 检查NFT所有者
	owner := OwnerOf(tokenID)
	if owner == nil {
		return framework.NewContractError(
			framework.ERROR_NOT_FOUND,
			"NFT not found",
		)
	}

	// 3. 验证发送者是所有者
	if !owner.Equal(from) {
		return framework.NewContractError(
			framework.ERROR_UNAUTHORIZED,
			"not the owner",
		)
	}

	// 4. 使用token.Transfer转移NFT（数量为1）
	err := token.Transfer(from, to, tokenID, framework.Amount(1))
	if err != nil {
		return err
	}

	// 5. 更新NFT所有权状态（使用StateOutput）
	{
		ownerStateID := buildOwnerStateID(tokenID)
		ownerHash := computeOwnerHash(ownerStateID, to)
		
		success, _, errCode := framework.BeginTransaction().
			AddStateOutput(ownerStateID, 1, ownerHash).
			Finalize()
		
		if !success {
			return framework.NewContractError(errCode, "failed to update owner")
		}
	}

	// 6. 发出NFT转移事件
	event := framework.NewEvent("NFTTransfer")
	event.AddAddressField("from", from)
	event.AddAddressField("to", to)
	event.AddStringField("token_id", string(tokenID))
	framework.EmitEvent(event)

	return nil
}

// validateTransferParams 验证转移参数
func validateTransferParams(from, to framework.Address, tokenID framework.TokenID) error {
	zeroAddr := framework.Address{}
	if from == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"from address cannot be zero",
		)
	}
	if to == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"to address cannot be zero",
		)
	}
	if from == to {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"from and to addresses cannot be the same",
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

