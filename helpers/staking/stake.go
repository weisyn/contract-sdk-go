//go:build tinygo || (js && wasm)

package staking

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Stake 合约内质押操作
//
// 🎯 **用途**：在合约代码中执行质押
//
// **参数**：
//   - staker: 质押者地址
//   - validator: 验证者地址
//   - tokenID: 代币ID（nil表示原生币）
//   - amount: 质押金额
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 质押操作会创建带ContractLock的UTXO输出
//   - 质押状态通过StateOutput记录
//   - 权限控制和锁定期管理是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func Stake() uint32 {
//	    caller := framework.GetCaller()
//	    
//	    // 权限检查（业务逻辑）
//	    if !isAuthorizedStaker(caller) {
//	        return framework.ERROR_UNAUTHORIZED
//	    }
//	    
//	    err := staking.Stake(
//	        caller,
//	        validatorAddr,
//	        nil,  // 原生币
//	        framework.Amount(10000),
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Stake(staker, validator framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 1. 参数验证
	if err := validateStakeParams(staker, validator, amount); err != nil {
		return err
	}

	// 2. 查询余额（通过framework）
	balance := framework.QueryUTXOBalance(staker, tokenID)
	if balance < amount {
		return framework.NewContractError(
			framework.ERROR_INSUFFICIENT_BALANCE,
			"insufficient balance to stake",
		)
	}

	// 3. 构建交易（使用internal包链式API）
	// 质押操作：将代币转移到验证者地址，并添加ContractLock
	success, _, errCode := framework.BeginTransaction().
		Stake(staker, amount, validator).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "stake failed")
	}

	// 4. 发出质押事件
	caller := framework.GetCaller()
	event := framework.NewEvent("Stake")
	event.AddAddressField("staker", staker)
	event.AddAddressField("validator", validator)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	event.AddAddressField("caller", caller)
	framework.EmitEvent(event)

	return nil
}

// validateStakeParams 验证质押参数
func validateStakeParams(staker, validator framework.Address, amount framework.Amount) error {
	// 验证地址
	zeroAddr := framework.Address{}
	if staker == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"staker address cannot be zero",
		)
	}
	if validator == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"validator address cannot be zero",
		)
	}
	if staker == validator {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"staker and validator addresses cannot be the same",
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

