//go:build tinygo || (js && wasm)

package nft

import (
	"github.com/weisyn/contract-sdk-go/framework"
	"github.com/weisyn/contract-sdk-go/helpers/token"
)

// Burn 合约内NFT销毁操作
//
// 🎯 **用途**：在合约代码中销毁NFT
//
// **参数**：
//   - from: 销毁者地址
//   - tokenID: NFT代币ID
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - NFT销毁通过token.Burn实现
//   - 权限控制是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func BurnNFT() uint32 {
//	    params := framework.GetContractParams()
//	    tokenIDStr := params.ParseJSON("token_id")
//	    
//	    caller := framework.GetCaller()
//	    err = nft.Burn(caller, framework.TokenID(tokenIDStr))
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Burn(from framework.Address, tokenID framework.TokenID) error {
	// 1. 参数验证
	if err := validateBurnParams(from, tokenID); err != nil {
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

	// 3. 验证销毁者是所有者
	if !owner.Equal(from) {
		return framework.NewContractError(
			framework.ERROR_UNAUTHORIZED,
			"not the owner",
		)
	}

	// 4. 使用token.Burn销毁NFT（数量为1）
	err := token.Burn(from, tokenID, framework.Amount(1))
	if err != nil {
		return err
	}

	// 5. 发出NFT销毁事件
	event := framework.NewEvent("NFTBurn")
	event.AddAddressField("from", from)
	event.AddStringField("token_id", string(tokenID))
	framework.EmitEvent(event)

	return nil
}

// validateBurnParams 验证销毁参数
func validateBurnParams(from framework.Address, tokenID framework.TokenID) error {
	zeroAddr := framework.Address{}
	if from == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"from address cannot be zero",
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

