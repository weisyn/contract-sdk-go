//go:build tinygo || (js && wasm)

package token

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Freeze 合约内代币冻结操作
//
// 🎯 **用途**：冻结指定地址的代币
//
// **参数**：
//   - target: 目标地址
//   - tokenID: 代币ID
//   - amount: 冻结数量
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 冻结信息需要存储在合约状态中
//   - 需要使用StateOutput来记录冻结状态
//
// **示例**：
//
//	func Freeze() uint32 {
//	    caller := framework.GetCaller()
//	    
//	    // 权限检查（业务逻辑）
//	    if !isAuthorizedFreezer(caller) {
//	        return framework.ERROR_UNAUTHORIZED
//	    }
//	    
//	    err := token.Freeze(
//	        targetAddr,
//	        framework.TokenID("my_token"),
//	        framework.Amount(1000),
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Freeze(target framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 1. 参数验证
	if err := validateFreezeParams(target, tokenID, amount); err != nil {
		return err
	}

	// 2. 查询余额（通过framework）
	balance := framework.QueryUTXOBalance(target, tokenID)
	if balance < amount {
		return framework.NewContractError(
			framework.ERROR_INSUFFICIENT_BALANCE,
			"insufficient balance to freeze",
		)
	}

	// 3. 构建冻结状态ID
	stateID := buildFreezeStateID(target, tokenID)

	// 4. 计算冻结状态哈希
	execHash := computeFreezeHash(stateID, amount)

	// 5. 构建交易（使用internal包链式API）
	// 使用StateOutput记录冻结状态
	success, _, errCode := framework.BeginTransaction().
		AddStateOutput(stateID, 1, execHash).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "freeze failed")
	}

	// 6. 发出冻结事件
	caller := framework.GetCaller()
	event := framework.NewEvent("Freeze")
	event.AddAddressField("target", target)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	event.AddAddressField("freezer", caller)
	framework.EmitEvent(event)

	return nil
}

// validateFreezeParams 验证冻结参数
func validateFreezeParams(target framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 验证地址
	zeroAddr := framework.Address{}
	if target == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"target address cannot be zero",
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

// buildFreezeStateID 构建冻结状态ID
func buildFreezeStateID(target framework.Address, tokenID framework.TokenID) []byte {
	stateID := "freeze:" + target.ToString() + ":" + string(tokenID)
	return []byte(stateID)
}

// computeFreezeHash 计算冻结状态哈希
func computeFreezeHash(stateID []byte, amount framework.Amount) []byte {
	hash := make([]byte, 32)
	copy(hash, stateID)
	if len(hash) > 32 {
		hash = hash[:32]
	}
	return hash
}

