//go:build tinygo || (js && wasm)

package market

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Escrow 合约内托管操作
//
// 🎯 **用途**：在合约代码中创建托管
//
// **参数**：
//   - buyer: 买方地址
//   - seller: 卖方地址
//   - tokenID: 代币ID（nil表示原生币）
//   - amount: 托管金额
//   - escrowID: 托管ID（由合约生成）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 托管状态通过StateOutput记录
//   - 权限控制和托管条件验证是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func CreateEscrow() uint32 {
//	    caller := framework.GetCaller()
//	    
//	    escrowID := generateEscrowID(caller, seller)
//	    
//	    err := market.Escrow(
//	        buyer,
//	        seller,
//	        nil,  // 原生币
//	        framework.Amount(10000),
//	        escrowID,
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Escrow(buyer, seller framework.Address, tokenID framework.TokenID, amount framework.Amount, escrowID []byte) error {
	// 1. 参数验证
	if err := validateEscrowParams(buyer, seller, amount, escrowID); err != nil {
		return err
	}

	// 2. 查询余额（通过framework）
	balance := framework.QueryUTXOBalance(buyer, tokenID)
	if balance < amount {
		return framework.NewContractError(
			framework.ERROR_INSUFFICIENT_BALANCE,
			"insufficient balance to escrow",
		)
	}

	// 3. 构建托管状态ID
	stateID := buildEscrowStateID(escrowID)

	// 4. 计算托管状态哈希
	execHash := computeEscrowHash(stateID, buyer, seller, amount)

	// 5. 构建交易（使用internal包链式API）
	// 将代币转移到托管地址（使用ContractLock）
	success, _, errCode := framework.BeginTransaction().
		Transfer(buyer, seller, tokenID, amount).
		AddStateOutput(stateID, 1, execHash).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "escrow failed")
	}

	// 6. 发出托管事件
	caller := framework.GetCaller()
	event := framework.NewEvent("Escrow")
	event.AddAddressField("buyer", buyer)
	event.AddAddressField("seller", seller)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	event.AddField("escrow_id", string(escrowID))
	event.AddAddressField("caller", caller)
	framework.EmitEvent(event)

	return nil
}

// validateEscrowParams 验证托管参数
func validateEscrowParams(buyer, seller framework.Address, amount framework.Amount, escrowID []byte) error {
	zeroAddr := framework.Address{}
	if buyer == zeroAddr {
		return framework.NewContractError(framework.ERROR_INVALID_PARAMS, "buyer address cannot be zero")
	}
	if seller == zeroAddr {
		return framework.NewContractError(framework.ERROR_INVALID_PARAMS, "seller address cannot be zero")
	}
	if buyer == seller {
		return framework.NewContractError(framework.ERROR_INVALID_PARAMS, "buyer and seller addresses cannot be the same")
	}
	if amount == 0 {
		return framework.NewContractError(framework.ERROR_INVALID_PARAMS, "amount must be greater than 0")
	}
	if len(escrowID) == 0 {
		return framework.NewContractError(framework.ERROR_INVALID_PARAMS, "escrowID cannot be empty")
	}
	return nil
}

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
// Release 函数已移至 release.go，请使用 market.Release()

// buildEscrowStateID 构建托管状态ID
func buildEscrowStateID(escrowID []byte) []byte {
	stateID := "escrow:" + string(escrowID)
	return []byte(stateID)
}

// computeEscrowHash 计算托管状态哈希
// 使用framework.ComputeHash计算真实哈希值
func computeEscrowHash(stateID []byte, buyer, seller framework.Address, amount framework.Amount) []byte {
	// 组合所有数据用于哈希计算
	data := make([]byte, 0, len(stateID)+40+8)
	data = append(data, stateID...)
	data = append(data, buyer.ToBytes()...)
	data = append(data, seller.ToBytes()...)
	amountBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		amountBytes[i] = byte(amount >> (i * 8))
	}
	data = append(data, amountBytes...)
	
	// 使用framework提供的真实哈希函数
	hash := framework.ComputeHash(data)
	return hash.ToBytes()
}

