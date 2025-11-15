//go:build tinygo || (js && wasm)

package token

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Transfer 合约内转账操作
//
// 🎯 **用途**：在合约代码中执行转账
//
// **参数**：
//   - from: 发送者地址
//   - to: 接收者地址
//   - tokenID: 代币ID（nil表示原生币）
//   - amount: 转账金额
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **示例**：
//
//	func Transfer() uint32 {
//	    err := token.Transfer(
//	        framework.GetCaller(),
//	        recipientAddr,
//	        nil,  // 原生币
//	        framework.Amount(100),
//	    )
//	    if err != nil {
//	        return framework.ERROR_INSUFFICIENT_BALANCE
//	    }
//	    return framework.SUCCESS
//	}
func Transfer(from, to framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 1. 参数验证
	if err := validateTransferParams(from, to, amount); err != nil {
		return err
	}

	// 2. 查询余额（通过framework）
	balance := framework.QueryUTXOBalance(from, tokenID)
	if balance < amount {
		return framework.NewContractError(
			framework.ERROR_INSUFFICIENT_BALANCE,
			"insufficient balance",
		)
	}

	// 3. 构建交易（使用internal包链式API）
	success, _, errCode := framework.BeginTransaction().
		Transfer(from, to, tokenID, amount).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "transfer failed")
	}

	// 4. 发出转账事件
	event := framework.NewEvent("Transfer")
	event.AddAddressField("from", from)
	event.AddAddressField("to", to)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return nil
}

// validateTransferParams 验证转账参数
func validateTransferParams(from, to framework.Address, amount framework.Amount) error {
	// 验证地址
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

	// 验证金额
	if amount == 0 {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"amount must be greater than 0",
		)
	}

	return nil
}

