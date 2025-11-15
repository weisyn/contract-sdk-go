//go:build tinygo || (js && wasm)

package token

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// MintRecipient 批量铸造接收者
type MintRecipient struct {
	Address framework.Address
	Amount  framework.Amount
}

// BatchMint 批量铸造代币
//
// 🎯 **用途**：一次性向多个地址铸造代币
//
// **参数**：
//   - recipients: 接收者列表，每个接收者包含地址和数量
//   - tokenID: 代币ID（空字符串表示原生币或合约代币）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 合约只能铸造自己的代币
//   - 权限控制和总量控制是业务逻辑，需要在合约代码中实现
//   - 批量铸造会在一次交易中创建多个AssetOutput（UTXO输出）
//
// **示例**：
//
//	func BatchMint() uint32 {
//	    recipients := []token.MintRecipient{
//	        {Address: addr1, Amount: framework.Amount(100)},
//	        {Address: addr2, Amount: framework.Amount(200)},
//	        {Address: addr3, Amount: framework.Amount(300)},
//	    }
//	    
//	    err := token.BatchMint(recipients, framework.TokenID("my_token"))
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func BatchMint(recipients []MintRecipient, tokenID framework.TokenID) error {
	// 1. 参数验证
	if err := validateBatchMintParams(recipients, tokenID); err != nil {
		return err
	}

	// 2. 构建交易（使用internal包链式API）
	// 注意：批量铸造操作实际上是创建多个UTXO输出
	builder := framework.BeginTransaction()

	// 为每个接收者创建AssetOutput
	for _, recipient := range recipients {
		builder.AddAssetOutput(recipient.Address, tokenID, recipient.Amount)
	}

	// 完成交易构建
	success, _, errCode := builder.Finalize()
	if !success {
		return framework.NewContractError(errCode, "batch mint failed")
	}

	// 3. 发出批量铸造事件
	caller := framework.GetCaller()
	event := framework.NewEvent("BatchMint")
	event.AddAddressField("minter", caller)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("recipient_count", uint64(len(recipients)))
	
	// 计算总金额
	var totalAmount framework.Amount
	for _, recipient := range recipients {
		totalAmount = totalAmount.Add(recipient.Amount)
	}
	event.AddUint64Field("total_amount", uint64(totalAmount))
	
	framework.EmitEvent(event)

	return nil
}

// validateBatchMintParams 验证批量铸造参数
func validateBatchMintParams(recipients []MintRecipient, tokenID framework.TokenID) error {
	// 验证接收者列表
	if len(recipients) == 0 {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"recipients list cannot be empty",
		)
	}

	// 验证每个接收者
	zeroAddr := framework.Address{}
	for i, recipient := range recipients {
		if recipient.Address == zeroAddr {
			return framework.NewContractError(
				framework.ERROR_INVALID_PARAMS,
				"recipient address cannot be zero",
			)
		}
		if recipient.Amount == 0 {
			return framework.NewContractError(
				framework.ERROR_INVALID_PARAMS,
				"recipient amount must be greater than 0",
			)
		}
		
		// 检查重复地址（可选，但建议避免）
		for j := i + 1; j < len(recipients); j++ {
			if recipient.Address == recipients[j].Address {
				return framework.NewContractError(
					framework.ERROR_INVALID_PARAMS,
					"duplicate recipient address",
				)
			}
		}
	}

	return nil
}

