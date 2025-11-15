//go:build tinygo || (js && wasm)

// Package main 提供委托质押合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建委托质押合约。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/staking 模块进行委托质押操作
//   - 如何实现质押权的委托和取消委托
//   - 如何管理委托关系和收益分配
//
// 🎯 核心功能
//
//  1. Delegate - 委托质押
//     - 使用 staking.Delegate() 将质押权委托给验证者
//     - 委托者仍持有代币，但质押权由被委托者行使
//
//  2. Undelegate - 取消委托
//     - 使用 staking.Undelegate() 取消委托
//     - 支持部分取消委托或全部取消委托
//
//  3. QueryDelegation - 查询委托信息
//     - 查询委托关系和委托数量
//     - 查询委托收益
//
// 📚 相关文档
//
//   - [Staking 模块文档](../../helpers/staking/README.md)
//   - [Framework 文档](../../framework/README.md)
//   - [示例总览](../README.md)
package main

import (
	"github.com/weisyn/contract-sdk-go/helpers/staking"
	"github.com/weisyn/contract-sdk-go/framework"
)

// DelegationContract 委托质押合约
//
// 本合约使用 helpers/staking 模块提供的业务语义API，
// 简化委托质押操作的实现，开发者只需关注业务逻辑。
//
// 委托质押特点：
//   - 委托者持有代币，但质押权由被委托者行使
//   - 委托者可以获得质押收益
//   - 支持部分委托和全部委托
type DelegationContract struct {
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
//       "contract": "Delegation",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "Delegation")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Delegate 委托质押
//
// 使用 helpers/staking 模块的 Delegate 函数将质押权委托给验证者。
// 委托者仍持有代币，但质押权由被委托者行使。
// SDK 内部会自动处理：
//   - 余额检查（确保委托者余额充足）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Delegate 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "validator": "validator_address",  // 验证者地址（Base58编码，必填）
//	  "delegatee": "delegatee_address", // 被委托者地址（Base58编码，必填）
//	  "amount": 10000                    // 委托数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析验证者和被委托者地址
//  3. 调用 staking.Delegate() 进行委托
//     - SDK 内部自动处理余额检查
//     - SDK 内部自动构建交易
//  4. 发出委托事件（自定义事件，包含更多信息）
//  5. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - 验证者有效性检查（验证者是否在验证者列表中）
//   - 最小委托数量检查
//   - 委托关系管理（使用状态输出存储）
//
// 返回：
//   - framework.SUCCESS - 委托成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Delegate - 委托事件（由 SDK 自动发出）
//     {
//       "delegator": "<委托者地址>",
//       "validator": "<验证者地址>",
//       "amount": 10000
//     }
//   - DelegationCreated - 委托创建事件（自定义）
//     {
//       "delegator": "<委托者地址>",
//       "validator": "<验证者地址>",
//       "delegatee": "<被委托者地址>",
//       "amount": 10000
//     }
//
//export Delegate
func Delegate() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	validatorStr := params.ParseJSON("validator")
	delegateeStr := params.ParseJSON("delegatee")
	amount := params.ParseJSONInt("amount")

	if validatorStr == "" || delegateeStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析验证者和被委托者地址
	validator, err1 := framework.ParseAddressBase58(validatorStr)
	delegatee, err2 := framework.ParseAddressBase58(delegateeStr)
	if err1 != nil || err2 != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：使用 SDK 基础能力进行委托
	//
	// SDK 提供的 staking.Delegate() 会自动处理：
	//   - 余额检查
	//   - 交易构建
	//   - 事件发出
	//
	// ⚠️ 注意：实际应用中需要业务规则检查
	//   验证者有效性、最小委托数量、委托关系管理等应在应用层实现
	caller := framework.GetCaller()
	err := staking.Delegate(caller, validator, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤4：发出委托创建事件（自定义事件，包含被委托者信息）
	event := framework.NewEvent("DelegationCreated")
	event.AddAddressField("delegator", caller)
	event.AddAddressField("validator", validator)
	event.AddAddressField("delegatee", delegatee)
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Undelegate 取消委托
//
// 使用 helpers/staking 模块的 Undelegate 函数取消委托。
// 支持部分取消委托或全部取消委托。
// SDK 内部会自动处理：
//   - 委托关系检查（确保存在委托关系）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Undelegate 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "validator": "validator_address",  // 验证者地址（Base58编码，必填）
//	  "amount": 5000                    // 取消委托数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析验证者地址
//  3. 查询委托关系（从状态输出）
//  4. 调用 staking.Undelegate() 取消委托
//     - SDK 内部自动处理委托关系检查
//     - SDK 内部自动构建交易
//  5. 发出取消委托事件（自定义事件）
//  6. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - 委托关系存在性检查
//   - 取消委托数量验证（不能超过已委托数量）
//   - 锁定期检查（业务逻辑）
//
// 返回：
//   - framework.SUCCESS - 取消委托成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_NOT_FOUND - 委托关系不存在
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Undelegate - 取消委托事件（由 SDK 自动发出）
//     {
//       "delegator": "<委托者地址>",
//       "validator": "<验证者地址>",
//       "amount": 5000
//     }
//   - DelegationCancelled - 委托取消事件（自定义）
//     {
//       "delegator": "<委托者地址>",
//       "validator": "<验证者地址>",
//       "amount": 5000
//     }
//
//export Undelegate
func Undelegate() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	validatorStr := params.ParseJSON("validator")
	amount := params.ParseJSONInt("amount")

	if validatorStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析验证者地址
	validator, err := framework.ParseAddressBase58(validatorStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：查询委托关系
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从状态输出查询委托关系
	//   检查委托关系是否存在，以及已委托数量

	// 步骤4：使用 SDK 基础能力取消委托
	//
	// SDK 提供的 staking.Undelegate() 会自动处理：
	//   - 委托关系检查
	//   - 交易构建
	//   - 事件发出
	//
	// ⚠️ 注意：实际应用中需要业务规则检查
	//   委托关系存在性、取消委托数量、锁定期等应在应用层实现
	caller := framework.GetCaller()
	err = staking.Undelegate(caller, validator, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤5：发出委托取消事件（自定义事件）
	event := framework.NewEvent("DelegationCancelled")
	event.AddAddressField("delegator", caller)
	event.AddAddressField("validator", validator)
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// QueryDelegation 查询委托信息
//
// 查询委托关系的详细信息，包括委托数量、委托收益等。
//
// 参数格式（JSON）:
//
//	{
//	  "validator": "validator_address"  // 验证者地址（Base58编码，必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 查询委托关系（从状态输出）
//  3. 计算委托收益（业务逻辑）
//  4. 返回查询结果
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中，应该从状态输出查询委托关系信息
//   包括委托数量、委托时间、委托收益等
//
// 返回：
//   - framework.SUCCESS - 查询成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_NOT_FOUND - 委托关系不存在
//
//export QueryDelegation
func QueryDelegation() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	validatorStr := params.ParseJSON("validator")

	if validatorStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析验证者地址
	validator, err := framework.ParseAddressBase58(validatorStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：查询委托关系
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从状态输出查询委托关系信息
	//   包括委托数量、委托时间、委托收益等
	caller := framework.GetCaller()

	// 简化实现：查询调用者的委托数量
	// 实际应该从状态输出查询
	delegatedAmount := framework.QueryUTXOBalance(caller, framework.TokenID(""))
	if delegatedAmount == 0 {
		return framework.ERROR_NOT_FOUND
	}

	// 步骤4：返回查询结果
	// 注意：实际应用中应该返回完整的委托信息
	result := `{"delegator":"` + caller.ToString() + `","validator":"` + validator.ToString() + `","amount":` + framework.Uint64ToString(uint64(delegatedAmount)) + `}`
	framework.SetReturnData([]byte(result))

	return framework.SUCCESS
}

func main() {}

