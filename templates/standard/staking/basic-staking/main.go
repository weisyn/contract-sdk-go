//go:build tinygo || (js && wasm)

// Package main 提供基础质押合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建质押和委托相关的智能合约。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/staking 模块进行质押和委托操作
//   - 如何使用业务语义API简化质押合约开发
//   - 如何实现完整的质押功能（Stake、Unstake、Delegate、Undelegate）
//
// 🎯 核心功能
//
//  1. Stake - 质押
//     - 使用 staking.Stake() 进行代币质押
//     - SDK 内部自动处理余额检查、交易构建、事件发出
//
//  2. Unstake - 解质押
//     - 使用 staking.Unstake() 解质押代币
//     - 支持部分解质押或全部解质押
//
//  3. Delegate - 委托
//     - 使用 staking.Delegate() 将质押权委托给验证者
//     - 适用于委托质押场景
//
//  4. Undelegate - 取消委托
//     - 使用 staking.Undelegate() 取消委托
//     - 支持部分取消委托或全部取消委托
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

// StakingContract 基础质押合约
//
// 本合约使用 helpers/staking 模块提供的业务语义API，
// 简化质押和委托操作的实现，开发者只需关注业务逻辑。
type StakingContract struct {
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
//       "contract": "Staking",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "Staking")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Stake 质押代币
//
// 使用 helpers/staking 模块的 Stake 函数进行代币质押。
// SDK 内部会自动处理：
//   - 余额检查（确保质押者余额充足）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Stake 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "validator": "validator_address", // 验证者地址（Base58编码，必填）
//	  "amount": 10000                  // 质押数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析验证者地址
//  3. 调用 staking.Stake() 进行质押
//     - SDK 内部自动处理余额检查
//     - SDK 内部自动构建交易
//  4. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - 验证者有效性检查（验证者是否在验证者列表中）
//   - 最小质押数量检查
//   - 锁定期检查（业务逻辑）
//
// 返回：
//   - framework.SUCCESS - 质押成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Stake - 质押事件（由 SDK 自动发出）
//     {
//       "staker": "<质押者地址>",
//       "validator": "<验证者地址>",
//       "amount": 10000
//     }
//
//export Stake
func Stake() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	validatorStr := params.ParseJSON("validator")
	amount := params.ParseJSONInt("amount")

	if validatorStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 解析验证者地址
	validator, err := framework.ParseAddressBase58(validatorStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：使用 SDK 基础能力进行代币质押
	//
	// SDK 提供的 staking.Stake() 会自动处理：
	//   - 余额检查
	//   - 交易构建
	//   - 事件发出
	//
	// ⚠️ 注意：实际应用中需要业务规则检查
	//   验证者有效性、最小质押数量、锁定期等应在应用层实现
	caller := framework.GetCaller()
	err = staking.Stake(caller, validator, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		// 检查错误类型
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Unstake 解质押
//
// 使用 helpers/staking 模块的 Unstake 函数解质押代币。
// SDK 内部会自动处理：
//   - 质押余额检查（确保有足够的质押余额）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Unstake 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "validator": "validator_address", // 验证者地址（Base58编码，必填）
//	  "amount": 5000                   // 解质押数量（可选，0表示全部解质押）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析验证者地址
//  3. 调用 staking.Unstake() 进行解质押
//     - SDK 内部自动处理质押余额检查
//     - SDK 内部自动构建交易
//  4. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - 锁定期检查（必须满足锁定期要求才能解质押）
//   - 解质押冷却期检查
//   - amount为0表示全部解质押
//
// 返回：
//   - framework.SUCCESS - 解质押成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 质押余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Unstake - 解质押事件（由 SDK 自动发出）
//     {
//       "staker": "<质押者地址>",
//       "validator": "<验证者地址>",
//       "amount": 5000
//     }
//
//export Unstake
func Unstake() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	validatorStr := params.ParseJSON("validator")
	amount := params.ParseJSONInt("amount")

	if validatorStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 解析验证者地址
	validator, err := framework.ParseAddressBase58(validatorStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 获取调用者
	caller := framework.GetCaller()

	// 使用helpers进行解质押
	// 注意：amount为0表示全部解质押
	err = staking.Unstake(caller, validator, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Delegate 委托质押
//
// 使用 helpers/staking 模块的 Delegate 函数将质押权委托给验证者。
// 适用于委托质押场景，允许用户将质押权委托给其他验证者。
// SDK 内部会自动处理：
//   - 余额检查（确保委托者余额充足）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Delegate 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "validator": "validator_address", // 验证者地址（Base58编码，必填）
//	  "amount": 5000                   // 委托数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析验证者地址
//  3. 调用 staking.Delegate() 进行委托
//     - SDK 内部自动处理余额检查
//     - SDK 内部自动构建交易
//  4. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - 验证者有效性检查
//   - 最小委托数量检查
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
//       "amount": 5000
//     }
//
//export Delegate
func Delegate() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	validatorStr := params.ParseJSON("validator")
	amount := params.ParseJSONInt("amount")

	if validatorStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 解析验证者地址
	validator, err := framework.ParseAddressBase58(validatorStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 获取调用者
	caller := framework.GetCaller()

	// 使用helpers进行委托
	err = staking.Delegate(caller, validator, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Undelegate 取消委托
//
// 使用 helpers/staking 模块的 Undelegate 函数取消委托。
// SDK 内部会自动处理：
//   - 委托余额检查（确保有足够的委托余额）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Undelegate 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "validator": "validator_address", // 验证者地址（Base58编码，必填）
//	  "amount": 2000                   // 取消委托数量（可选，0表示全部取消委托）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析验证者地址
//  3. 调用 staking.Undelegate() 进行取消委托
//     - SDK 内部自动处理委托余额检查
//     - SDK 内部自动构建交易
//  4. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - 锁定期检查（必须满足锁定期要求才能取消委托）
//   - 取消委托冷却期检查
//   - amount为0表示全部取消委托
//
// 返回：
//   - framework.SUCCESS - 取消委托成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 委托余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Undelegate - 取消委托事件（由 SDK 自动发出）
//     {
//       "delegator": "<委托者地址>",
//       "validator": "<验证者地址>",
//       "amount": 2000
//     }
//
//export Undelegate
func Undelegate() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	validatorStr := params.ParseJSON("validator")
	amount := params.ParseJSONInt("amount")

	if validatorStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 解析验证者地址
	validator, err := framework.ParseAddressBase58(validatorStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 获取调用者
	caller := framework.GetCaller()

	// 使用helpers进行取消委托
	// 注意：amount为0表示全部取消委托
	err = staking.Undelegate(caller, validator, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

func main() {}

