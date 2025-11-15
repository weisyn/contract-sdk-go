//go:build tinygo || (js && wasm)

package token

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Approve 合约内代币授权操作
//
// 🎯 **用途**：授权指定地址使用代币
//
// **参数**：
//   - owner: 代币所有者地址
//   - spender: 被授权地址
//   - tokenID: 代币ID
//   - amount: 授权数量
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 授权信息需要存储在合约状态中
//   - 需要使用StateOutput来记录授权状态
//
// **示例**：
//
//	func Approve() uint32 {
//	    caller := framework.GetCaller()
//	    
//	    err := token.Approve(
//	        caller,
//	        spenderAddr,
//	        framework.TokenID("my_token"),
//	        framework.Amount(1000),
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Approve(owner, spender framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 1. 参数验证
	if err := validateApproveParams(owner, spender, tokenID, amount); err != nil {
		return err
	}

	// 2. 查询余额（通过framework）
	balance := framework.QueryUTXOBalance(owner, tokenID)
	if balance < amount {
		return framework.NewContractError(
			framework.ERROR_INSUFFICIENT_BALANCE,
			"insufficient balance to approve",
		)
	}

	// 3. 构建授权状态ID
	// 格式：approve:{owner}:{spender}:{tokenID}
	stateID := buildApproveStateID(owner, spender, tokenID)

	// 4. 计算授权状态哈希
	// 使用状态ID和金额构建哈希，用于StateOutput的execHash字段
	execHash := computeApproveHash(stateID, amount)

	// 5. 构建交易（使用internal包链式API）
	// 使用StateOutput记录授权状态
	success, _, errCode := framework.BeginTransaction().
		AddStateOutput(stateID, 1, execHash).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "approve failed")
	}

	// 6. 发出授权事件
	event := framework.NewEvent("Approve")
	event.AddAddressField("owner", owner)
	event.AddAddressField("spender", spender)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return nil
}

// validateApproveParams 验证授权参数
func validateApproveParams(owner, spender framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 验证地址
	zeroAddr := framework.Address{}
	if owner == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"owner address cannot be zero",
		)
	}
	if spender == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"spender address cannot be zero",
		)
	}
	if owner == spender {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"owner and spender addresses cannot be the same",
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

// buildApproveStateID 构建授权状态ID
func buildApproveStateID(owner, spender framework.Address, tokenID framework.TokenID) []byte {
	// 使用规范的格式构建状态ID
	stateID := "approve:" + owner.ToString() + ":" + spender.ToString() + ":" + string(tokenID)
	return []byte(stateID)
}

// computeApproveHash 计算授权状态哈希
// 使用framework.ComputeHash计算真实哈希值
func computeApproveHash(stateID []byte, amount framework.Amount) []byte {
	// 组合所有数据用于哈希计算
	data := make([]byte, 0, len(stateID)+8)
	data = append(data, stateID...)
	amountBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		amountBytes[i] = byte(amount >> (i * 8))
	}
	data = append(data, amountBytes...)
	
	// 使用framework提供的真实哈希函数
	hash := framework.ComputeHash(data)
	return hash.ToBytes()
}

