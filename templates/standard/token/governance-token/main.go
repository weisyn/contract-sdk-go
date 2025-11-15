//go:build tinygo || (js && wasm)

// Package main 提供治理代币合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建治理代币合约。
// 治理代币是一种特殊的代币，持有者可以使用代币进行投票，参与去中心化治理。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/token 模块创建治理代币
//   - 如何将代币持有量与投票权关联
//   - 如何实现治理代币的铸造、转移、投票等功能
//
// 🎯 核心功能
//
//  1. Mint - 铸造治理代币
//     - 使用 token.Mint() 铸造治理代币
//     - 持有代币即拥有投票权
//
//  2. Transfer - 转移治理代币
//     - 使用 token.Transfer() 转移代币
//     - 转移代币时，投票权也随之转移
//
//  3. DelegateVotingPower - 委托投票权
//     - 将投票权委托给其他地址
//     - 委托后，被委托者可以使用委托的代币进行投票
//
//  4. VoteWithTokens - 使用代币投票
//     - 使用治理代币进行投票
//     - 投票权重 = 持有的代币数量
//
// 📚 相关文档
//
//   - [Token 模块文档](../../helpers/token/README.md)
//   - [Governance 模块文档](../../helpers/governance/README.md)
//   - [Framework 文档](../../framework/README.md)
//   - [示例总览](../README.md)
package main

import (
	"github.com/weisyn/contract-sdk-go/helpers/governance"
	"github.com/weisyn/contract-sdk-go/helpers/token"
	"github.com/weisyn/contract-sdk-go/framework"
)

// GovernanceTokenContract 治理代币合约
//
// 本合约使用 helpers/token 和 helpers/governance 模块提供的业务语义API，
// 简化治理代币操作的实现，开发者只需关注业务逻辑。
//
// 治理代币特点：
//   - 持有代币即拥有投票权
//   - 投票权重 = 持有的代币数量
//   - 支持投票权委托
type GovernanceTokenContract struct {
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
//       "contract": "GovernanceToken",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "GovernanceToken")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Mint 铸造治理代币
//
// 使用 helpers/token 模块的 Mint 函数铸造治理代币。
// 持有代币即拥有投票权，投票权重等于持有的代币数量。
// SDK 内部会自动处理：
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Mint 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "to": "receiver_address",    // 接收者地址（Base58编码，必填）
//	  "amount": 1000               // 铸造数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析接收者地址
//  3. 调用 token.Mint() 进行铸造
//     - SDK 内部自动构建交易
//  4. 返回执行结果
//
// ⚠️ 注意：实际应用中需要权限检查
//   - 只有授权地址才能调用 Mint
//   - 权限检查逻辑应在应用层实现
//
// 返回：
//   - framework.SUCCESS - 铸造成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Mint - 铸造事件（由 SDK 自动发出）
//     {
//       "to": "<接收者地址>",
//       "amount": 1000
//     }
//
//export Mint
func Mint() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	toStr := params.ParseJSON("to")
	amount := params.ParseJSONInt("amount")

	if toStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析接收者地址
	to, err := framework.ParseAddressBase58(toStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：使用 SDK 基础能力进行代币铸造
	//
	// SDK 提供的 token.Mint() 会自动处理：
	//   - 交易构建
	//   - 事件发出
	//
	// ⚠️ 注意：实际应用中需要权限检查
	//   只有授权地址才能调用 Mint，权限检查逻辑应在应用层实现
	err = token.Mint(to, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Transfer 转移治理代币
//
// 使用 helpers/token 模块的 Transfer 函数转移治理代币。
// 转移代币时，投票权也随之转移。
// SDK 内部会自动处理：
//   - 余额检查（确保发送者余额充足）
//   - 交易构建（自动构建 UTXO 交易）
//   - 找零处理（自动处理找零 UTXO）
//   - 事件发出（自动发出 Transfer 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "to": "receiver_address",    // 接收者地址（Base58编码，必填）
//	  "amount": 100                // 转账数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析接收者地址
//  3. 调用 token.Transfer() 进行转账
//     - SDK 内部自动处理余额检查
//     - SDK 内部自动构建交易
//  4. 返回执行结果
//
// 返回：
//   - framework.SUCCESS - 转账成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Transfer - 转账事件（由 SDK 自动发出）
//     {
//       "from": "<发送者地址>",
//       "to": "<接收者地址>",
//       "amount": 100
//     }
//
//export Transfer
func Transfer() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	toStr := params.ParseJSON("to")
	amount := params.ParseJSONInt("amount")

	if toStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析接收者地址
	to, err := framework.ParseAddressBase58(toStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：使用 SDK 基础能力进行代币转账
	//
	// SDK 提供的 token.Transfer() 会自动处理：
	//   - 余额检查
	//   - 交易构建
	//   - 事件发出
	caller := framework.GetCaller()
	err = token.Transfer(caller, to, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// DelegateVotingPower 委托投票权
//
// 将投票权委托给其他地址。
// 委托后，被委托者可以使用委托的代币进行投票。
//
// 参数格式（JSON）:
//
//	{
//	  "delegate": "delegate_address",  // 被委托者地址（Base58编码，必填）
//	  "amount": 500                   // 委托的代币数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析被委托者地址
//  3. 检查委托者余额
//  4. 记录委托关系（使用状态输出）
//  5. 发出委托事件
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中，应该使用状态输出存储委托关系
//   并在投票时检查委托的代币数量
//
// 返回：
//   - framework.SUCCESS - 委托成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - VotingPowerDelegated - 投票权委托事件
//     {
//       "delegator": "<委托者地址>",
//       "delegate": "<被委托者地址>",
//       "amount": 500
//     }
//
//export DelegateVotingPower
func DelegateVotingPower() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	delegateStr := params.ParseJSON("delegate")
	amount := params.ParseJSONInt("amount")

	if delegateStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析被委托者地址
	delegate, err := framework.ParseAddressBase58(delegateStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：检查委托者余额
	caller := framework.GetCaller()
	balance := framework.QueryUTXOBalance(caller, framework.TokenID(""))
	if balance < framework.Amount(amount) {
		return framework.ERROR_INSUFFICIENT_BALANCE
	}

	// 步骤4：记录委托关系
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该使用状态输出存储委托关系
	//   并在投票时检查委托的代币数量
	//   这里只发出事件，实际委托关系应该在应用层维护

	// 步骤5：发出委托事件
	event := framework.NewEvent("VotingPowerDelegated")
	event.AddAddressField("delegator", caller)
	event.AddAddressField("delegate", delegate)
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// VoteWithTokens 使用代币投票
//
// 使用治理代币进行投票。
// 投票权重 = 持有的代币数量 + 委托的代币数量
//
// 参数格式（JSON）:
//
//	{
//	  "proposal_id": "proposal_123",  // 提案ID（必填）
//	  "support": true                 // 是否支持（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 计算投票权重（持有的代币数量）
//  3. 调用 governance.Vote() 进行投票
//     - SDK 内部自动构建状态输出
//  4. 发出投票事件（包含投票权重）
//  5. 返回执行结果
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中，应该考虑委托的代币数量
//   投票权重 = 持有的代币数量 + 委托的代币数量
//
// 返回：
//   - framework.SUCCESS - 投票成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Vote - 投票事件（由 SDK 自动发出）
//     {
//       "voter": "<投票者地址>",
//       "proposal_id": "proposal_123",
//       "support": true
//     }
//   - TokenVote - 代币投票事件（自定义）
//     {
//       "voter": "<投票者地址>",
//       "proposal_id": "proposal_123",
//       "support": true,
//       "voting_power": 1000
//     }
//
//export VoteWithTokens
func VoteWithTokens() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	proposalIDStr := params.ParseJSON("proposal_id")
	supportStr := params.ParseJSON("support")

	if proposalIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析支持/反对
	support := supportStr == "true" || supportStr == "1"

	// 步骤3：计算投票权重（持有的代币数量）
	caller := framework.GetCaller()
	votingPower := framework.QueryUTXOBalance(caller, framework.TokenID(""))

	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该考虑委托的代币数量
	//   投票权重 = 持有的代币数量 + 委托的代币数量

	// 步骤4：使用 SDK 基础能力进行投票
	//
	// SDK 提供的 governance.Vote() 会自动处理：
	//   - 状态输出构建
	//   - 事件发出
	err := governance.Vote(
		caller,
		[]byte(proposalIDStr),
		support,
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤5：发出代币投票事件（包含投票权重）
	event := framework.NewEvent("TokenVote")
	event.AddAddressField("voter", caller)
	event.AddStringField("proposal_id", proposalIDStr)
	event.AddField("support", support) // 使用 AddField 支持 bool 类型
	event.AddUint64Field("voting_power", uint64(votingPower))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

func main() {}

