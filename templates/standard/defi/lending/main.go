//go:build tinygo || (js && wasm)

// Package main 提供借贷协议合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建 DeFi 借贷协议合约。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/token 模块进行代币操作
//   - 如何使用 helpers/market 模块进行资产托管
//   - 如何实现完整的借贷功能（存款、借款、还款、取款）
//
// 🎯 核心功能
//
//  1. Deposit - 存款
//     - 用户存入代币作为抵押品
//     - 获得存款凭证代币
//
//  2. Borrow - 借款
//     - 使用抵押品借出代币
//     - 需要满足抵押率要求
//
//  3. Repay - 还款
//     - 偿还借款本金和利息
//     - 释放抵押品
//
//  4. Withdraw - 取款
//     - 取出存款和收益
//
// ⚠️ 注意：本示例是简化实现
//   实际应用中需要实现：
//   - 利率计算（根据市场供需动态调整）
//   - 抵押率检查（确保抵押品价值足够）
//   - 清算机制（抵押品不足时清算）
//   - 存款凭证代币管理
//
// 📚 相关文档
//
//   - [Token 模块文档](../../../helpers/token/README.md)
//   - [Market 模块文档](../../../helpers/market/README.md)
//   - [Framework 文档](../../../framework/README.md)
//   - [示例总览](../../README.md)
package main

import (
	"github.com/weisyn/contract-sdk-go/helpers/market"
	"github.com/weisyn/contract-sdk-go/helpers/token"
	"github.com/weisyn/contract-sdk-go/framework"
)

// LendingContract 借贷协议合约
//
// 本合约使用 helpers/token 和 helpers/market 模块提供的业务语义API，
// 简化借贷操作的实现，开发者只需关注业务逻辑。
//
// 借贷协议特点：
//   - 用户存入代币作为抵押品
//   - 根据抵押品价值借出代币
//   - 需要支付利息
//   - 抵押率不足时会被清算
type LendingContract struct {
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
//       "contract": "Lending",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "Lending")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Deposit 存款
//
// 用户存入代币作为抵押品，获得存款凭证代币。
// SDK 内部会自动处理：
//   - 余额检查（确保用户余额充足）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Deposit 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "token_id": "TOKEN_001",  // 代币ID（可选，nil表示原生代币）
//	  "amount": 10000            // 存款数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 检查用户余额
//  3. 转移代币到合约（使用托管）
//  4. 铸造存款凭证代币
//  5. 发出存款事件
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中需要实现：
//   - 存款凭证代币的铸造和管理
//   - 存款利率计算
//   - 存款总额统计
//
// 返回：
//   - framework.SUCCESS - 存款成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Deposit - 存款事件
//     {
//       "depositor": "<存款者地址>",
//       "token_id": "TOKEN_001",
//       "amount": 10000
//     }
//
//export Deposit
func Deposit() uint32 {
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

	// 步骤5：转移代币到合约（使用托管）
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该将代币转移到合约地址，并记录存款信息
	//   这里使用托管作为示例
	contractAddr := framework.GetContractAddress()
	err := market.Escrow(
		caller,                        // 存款者
		contractAddr,                  // 合约地址（作为托管方）
		tokenID,                       // 代币ID
		framework.Amount(amount),      // 存款数量
		[]byte("deposit"),            // 托管ID
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤6：铸造存款凭证代币
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该铸造存款凭证代币（cToken）给用户
	//   凭证代币数量 = 存款数量 * 凭证汇率
	//   这里简化处理，不实际铸造凭证代币

	// 步骤7：发出存款事件
	event := framework.NewEvent("Deposit")
	event.AddAddressField("depositor", caller)
	if tokenIDStr != "" {
		event.AddStringField("token_id", tokenIDStr)
	}
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Borrow 借款
//
// 使用抵押品借出代币。
// SDK 内部会自动处理：
//   - 余额检查（确保合约有足够的代币）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Borrow 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "token_id": "TOKEN_002",  // 借款代币ID（可选，nil表示原生代币）
//	  "amount": 5000            // 借款数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 检查抵押品价值
//  3. 检查抵押率（抵押品价值 / 借款价值 >= 抵押率）
//  4. 转移代币给借款人
//  5. 记录借款信息
//  6. 发出借款事件
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中需要实现：
//   - 抵押品价值查询（需要价格预言机或ISPC受控机制）
//   - 抵押率检查（确保抵押品价值足够）
//   - 借款利率计算
//   - 借款总额统计
//
// 返回：
//   - framework.SUCCESS - 借款成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Borrow - 借款事件
//     {
//       "borrower": "<借款人地址>",
//       "token_id": "TOKEN_002",
//       "amount": 5000
//     }
//
//export Borrow
func Borrow() uint32 {
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

	// 步骤4：检查抵押品价值
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该查询用户的抵押品价值
	//   抵押品价值 = 存款数量 * 代币价格
	//   这里简化处理，假设抵押品价值足够

	// 步骤5：检查抵押率
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该检查：
	//   抵押品价值 / 借款价值 >= 抵押率（如150%）
	//   这里简化处理，假设抵押率满足要求

	// 步骤6：检查合约余额
	contractAddr := framework.GetContractAddress()
	contractBalance := framework.QueryUTXOBalance(contractAddr, tokenID)
	if contractBalance < framework.Amount(amount) {
		return framework.ERROR_INSUFFICIENT_BALANCE
	}

	// 步骤7：转移代币给借款人
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从合约地址转移代币给借款人
	//   这里使用 token.Transfer 作为示例
	err := token.Transfer(
		contractAddr,                 // 从合约地址
		caller,                       // 到借款人地址
		tokenID,                      // 代币ID
		framework.Amount(amount),     // 借款数量
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤8：记录借款信息
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该使用状态输出记录借款信息
	//   包括借款数量、利率、到期时间等

	// 步骤9：发出借款事件
	event := framework.NewEvent("Borrow")
	event.AddAddressField("borrower", caller)
	if tokenIDStr != "" {
		event.AddStringField("token_id", tokenIDStr)
	}
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Repay 还款
//
// 偿还借款本金和利息，释放抵押品。
// SDK 内部会自动处理：
//   - 余额检查（确保用户余额充足）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Repay 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "token_id": "TOKEN_002",  // 借款代币ID（可选，nil表示原生代币）
//	  "amount": 5500            // 还款数量（本金+利息，必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 查询借款信息
//  3. 计算应还金额（本金+利息）
//  4. 转移代币到合约
//  5. 更新借款信息
//  6. 释放抵押品（部分或全部）
//  7. 发出还款事件
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中需要实现：
//   - 借款信息查询（从状态输出）
//   - 利息计算（根据借款时间和利率）
//   - 抵押品释放逻辑
//
// 返回：
//   - framework.SUCCESS - 还款成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Repay - 还款事件
//     {
//       "borrower": "<借款人地址>",
//       "token_id": "TOKEN_002",
//       "amount": 5500
//     }
//
//export Repay
func Repay() uint32 {
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

	// 步骤5：查询借款信息
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从状态输出查询借款信息
	//   包括借款数量、利率、到期时间等

	// 步骤6：计算应还金额
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该计算：
	//   应还金额 = 本金 + 利息
	//   利息 = 本金 * 利率 * 时间

	// 步骤7：转移代币到合约
	contractAddr := framework.GetContractAddress()
	err := token.Transfer(
		caller,                       // 从借款人地址
		contractAddr,                 // 到合约地址
		tokenID,                      // 代币ID
		framework.Amount(amount),     // 还款数量
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤8：更新借款信息
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该使用状态输出更新借款信息
	//   减少借款余额，更新还款时间等

	// 步骤9：释放抵押品
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该根据还款比例释放抵押品
	//   释放数量 = 还款数量 / 抵押率

	// 步骤10：发出还款事件
	event := framework.NewEvent("Repay")
	event.AddAddressField("borrower", caller)
	if tokenIDStr != "" {
		event.AddStringField("token_id", tokenIDStr)
	}
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Withdraw 取款
//
// 取出存款和收益。
// SDK 内部会自动处理：
//   - 余额检查（确保合约有足够的代币）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Withdraw 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "token_id": "TOKEN_001",  // 代币ID（可选，nil表示原生代币）
//	  "amount": 10500           // 取款数量（本金+收益，必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 查询存款信息
//  3. 计算可取金额（本金+收益）
//  4. 销毁存款凭证代币
//  5. 转移代币给用户
//  6. 发出取款事件
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中需要实现：
//   - 存款信息查询（从状态输出）
//   - 收益计算（根据存款时间和利率）
//   - 存款凭证代币的销毁
//
// 返回：
//   - framework.SUCCESS - 取款成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Withdraw - 取款事件
//     {
//       "depositor": "<存款者地址>",
//       "token_id": "TOKEN_001",
//       "amount": 10500
//     }
//
//export Withdraw
func Withdraw() uint32 {
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

	// 步骤4：查询存款信息
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从状态输出查询存款信息
	//   包括存款数量、存款凭证代币数量等

	// 步骤5：计算可取金额
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该计算：
	//   可取金额 = 存款数量 + 收益
	//   收益 = 存款数量 * 利率 * 时间

	// 步骤6：检查合约余额
	contractAddr := framework.GetContractAddress()
	contractBalance := framework.QueryUTXOBalance(contractAddr, tokenID)
	if contractBalance < framework.Amount(amount) {
		return framework.ERROR_INSUFFICIENT_BALANCE
	}

	// 步骤7：销毁存款凭证代币
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该销毁存款凭证代币（cToken）
	//   销毁数量 = 取款数量 / 凭证汇率

	// 步骤8：转移代币给用户
	err := token.Transfer(
		contractAddr,                 // 从合约地址
		caller,                       // 到存款者地址
		tokenID,                      // 代币ID
		framework.Amount(amount),     // 取款数量
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤9：发出取款事件
	event := framework.NewEvent("Withdraw")
	event.AddAddressField("depositor", caller)
	if tokenIDStr != "" {
		event.AddStringField("token_id", tokenIDStr)
	}
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

func main() {}

