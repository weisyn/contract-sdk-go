//go:build tinygo || (js && wasm)

// Package main 提供提案投票治理合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建去中心化治理相关的智能合约。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/governance 模块进行治理操作
//   - 如何使用业务语义API简化治理合约开发
//   - 如何实现完整的治理功能（Propose、Vote）
//
// 🎯 核心功能
//
//  1. Propose - 创建提案
//     - 使用 governance.Propose() 创建治理提案
//     - SDK 内部自动处理状态输出、事件发出
//
//  2. Vote - 投票
//     - 使用 governance.Vote() 对提案进行投票
//     - 支持支持/反对两种投票方式
//
// 📚 相关文档
//
//   - [Governance 模块文档](../../helpers/governance/README.md)
//   - [Framework 文档](../../framework/README.md)
//   - [示例总览](../README.md)
package main

import (
	"github.com/weisyn/contract-sdk-go/helpers/governance"
	"github.com/weisyn/contract-sdk-go/framework"
)

// GovernanceContract 提案投票治理合约
//
// 本合约使用 helpers/governance 模块提供的业务语义API，
// 简化治理操作的实现，开发者只需关注业务逻辑。
type GovernanceContract struct {
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
//       "contract": "Governance",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "Governance")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Propose 创建提案
//
// 使用 helpers/governance 模块的 Propose 函数创建治理提案。
// SDK 内部会自动处理：
//   - 状态输出构建（自动构建提案状态输出）
//   - 事件发出（自动发出 Propose 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "proposal_id": "proposal_123",        // 提案ID（必填）
//	  "proposal_data": "proposal content"  // 提案内容（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 调用 governance.Propose() 创建提案
//     - SDK 内部自动构建状态输出
//     - SDK 内部自动发出事件
//  3. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - 提案创建权限检查（谁可以创建提案）
//   - 提案格式验证（提案内容是否符合规范）
//   - 提案ID唯一性检查
//
// 返回：
//   - framework.SUCCESS - 提案创建成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Propose - 提案创建事件（由 SDK 自动发出）
//     {
//       "proposer": "<提案创建者地址>",
//       "proposal_id": "proposal_123",
//       "proposal_data": "proposal content"
//     }
//
//export Propose
func Propose() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	proposalIDStr := params.ParseJSON("proposal_id")
	proposalDataStr := params.ParseJSON("proposal_data")

	if proposalIDStr == "" || proposalDataStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：使用 SDK 基础能力创建提案
	//
	// SDK 提供的 governance.Propose() 会自动处理：
	//   - 状态输出构建
	//   - 事件发出
	//
	// ⚠️ 注意：实际应用中需要业务规则检查
	//   提案创建权限、提案格式验证、提案ID唯一性等应在应用层实现
	caller := framework.GetCaller()
	err := governance.Propose(
		caller,
		[]byte(proposalIDStr),
		[]byte(proposalDataStr),
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Vote 投票
//
// 使用 helpers/governance 模块的 Vote 函数对提案进行投票。
// SDK 内部会自动处理：
//   - 状态输出构建（自动构建投票状态输出）
//   - 事件发出（自动发出 Vote 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "proposal_id": "proposal_123",  // 提案ID（必填）
//	  "support": true                 // 是否支持（必填，true表示支持，false表示反对）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析支持/反对（支持 "true"/"1" 表示支持，"false"/"0" 表示反对）
//  3. 调用 governance.Vote() 进行投票
//     - SDK 内部自动构建状态输出
//     - SDK 内部自动发出事件
//  4. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - 提案存在性检查（提案是否已创建）
//   - 投票权限检查（谁可以投票，是否已投票）
//   - 投票时间窗口检查（是否在投票期内）
//   - 投票权重计算（基于代币持有量或其他规则）
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
//
//export Vote
func Vote() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	proposalIDStr := params.ParseJSON("proposal_id")
	supportStr := params.ParseJSON("support")

	if proposalIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析支持/反对
	// 支持 "true"/"1" 表示支持，"false"/"0" 表示反对
	support := supportStr == "true" || supportStr == "1"

	// 步骤3：使用 SDK 基础能力进行投票
	//
	// SDK 提供的 governance.Vote() 会自动处理：
	//   - 状态输出构建
	//   - 事件发出
	//
	// ⚠️ 注意：实际应用中需要业务规则检查
	//   提案存在性、投票权限、投票时间窗口、投票权重等应在应用层实现
	caller := framework.GetCaller()
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

	return framework.SUCCESS
}

func main() {}

