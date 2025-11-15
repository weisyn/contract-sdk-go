//go:build tinygo || (js && wasm)

package token

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// AirdropRecipient 空投接收者
type AirdropRecipient struct {
	Address framework.Address
	Amount  framework.Amount
}

// Airdrop 合约内批量空投操作
//
// 🎯 **用途**：批量转账到多个地址
//
// **参数**：
//   - from: 发送者地址
//   - recipients: 接收者列表
//   - tokenID: 代币ID（nil表示原生币）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **示例**：
//
//	func Airdrop() uint32 {
//	    caller := framework.GetCaller()
//	    
//	    recipients := []token.AirdropRecipient{
//	        {Address: addr1, Amount: framework.Amount(100)},
//	        {Address: addr2, Amount: framework.Amount(200)},
//	    }
//	    
//	    err := token.Airdrop(caller, recipients, framework.TokenID("my_token"))
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Airdrop(from framework.Address, recipients []AirdropRecipient, tokenID framework.TokenID) error {
	// 1. 参数验证
	if err := validateAirdropParams(from, recipients, tokenID); err != nil {
		return err
	}

	// 2. 计算总金额
	var totalAmount framework.Amount
	for _, recipient := range recipients {
		totalAmount = totalAmount.Add(recipient.Amount)
	}

	// 3. 查询余额（通过framework）
	balance := framework.QueryUTXOBalance(from, tokenID)
	if balance < totalAmount {
		return framework.NewContractError(
			framework.ERROR_INSUFFICIENT_BALANCE,
			"insufficient balance for airdrop",
		)
	}

	// 4. 构建交易（使用internal包链式API）
	builder := framework.BeginTransaction()

	// 添加所有接收者的输出
	for _, recipient := range recipients {
		builder.AddAssetOutput(recipient.Address, tokenID, recipient.Amount)
	}

	// 完成交易构建
	success, _, errCode := builder.Finalize()
	if !success {
		return framework.NewContractError(errCode, "airdrop failed")
	}

	// 5. 发出空投事件
	event := framework.NewEvent("Airdrop")
	event.AddAddressField("from", from)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("total_amount", uint64(totalAmount))
	event.AddUint64Field("recipient_count", uint64(len(recipients)))
	framework.EmitEvent(event)

	return nil
}

// validateAirdropParams 验证空投参数
func validateAirdropParams(from framework.Address, recipients []AirdropRecipient, tokenID framework.TokenID) error {
	// 验证发送者地址
	zeroAddr := framework.Address{}
	if from == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"from address cannot be zero",
		)
	}

	// 验证接收者列表
	if len(recipients) == 0 {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"recipients list cannot be empty",
		)
	}

	// 验证每个接收者
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
		// 检查重复地址（可选）
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

