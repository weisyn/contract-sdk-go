//go:build tinygo || (js && wasm)

package staking

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Unstake 合约内解质押操作
//
// 🎯 **用途**：在合约代码中执行解质押
//
// **参数**：
//   - staker: 质押者地址
//   - validator: 验证者地址
//   - tokenID: 代币ID（nil表示原生币）
//   - amount: 解质押金额（0表示全部解质押）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 解质押操作需要解锁ContractLock的UTXO
//   - 解质押状态通过StateOutput记录
//   - 锁定期检查和权限控制是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func Unstake() uint32 {
//	    caller := framework.GetCaller()
//	    
//	    // 锁定期检查（业务逻辑）
//	    if !isUnlockPeriodReached(caller, validator) {
//	        return framework.ERROR_INVALID_STATE
//	    }
//	    
//	    err := staking.Unstake(
//	        caller,
//	        validatorAddr,
//	        nil,  // 原生币
//	        framework.Amount(5000),  // 部分解质押
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Unstake(staker, validator framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 1. 参数验证
	if err := validateUnstakeParams(staker, validator, amount); err != nil {
		return err
	}

	// 2. 构建交易（使用internal包链式API）
	// 解质押操作：从验证者地址转回质押者，解锁ContractLock
	// 注意：实际实现中需要查询质押UTXO并解锁
	success, _, errCode := framework.BeginTransaction().
		Transfer(validator, staker, tokenID, amount).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "unstake failed")
	}

	// 3. 发出解质押事件
	caller := framework.GetCaller()
	event := framework.NewEvent("Unstake")
	event.AddAddressField("staker", staker)
	event.AddAddressField("validator", validator)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	event.AddAddressField("caller", caller)
	framework.EmitEvent(event)

	return nil
}

// validateUnstakeParams 验证解质押参数
func validateUnstakeParams(staker, validator framework.Address, amount framework.Amount) error {
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

	// 金额可以为0（表示全部解质押），但不能为负数（由类型系统保证）

	return nil
}

