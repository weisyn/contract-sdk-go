//go:build tinygo || (js && wasm)

package token

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Burn 合约内代币销毁操作
//
// 🎯 **用途**：在合约代码中销毁代币
//
// **参数**：
//   - from: 销毁者地址
//   - tokenID: 代币ID
//   - amount: 销毁数量
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 销毁操作需要先消费UTXO，然后不创建新的输出
//   - 权限控制是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func Burn() uint32 {
//	    caller := framework.GetCaller()
//	    
//	    err := token.Burn(
//	        caller,
//	        framework.TokenID("my_token"),
//	        framework.Amount(500),
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Burn(from framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 1. 参数验证
	if err := validateBurnParams(from, tokenID, amount); err != nil {
		return err
	}

	// 2. 查询余额（通过framework）
	balance := framework.QueryUTXOBalance(from, tokenID)
	if balance < amount {
		return framework.NewContractError(
			framework.ERROR_INSUFFICIENT_BALANCE,
			"insufficient balance to burn",
		)
	}

	// 3. 构建交易（使用framework链式API）
	// 注意：在UTXO模型中，销毁代币的标准方式是将其转移到零地址
	// 零地址是一个特殊的地址，代币一旦转移到零地址，就无法再被使用
	// 这是UTXO模型中的标准销毁方式，符合区块链的去中心化原则
	zeroAddr := framework.Address{}
	success, _, errCode := framework.BeginTransaction().
		Transfer(from, zeroAddr, tokenID, amount).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "burn failed")
	}

	// 4. 发出销毁事件
	event := framework.NewEvent("Burn")
	event.AddAddressField("from", from)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return nil
}

// validateBurnParams 验证销毁参数
func validateBurnParams(from framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 验证地址
	zeroAddr := framework.Address{}
	if from == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"from address cannot be zero",
		)
	}

	// 验证代币ID
	if tokenID == "" {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"tokenID cannot be empty",
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

