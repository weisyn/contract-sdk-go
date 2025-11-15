//go:build tinygo || (js && wasm)

package token

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Mint 合约内代币铸造操作
//
// 🎯 **用途**：在合约代码中铸造新代币
//
// **参数**：
//   - to: 接收者地址
//   - tokenID: 代币ID
//   - amount: 铸造数量
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 合约只能铸造自己的代币
//   - 权限控制和总量控制是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func Mint() uint32 {
//	    caller := framework.GetCaller()
//	    contractAddr := framework.GetContractAddress()
//	    
//	    // 权限检查（业务逻辑）
//	    if !isAuthorizedMinter(caller) {
//	        return framework.ERROR_UNAUTHORIZED
//	    }
//	    
//	    err := token.Mint(
//	        recipientAddr,
//	        framework.TokenID("my_token"),
//	        framework.Amount(1000),
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Mint(to framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 1. 参数验证
	if err := validateMintParams(to, tokenID, amount); err != nil {
		return err
	}

	// 2. 构建交易（使用internal包链式API）
	// 注意：Mint操作实际上是创建新的UTXO输出
	success, _, errCode := framework.BeginTransaction().
		AddAssetOutput(to, tokenID, amount).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "mint failed")
	}

	// 3. 发出铸造事件
	caller := framework.GetCaller()
	event := framework.NewEvent("Mint")
	event.AddAddressField("to", to)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	event.AddAddressField("minter", caller)
	framework.EmitEvent(event)

	return nil
}

// validateMintParams 验证铸造参数
func validateMintParams(to framework.Address, tokenID framework.TokenID, amount framework.Amount) error {
	// 验证地址
	zeroAddr := framework.Address{}
	if to == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"to address cannot be zero",
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

