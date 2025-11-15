//go:build tinygo || (js && wasm)

// Package main 提供AMM（自动化做市商）合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建 AMM（Automated Market Maker）合约。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/token 模块进行代币操作
//   - 如何使用 helpers/market 模块进行资产托管
//   - 如何实现完整的AMM功能（添加流动性、移除流动性、代币交换）
//
// 🎯 核心功能
//
//  1. AddLiquidity - 添加流动性
//     - 向流动性池添加代币对
//     - 获得流动性凭证代币（LP Token）
//
//  2. RemoveLiquidity - 移除流动性
//     - 从流动性池移除代币对
//     - 销毁流动性凭证代币
//
//  3. SwapTokens - 代币交换
//     - 使用恒定乘积公式（x*y=k）进行代币交换
//     - 自动计算交换价格和滑点
//
// ⚠️ 注意：本示例是简化实现
//   实际应用中需要实现：
//   - 恒定乘积公式（x*y=k）价格计算
//   - 滑点保护机制
//   - 手续费分成（给流动性提供者）
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

// AMMContract AMM（自动化做市商）合约
//
// 本合约使用 helpers/token 和 helpers/market 模块提供的业务语义API，
// 简化AMM操作的实现，开发者只需关注业务逻辑。
//
// AMM特点：
//   - 使用恒定乘积公式（x*y=k）进行价格发现
//   - 流动性提供者获得交易手续费分成
//   - 无需订单簿，自动匹配交易
type AMMContract struct {
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
//       "contract": "AMM",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "AMM")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// AddLiquidity 添加流动性
//
// 向流动性池添加代币对，获得流动性凭证代币（LP Token）。
// SDK 内部会自动处理：
//   - 余额检查（确保用户余额充足）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 AddLiquidity 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "token_a_id": "TOKEN_A",  // 代币A ID（必填）
//	  "token_b_id": "TOKEN_B",  // 代币B ID（必填）
//	  "amount_a": 1000,         // 代币A数量（必填）
//	  "amount_b": 2000          // 代币B数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 检查用户余额
//  3. 计算流动性凭证代币数量（根据恒定乘积公式）
//  4. 转移代币到合约
//  5. 铸造流动性凭证代币
//  6. 发出添加流动性事件
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中需要实现：
//   - 恒定乘积公式计算（x*y=k）
//   - 流动性凭证代币数量计算
//   - 首次添加流动性的特殊处理
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
//       "token_a_id": "TOKEN_A",
//       "token_b_id": "TOKEN_B",
//       "amount_a": 1000,
//       "amount_b": 2000
//     }
//
//export AddLiquidity
func AddLiquidity() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	tokenAIDStr := params.ParseJSON("token_a_id")
	tokenBIDStr := params.ParseJSON("token_b_id")
	amountA := params.ParseJSONInt("amount_a")
	amountB := params.ParseJSONInt("amount_b")

	if tokenAIDStr == "" || tokenBIDStr == "" || amountA == 0 || amountB == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析代币ID
	tokenAID := framework.TokenID(tokenAIDStr)
	tokenBID := framework.TokenID(tokenBIDStr)

	// 步骤3：获取调用者
	caller := framework.GetCaller()

	// 步骤4：检查余额
	balanceA := framework.QueryUTXOBalance(caller, tokenAID)
	balanceB := framework.QueryUTXOBalance(caller, tokenBID)
	if balanceA < framework.Amount(amountA) || balanceB < framework.Amount(amountB) {
		return framework.ERROR_INSUFFICIENT_BALANCE
	}

	// 步骤5：计算流动性凭证代币数量
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该使用恒定乘积公式计算：
	//   LP Token数量 = sqrt(amountA * amountB) - 最小流动性
	//   首次添加流动性时，LP Token数量 = sqrt(amountA * amountB)

	// 步骤6：转移代币到合约
	contractAddr := framework.GetContractAddress()
	err := token.Transfer(
		caller,
		contractAddr,
		tokenAID,
		framework.Amount(amountA),
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	err = token.Transfer(
		caller,
		contractAddr,
		tokenBID,
		framework.Amount(amountB),
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
	//   LP Token数量 = sqrt(amountA * amountB) - 最小流动性

	// 步骤8：发出添加流动性事件
	event := framework.NewEvent("AddLiquidity")
	event.AddAddressField("provider", caller)
	event.AddStringField("token_a_id", tokenAIDStr)
	event.AddStringField("token_b_id", tokenBIDStr)
	event.AddUint64Field("amount_a", uint64(amountA))
	event.AddUint64Field("amount_b", uint64(amountB))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// RemoveLiquidity 移除流动性
//
// 从流动性池移除代币对，销毁流动性凭证代币。
// SDK 内部会自动处理：
//   - 余额检查（确保用户有足够的LP Token）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 RemoveLiquidity 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "token_a_id": "TOKEN_A",  // 代币A ID（必填）
//	  "token_b_id": "TOKEN_B",  // 代币B ID（必填）
//	  "lp_token_amount": 100    // LP Token数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 检查LP Token余额
//  3. 计算应返还的代币数量（根据恒定乘积公式）
//  4. 销毁LP Token
//  5. 转移代币给用户
//  6. 发出移除流动性事件
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中需要实现：
//   - 恒定乘积公式计算（x*y=k）
//   - 应返还代币数量计算
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
//       "token_a_id": "TOKEN_A",
//       "token_b_id": "TOKEN_B",
//       "amount_a": 1000,
//       "amount_b": 2000
//     }
//
//export RemoveLiquidity
func RemoveLiquidity() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	tokenAIDStr := params.ParseJSON("token_a_id")
	tokenBIDStr := params.ParseJSON("token_b_id")
	lpTokenAmount := params.ParseJSONInt("lp_token_amount")

	if tokenAIDStr == "" || tokenBIDStr == "" || lpTokenAmount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析代币ID
	tokenAID := framework.TokenID(tokenAIDStr)
	tokenBID := framework.TokenID(tokenBIDStr)

	// 步骤3：获取调用者
	caller := framework.GetCaller()

	// 步骤4：检查LP Token余额
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该检查用户的LP Token余额
	//   这里简化处理，假设余额足够

	// 步骤5：计算应返还的代币数量
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该使用恒定乘积公式计算：
	//   amountA = (lpTokenAmount / totalLPTokens) * reserveA
	//   amountB = (lpTokenAmount / totalLPTokens) * reserveB

	// 步骤6：销毁LP Token
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该销毁LP Token
	//   这里简化处理，不实际销毁

	// 步骤7：转移代币给用户
	contractAddr := framework.GetContractAddress()
	// 简化处理：假设返还相同比例
	amountA := lpTokenAmount
	amountB := lpTokenAmount * 2

	err := token.Transfer(
		contractAddr,
		caller,
		tokenAID,
		framework.Amount(amountA),
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	err = token.Transfer(
		contractAddr,
		caller,
		tokenBID,
		framework.Amount(amountB),
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤8：发出移除流动性事件
	event := framework.NewEvent("RemoveLiquidity")
	event.AddAddressField("provider", caller)
	event.AddStringField("token_a_id", tokenAIDStr)
	event.AddStringField("token_b_id", tokenBIDStr)
	event.AddUint64Field("amount_a", uint64(amountA))
	event.AddUint64Field("amount_b", uint64(amountB))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// SwapTokens 代币交换
//
// 使用恒定乘积公式（x*y=k）进行代币交换。
// SDK 内部会自动处理：
//   - 余额检查（确保用户余额充足）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 SwapTokens 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "token_in_id": "TOKEN_A",   // 输入代币ID（必填）
//	  "token_out_id": "TOKEN_B",  // 输出代币ID（必填）
//	  "amount_in": 1000,          // 输入数量（必填）
//	  "min_amount_out": 1800     // 最小输出数量（必填，滑点保护）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 检查用户余额
//  3. 计算输出数量（使用恒定乘积公式）
//  4. 检查滑点（确保输出数量 >= min_amount_out）
//  5. 转移输入代币到合约
//  6. 转移输出代币给用户
//  7. 计算手续费（给流动性提供者）
//  8. 发出交换事件
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中需要实现：
//   - 恒定乘积公式计算（x*y=k）
//   - 滑点保护机制
//   - 手续费分成（给流动性提供者）
//
// 返回：
//   - framework.SUCCESS - 交换成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_SLIPPAGE_EXCEEDED - 滑点过大
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - SwapTokens - 代币交换事件
//     {
//       "trader": "<交易者地址>",
//       "token_in_id": "TOKEN_A",
//       "token_out_id": "TOKEN_B",
//       "amount_in": 1000,
//       "amount_out": 1900
//     }
//
//export SwapTokens
func SwapTokens() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	tokenInIDStr := params.ParseJSON("token_in_id")
	tokenOutIDStr := params.ParseJSON("token_out_id")
	amountIn := params.ParseJSONInt("amount_in")
	minAmountOut := params.ParseJSONInt("min_amount_out")

	if tokenInIDStr == "" || tokenOutIDStr == "" || amountIn == 0 || minAmountOut == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析代币ID
	tokenInID := framework.TokenID(tokenInIDStr)
	tokenOutID := framework.TokenID(tokenOutIDStr)

	// 步骤3：获取调用者
	caller := framework.GetCaller()

	// 步骤4：检查余额
	balance := framework.QueryUTXOBalance(caller, tokenInID)
	if balance < framework.Amount(amountIn) {
		return framework.ERROR_INSUFFICIENT_BALANCE
	}

	// 步骤5：计算输出数量（使用恒定乘积公式）
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该使用恒定乘积公式计算：
	//   amountOut = (reserveOut * amountIn) / (reserveIn + amountIn)
	//   这里简化处理，假设固定比例
	amountOut := amountIn * 2

	// 步骤6：检查滑点
	if amountOut < minAmountOut {
		return framework.ERROR_EXECUTION_FAILED // 滑点过大
	}

	// 步骤7：转移输入代币到合约
	contractAddr := framework.GetContractAddress()
	err := token.Transfer(
		caller,
		contractAddr,
		tokenInID,
		framework.Amount(amountIn),
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤8：计算手续费（给流动性提供者）
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该计算手续费（如0.3%）
	//   手续费 = amountOut * 0.003
	//   实际输出 = amountOut - 手续费
	actualAmountOut := amountOut

	// 步骤9：转移输出代币给用户
	err = token.Transfer(
		contractAddr,
		caller,
		tokenOutID,
		framework.Amount(actualAmountOut),
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤10：发出交换事件
	event := framework.NewEvent("SwapTokens")
	event.AddAddressField("trader", caller)
	event.AddStringField("token_in_id", tokenInIDStr)
	event.AddStringField("token_out_id", tokenOutIDStr)
	event.AddUint64Field("amount_in", uint64(amountIn))
	event.AddUint64Field("amount_out", uint64(actualAmountOut))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

func main() {}

