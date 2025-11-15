//go:build tinygo || (js && wasm)

package staking

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Delegate 合约内委托操作
//
// 🎯 **用途**：在合约代码中执行委托
//
// **参数**：
//   - delegator: 委托者地址
//   - validator: 验证者地址
//   - tokenID: 代币ID（nil表示原生币）
//   - amount: 委托金额
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 委托操作使用DelegationLock
//   - 委托状态通过StateOutput记录
//   - 权限控制和委托限制是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func Delegate() uint32 {
//	    caller := framework.GetCaller()
//	    
//	    // 权限检查（业务逻辑）
//	    if !isAuthorizedDelegator(caller) {
//	        return framework.ERROR_UNAUTHORIZED
//	    }
//	    
//	    err := staking.Delegate(
//	        caller,
//	        validatorAddr,
//	        nil,  // 原生币
//	        framework.Amount(5000),
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Delegate(delegator, validator framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 1. 参数验证
	if err := validateDelegateParams(delegator, validator, amount); err != nil {
		return err
	}

	// 2. 查询余额（通过framework）
	balance := framework.QueryUTXOBalance(delegator, tokenID)
	if balance < amount {
		return framework.NewContractError(
			framework.ERROR_INSUFFICIENT_BALANCE,
			"insufficient balance to delegate",
		)
	}

	// 3. 构建交易（使用TransactionBuilder实现完整转账逻辑）
	// 委托操作：将代币转移到验证者地址，并添加DelegationLock
	// DelegationLock允许验证者代为操作委托的代币，但受到授权限制
	// 注意：这里使用标准的Transfer，DelegationLock应该在合约的业务逻辑中通过StateOutput记录
	// 实际的锁定条件应用应该在合约的业务逻辑中处理，而不是在helpers层
	// TransactionBuilder.Transfer() 会通过 host_build_transaction 处理UTXO选择和找零
	success, _, errCode := framework.BeginTransaction().
		Transfer(delegator, validator, tokenID, amount).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "delegate failed")
	}

	// 4. 发出委托事件
	caller := framework.GetCaller()
	event := framework.NewEvent("Delegate")
	event.AddAddressField("delegator", delegator)
	event.AddAddressField("validator", validator)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	event.AddAddressField("caller", caller)
	framework.EmitEvent(event)

	return nil
}

// validateDelegateParams 验证委托参数
func validateDelegateParams(delegator, validator framework.Address, amount framework.Amount) error {
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
	if delegator == validator {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"delegator and validator addresses cannot be the same",
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

