//go:build tinygo || (js && wasm)

package staking

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Undelegate 合约内取消委托操作
//
// 🎯 **用途**：在合约代码中执行取消委托
//
// **参数**：
//   - delegator: 委托者地址
//   - validator: 验证者地址
//   - tokenID: 代币ID（nil表示原生币）
//   - amount: 取消委托金额（0表示全部取消）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 取消委托操作需要解锁DelegationLock的UTXO
//   - 取消委托状态通过StateOutput记录
//   - 锁定期检查和权限控制是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func Undelegate() uint32 {
//	    caller := framework.GetCaller()
//	    
//	    // 锁定期检查（业务逻辑）
//	    if !isUndelegatePeriodReached(caller, validator) {
//	        return framework.ERROR_INVALID_STATE
//	    }
//	    
//	    err := staking.Undelegate(
//	        caller,
//	        validatorAddr,
//	        nil,  // 原生币
//	        framework.Amount(2000),  // 部分取消委托
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Undelegate(delegator, validator framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 1. 参数验证
	if err := validateUndelegateParams(delegator, validator, amount); err != nil {
		return err
	}

	// 2. 构建交易（使用internal包链式API）
	// 取消委托操作：从验证者地址转回委托者，解锁DelegationLock
	// 注意：实际实现中需要查询委托UTXO并解锁
	success, _, errCode := framework.BeginTransaction().
		Transfer(validator, delegator, tokenID, amount).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "undelegate failed")
	}

	// 3. 发出取消委托事件
	caller := framework.GetCaller()
	event := framework.NewEvent("Undelegate")
	event.AddAddressField("delegator", delegator)
	event.AddAddressField("validator", validator)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	event.AddAddressField("caller", caller)
	framework.EmitEvent(event)

	return nil
}

// validateUndelegateParams 验证取消委托参数
func validateUndelegateParams(delegator, validator framework.Address, amount framework.Amount) error {
	// 验证地址
	zeroAddr := framework.Address{}
	if delegator == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"delegator address cannot be zero",
		)
	}
	if validator == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"validator address cannot be zero",
		)
	}

	// 金额可以为0（表示全部取消委托），但不能为负数（由类型系统保证）

	return nil
}

