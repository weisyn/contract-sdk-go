//go:build tinygo || (js && wasm)

package market

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Release 合约内分阶段释放操作
//
// 🎯 **用途**：在合约代码中创建分阶段释放计划
//
// **参数**：
//   - from: 释放者地址
//   - beneficiary: 受益人地址
//   - tokenID: 代币ID（nil表示原生币）
//   - totalAmount: 总释放金额
//   - vestingID: 释放计划ID（由合约生成）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 释放计划状态通过StateOutput记录
//   - 时间锁和释放条件验证是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func CreateVesting() uint32 {
//	    caller := framework.GetCaller()
//	    
//	    vestingID := generateVestingID(caller, beneficiary)
//	    
//	    err := market.Release(
//	        caller,
//	        beneficiary,
//	        nil,  // 原生币
//	        framework.Amount(100000),
//	        vestingID,
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Release(from, beneficiary framework.Address, tokenID framework.TokenID, totalAmount framework.Amount, vestingID []byte) error {
	// 1. 参数验证
	if err := validateReleaseParams(from, beneficiary, totalAmount, vestingID); err != nil {
		return err
	}

	// 2. 查询余额（通过framework）
	balance := framework.QueryUTXOBalance(from, tokenID)
	if balance < totalAmount {
		return framework.NewContractError(
			framework.ERROR_INSUFFICIENT_BALANCE,
			"insufficient balance to release",
		)
	}

	// 3. 构建释放计划状态ID
	stateID := buildVestingStateID(vestingID)

	// 4. 计算释放计划状态哈希
	execHash := computeVestingHash(stateID, from, beneficiary, totalAmount)

	// 5. 构建交易（使用internal包链式API）
	// 将代币转移到受益人地址（使用TimeLock或HeightLock）
	success, _, errCode := framework.BeginTransaction().
		Transfer(from, beneficiary, tokenID, totalAmount).
		AddStateOutput(stateID, 1, execHash).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "release failed")
	}

	// 6. 发出释放事件
	caller := framework.GetCaller()
	event := framework.NewEvent("Release")
	event.AddAddressField("from", from)
	event.AddAddressField("beneficiary", beneficiary)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("total_amount", uint64(totalAmount))
	event.AddField("vesting_id", string(vestingID))
	event.AddAddressField("caller", caller)
	framework.EmitEvent(event)

	return nil
}

// validateReleaseParams 验证释放参数
func validateReleaseParams(from, beneficiary framework.Address, totalAmount framework.Amount, vestingID []byte) error {
	// 验证地址
	zeroAddr := framework.Address{}
	if from == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"from address cannot be zero",
		)
	}
	if beneficiary == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"beneficiary address cannot be zero",
		)
	}
	if from == beneficiary {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"from and beneficiary addresses cannot be the same",
		)
	}

	// 验证金额
	if totalAmount == 0 {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"totalAmount must be greater than 0",
		)
	}

	// 验证释放计划ID
	if len(vestingID) == 0 {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"vestingID cannot be empty",
		)
	}

	return nil
}

// buildVestingStateID 构建释放计划状态ID
func buildVestingStateID(vestingID []byte) []byte {
	stateID := "vesting:" + string(vestingID)
	return []byte(stateID)
}

// computeVestingHash 计算释放计划状态哈希
// 使用framework.ComputeHash计算真实哈希值
func computeVestingHash(stateID []byte, from, beneficiary framework.Address, totalAmount framework.Amount) []byte {
	// 组合所有数据用于哈希计算
	data := make([]byte, 0, len(stateID)+40+8)
	data = append(data, stateID...)
	data = append(data, from.ToBytes()...)
	data = append(data, beneficiary.ToBytes()...)
	amountBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		amountBytes[i] = byte(totalAmount >> (i * 8))
	}
	data = append(data, amountBytes...)
	
	// 使用framework提供的真实哈希函数
	hash := framework.ComputeHash(data)
	return hash.ToBytes()
}

