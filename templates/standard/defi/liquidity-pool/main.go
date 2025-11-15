//go:build tinygo || (js && wasm)

// Package main 提供流动性池合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建流动性池合约。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/token 模块进行代币操作
//   - 如何使用 helpers/market 模块进行资产托管
//   - 如何实现完整的流动性池功能（添加流动性、移除流动性、查询池信息）
//
// 🎯 核心功能
//
//  1. AddLiquidity - 添加流动性
//     - 向流动性池添加代币
//     - 获得流动性凭证代币（LP Token）
//
//  2. RemoveLiquidity - 移除流动性
//     - 从流动性池移除代币
//     - 销毁流动性凭证代币
//
//  3. QueryPoolInfo - 查询池信息
//     - 查询流动性池的详细信息
//     - 查询池中代币余额和LP Token总量
//
// ⚠️ 注意：本示例是简化实现
//   实际应用中需要实现：
//   - 流动性份额计算
//   - 收益分配机制
//   - 流动性凭证代币管理
//
// 📚 相关文档
//
//   - [Token 模块文档](../../../helpers/token/README.md)
//   - [Market 模块文档](../../../helpers/market/README.md)
//   - [Framework 文档](../../../framework/README.md)
//   - [示例总览](../../README.md)
package main

import (
	"github.com/weisyn/contract-sdk-go/helpers/token"
	"github.com/weisyn/contract-sdk-go/framework"
)

// LiquidityPoolContract 流动性池合约
//
// 本合约使用 helpers/token 和 helpers/market 模块提供的业务语义API，
// 简化流动性池操作的实现，开发者只需关注业务逻辑。
//
// 流动性池特点：
//   - 用户存入代币，获得流动性凭证代币（LP Token）
//   - LP Token代表用户在池中的份额
//   - 流动性提供者获得收益分成
type LiquidityPoolContract struct {
	framework.ContractBase
}

// Initialize 初始化合约
//
// 合约部署时自动调用，用于初始化合约状态。
//
// 工作流程：
//  1. 获取合约调用者（部署者）
//  2. 发出合约初始化事件
//
// 返回：
//   - framework.SUCCESS - 初始化成功
//
// 事件：
//   - ContractInitialized - 合约初始化事件
//     {
//       "contract": "LiquidityPool",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "LiquidityPool")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// AddLiquidity 添加流动性
//
// 向流动性池添加代币，获得流动性凭证代币（LP Token）。
// SDK 内部会自动处理：
//   - 余额检查（确保用户余额充足）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 AddLiquidity 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "token_id": "TOKEN_001",  // 代币ID（可选，nil表示原生代币）
//	  "amount": 10000           // 添加数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 检查用户余额
//  3. 计算流动性份额（根据池中总代币和LP Token总量）
//  4. 转移代币到合约
//  5. 铸造流动性凭证代币
//  6. 发出添加流动性事件
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中需要实现：
//   - 流动性份额计算（LP Token数量 = (amount / totalReserve) * totalLPTokens）
//   - 首次添加流动性的特殊处理
//   - 流动性凭证代币的铸造
//
// 返回：
//   - framework.SUCCESS - 添加成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - AddLiquidity - 添加流动性事件
//     {
//       "provider": "<流动性提供者地址>",
//       "token_id": "TOKEN_001",
//       "amount": 10000,
//       "lp_token_amount": 100
//     }
//
//export AddLiquidity
func AddLiquidity() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	tokenIDStr := params.ParseJSON("token_id")
	amount := params.ParseJSONInt("amount")

	if amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析代币ID（可选）
	var tokenID framework.TokenID
	if tokenIDStr != "" {
		tokenID = framework.TokenID(tokenIDStr)
	}

	// 步骤3：获取调用者
	caller := framework.GetCaller()

	// 步骤4：检查余额
	balance := framework.QueryUTXOBalance(caller, tokenID)
	if balance < framework.Amount(amount) {
		return framework.ERROR_INSUFFICIENT_BALANCE
	}

	// 步骤5：计算流动性份额
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该计算：
	//   LP Token数量 = (amount / totalReserve) * totalLPTokens
	//   首次添加流动性时，LP Token数量 = amount
	lpTokenAmount := amount / 100

	// 步骤6：转移代币到合约
	contractAddr := framework.GetContractAddress()
	err := token.Transfer(
		caller,
		contractAddr,
		tokenID,
		framework.Amount(amount),
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤7：铸造流动性凭证代币
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该铸造流动性凭证代币（LP Token）给用户
	//   这里简化处理，不实际铸造

	// 步骤8：发出添加流动性事件
	event := framework.NewEvent("AddLiquidity")
	event.AddAddressField("provider", caller)
	if tokenIDStr != "" {
		event.AddStringField("token_id", tokenIDStr)
	}
	event.AddUint64Field("amount", uint64(amount))
	event.AddUint64Field("lp_token_amount", uint64(lpTokenAmount))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// RemoveLiquidity 移除流动性
//
// 从流动性池移除代币，销毁流动性凭证代币。
// SDK 内部会自动处理：
//   - 余额检查（确保用户有足够的LP Token）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 RemoveLiquidity 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "token_id": "TOKEN_001",  // 代币ID（可选，nil表示原生代币）
//	  "lp_token_amount": 100    // LP Token数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 检查LP Token余额
//  3. 计算应返还的代币数量（根据LP Token份额）
//  4. 销毁LP Token
//  5. 转移代币给用户
//  6. 发出移除流动性事件
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中需要实现：
//   - 应返还代币数量计算（amount = (lpTokenAmount / totalLPTokens) * totalReserve）
//   - LP Token销毁
//
// 返回：
//   - framework.SUCCESS - 移除成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - RemoveLiquidity - 移除流动性事件
//     {
//       "provider": "<流动性提供者地址>",
//       "token_id": "TOKEN_001",
//       "amount": 10000,
//       "lp_token_amount": 100
//     }
//
//export RemoveLiquidity
func RemoveLiquidity() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	tokenIDStr := params.ParseJSON("token_id")
	lpTokenAmount := params.ParseJSONInt("lp_token_amount")

	if lpTokenAmount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析代币ID（可选）
	var tokenID framework.TokenID
	if tokenIDStr != "" {
		tokenID = framework.TokenID(tokenIDStr)
	}

	// 步骤3：获取调用者
	caller := framework.GetCaller()

	// 步骤4：检查LP Token余额
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该检查用户的LP Token余额
	//   这里简化处理，假设余额足够

	// 步骤5：计算应返还的代币数量
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该计算：
	//   amount = (lpTokenAmount / totalLPTokens) * totalReserve
	amount := lpTokenAmount * 100

	// 步骤6：销毁LP Token
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该销毁LP Token
	//   这里简化处理，不实际销毁

	// 步骤7：检查合约余额
	contractAddr := framework.GetContractAddress()
	contractBalance := framework.QueryUTXOBalance(contractAddr, tokenID)
	if contractBalance < framework.Amount(amount) {
		return framework.ERROR_INSUFFICIENT_BALANCE
	}

	// 步骤8：转移代币给用户
	err := token.Transfer(
		contractAddr,
		caller,
		tokenID,
		framework.Amount(amount),
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤9：发出移除流动性事件
	event := framework.NewEvent("RemoveLiquidity")
	event.AddAddressField("provider", caller)
	if tokenIDStr != "" {
		event.AddStringField("token_id", tokenIDStr)
	}
	event.AddUint64Field("amount", uint64(amount))
	event.AddUint64Field("lp_token_amount", uint64(lpTokenAmount))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// QueryPoolInfo 查询池信息
//
// 查询流动性池的详细信息，包括池中代币余额、LP Token总量等。
//
// 参数格式（JSON）:
//
//	{
//	  "token_id": "TOKEN_001"  // 代币ID（可选，nil表示原生代币）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 查询池中代币余额
//  3. 查询LP Token总量
//  4. 返回池信息
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中，应该从状态输出查询池信息
//   包括池中代币余额、LP Token总量、总流动性等
//
// 返回：
//   - framework.SUCCESS - 查询成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//
//export QueryPoolInfo
func QueryPoolInfo() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	tokenIDStr := params.ParseJSON("token_id")

	// 步骤2：解析代币ID（可选）
	var tokenID framework.TokenID
	if tokenIDStr != "" {
		tokenID = framework.TokenID(tokenIDStr)
	}

	// 步骤3：查询池中代币余额
	contractAddr := framework.GetContractAddress()
	totalReserve := framework.QueryUTXOBalance(contractAddr, tokenID)

	// 步骤4：查询LP Token总量
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从状态输出查询LP Token总量
	totalLPTokens := totalReserve / 100

	// 步骤5：返回池信息
	// 注意：实际应用中应该返回完整的池信息
	result := `{"token_id":"` + tokenIDStr + `","total_reserve":` + framework.Uint64ToString(uint64(totalReserve)) + `,"total_lp_tokens":` + framework.Uint64ToString(uint64(totalLPTokens)) + `}`
	framework.SetReturnData([]byte(result))

	return framework.SUCCESS
}

func main() {}

