//go:build tinygo || (js && wasm)

// Package main 提供DAO治理合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建DAO（去中心化自治组织）治理合约。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/governance 模块进行DAO治理
//   - 如何使用 helpers/token 模块进行治理代币管理
//   - 如何实现完整的DAO功能（提案创建、投票、执行等）
//
// 🎯 核心功能
//
//  1. CreateProposal - 创建提案
//     - 创建新的治理提案
//     - 设置提案内容和投票参数
//
//  2. Vote - 投票
//     - 使用治理代币进行投票
//     - 投票权重 = 持有的代币数量
//
//  3. ExecuteProposal - 执行提案
//     - 执行已通过的提案
//     - 自动检查提案状态和投票结果
//
//  4. QueryProposal - 查询提案
//     - 查询提案的详细信息
//     - 查询投票结果
//
// 📚 相关文档
//
//   - [Governance 模块文档](../../helpers/governance/README.md)
//   - [Token 模块文档](../../helpers/token/README.md)
//   - [Framework 文档](../../framework/README.md)
//   - [示例总览](../README.md)
package main

import (
	"github.com/weisyn/contract-sdk-go/helpers/governance"
	"github.com/weisyn/contract-sdk-go/framework"
)

// DAOContract DAO治理合约
//
// 本合约使用 helpers/governance 和 helpers/token 模块提供的业务语义API，
// 简化DAO治理操作的实现，开发者只需关注业务逻辑。
//
// DAO特点：
//   - 去中心化决策
//   - 代币持有者投票
//   - 提案执行自动化
type DAOContract struct {
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
//       "contract": "DAO",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "DAO")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// CreateProposal 创建提案
//
// 创建新的治理提案。
// 提案包含提案内容、投票期限、通过阈值等信息。
//
// 参数格式（JSON）:
//
//	{
//	  "proposal_id": "proposal_001",     // 提案ID（必填）
//	  "title": "Proposal Title",        // 提案标题（必填）
//	  "description": "Proposal desc",   // 提案描述（可选）
//	  "voting_period": 604800,          // 投票期限（秒，可选）
//	  "threshold": 50                   // 通过阈值（百分比，可选）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 检查提案ID唯一性
//  3. 创建提案状态（使用状态输出）
//  4. 发出提案创建事件
//  5. 返回执行结果
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中，应该使用状态输出存储提案信息
//   包括提案内容、投票期限、通过阈值、投票结果等
//
// 返回：
//   - framework.SUCCESS - 创建成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_ALREADY_EXISTS - 提案已存在
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - ProposalCreated - 提案创建事件
//     {
//       "creator": "<创建者地址>",
//       "proposal_id": "proposal_001",
//       "title": "Proposal Title"
//     }
//
//export CreateProposal
func CreateProposal() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	proposalIDStr := params.ParseJSON("proposal_id")
	titleStr := params.ParseJSON("title")

	if proposalIDStr == "" || titleStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：检查提案ID唯一性
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从状态输出查询提案是否存在

	// 步骤3：创建提案状态
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该使用状态输出存储提案信息
	//   包括提案内容、投票期限、通过阈值、投票结果等

	// 步骤4：发出提案创建事件
	caller := framework.GetCaller()
	descriptionStr := params.ParseJSON("description")
	votingPeriodStr := params.ParseJSON("voting_period")
	thresholdStr := params.ParseJSON("threshold")

	event := framework.NewEvent("ProposalCreated")
	event.AddAddressField("creator", caller)
	event.AddStringField("proposal_id", proposalIDStr)
	event.AddStringField("title", titleStr)
	if descriptionStr != "" {
		event.AddStringField("description", descriptionStr)
	}
	if votingPeriodStr != "" {
		event.AddStringField("voting_period", votingPeriodStr)
	}
	if thresholdStr != "" {
		event.AddStringField("threshold", thresholdStr)
	}
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Vote 投票
//
// 使用治理代币进行投票。
// 投票权重 = 持有的代币数量。
// SDK 内部会自动处理：
//   - 状态输出构建（自动构建投票状态）
//   - 事件发出（自动发出 Vote 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "proposal_id": "proposal_001",  // 提案ID（必填）
//	  "support": true                 // 是否支持（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 检查提案是否存在
//  3. 检查投票期限
//  4. 计算投票权重（持有的代币数量）
//  5. 调用 governance.Vote() 进行投票
//     - SDK 内部自动构建状态输出
//  6. 发出投票事件（自定义事件，包含投票权重）
//  7. 返回执行结果
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中，应该：
//   - 检查提案是否存在
//   - 检查投票期限
//   - 检查是否已投票（防止重复投票）
//   - 计算投票权重（持有的代币数量 + 委托的代币数量）
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
//       "proposal_id": "proposal_001",
//       "support": true
//     }
//   - DAOVote - DAO投票事件（自定义）
//     {
//       "voter": "<投票者地址>",
//       "proposal_id": "proposal_001",
//       "support": true,
//       "voting_power": 1000
//     }
//
//export Vote
func Vote() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	proposalIDStr := params.ParseJSON("proposal_id")
	supportStr := params.ParseJSON("support")

	if proposalIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析支持/反对
	support := supportStr == "true" || supportStr == "1"

	// 步骤3：检查提案是否存在
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从状态输出查询提案是否存在

	// 步骤4：检查投票期限
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该检查投票期限是否已过

	// 步骤5：计算投票权重（持有的代币数量）
	caller := framework.GetCaller()
	votingPower := framework.QueryUTXOBalance(caller, framework.TokenID(""))

	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该考虑委托的代币数量
	//   投票权重 = 持有的代币数量 + 委托的代币数量

	// 步骤6：使用 SDK 基础能力进行投票
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

	// 步骤7：发出DAO投票事件（自定义事件，包含投票权重）
	event := framework.NewEvent("DAOVote")
	event.AddAddressField("voter", caller)
	event.AddStringField("proposal_id", proposalIDStr)
	event.AddField("support", support)
	event.AddUint64Field("voting_power", uint64(votingPower))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// ExecuteProposal 执行提案
//
// 执行已通过的提案。
// 实际应用中，应该检查提案状态和投票结果，然后执行提案内容。
//
// 参数格式（JSON）:
//
//	{
//	  "proposal_id": "proposal_001"  // 提案ID（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 查询提案状态和投票结果
//  3. 检查提案是否已通过
//  4. 检查提案是否已执行
//  5. 执行提案内容
//  6. 更新提案状态
//  7. 发出提案执行事件
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中，应该：
//   - 从状态输出查询提案状态和投票结果
//   - 检查提案是否已通过（支持票数 >= 阈值）
//   - 检查提案是否已执行（防止重复执行）
//   - 执行提案内容（调用其他合约、转移资金等）
//   - 更新提案状态（使用状态输出）
//
// 返回：
//   - framework.SUCCESS - 执行成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_NOT_FOUND - 提案不存在
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - ProposalExecuted - 提案执行事件
//     {
//       "executor": "<执行者地址>",
//       "proposal_id": "proposal_001"
//     }
//
//export ExecuteProposal
func ExecuteProposal() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	proposalIDStr := params.ParseJSON("proposal_id")

	if proposalIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：查询提案状态和投票结果
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从状态输出查询提案状态和投票结果

	// 步骤3：检查提案是否已通过
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该检查支持票数是否 >= 阈值

	// 步骤4：检查提案是否已执行
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该检查提案是否已执行（防止重复执行）

	// 步骤5：执行提案内容
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该执行提案内容（调用其他合约、转移资金等）

	// 步骤6：更新提案状态
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该使用状态输出更新提案状态

	// 步骤7：发出提案执行事件
	caller := framework.GetCaller()
	event := framework.NewEvent("ProposalExecuted")
	event.AddAddressField("executor", caller)
	event.AddStringField("proposal_id", proposalIDStr)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// QueryProposal 查询提案
//
// 查询提案的详细信息，包括提案内容、投票结果、执行状态等。
//
// 参数格式（JSON）:
//
//	{
//	  "proposal_id": "proposal_001"  // 提案ID（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 查询提案信息（从状态输出）
//  3. 查询投票结果（从状态输出）
//  4. 返回查询结果
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中，应该从状态输出查询提案信息
//   包括提案内容、投票期限、通过阈值、投票结果、执行状态等
//
// 返回：
//   - framework.SUCCESS - 查询成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_NOT_FOUND - 提案不存在
//
//export QueryProposal
func QueryProposal() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	proposalIDStr := params.ParseJSON("proposal_id")

	if proposalIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：查询提案信息
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从状态输出查询提案信息
	//   包括提案内容、投票期限、通过阈值、投票结果、执行状态等

	// 步骤3：返回查询结果
	// 注意：实际应用中应该返回完整的提案信息
	result := `{"proposal_id":"` + proposalIDStr + `","status":"active"}`
	framework.SetReturnData([]byte(result))

	return framework.SUCCESS
}

func main() {}

