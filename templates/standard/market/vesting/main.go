//go:build tinygo || (js && wasm)

// Package main 提供分阶段释放合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建分阶段释放（Vesting）合约。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/market 模块进行分阶段释放
//   - 如何实现代币的分阶段解锁和释放
//   - 如何管理释放计划和时间表
//
// 🎯 核心功能
//
//  1. CreateVesting - 创建释放计划
//     - 使用 market.Release() 创建分阶段释放计划
//     - 支持设置释放时间表和释放条件
//
//  2. ClaimVesting - 领取释放的代币
//     - 根据释放计划领取已解锁的代币
//     - 自动检查释放条件和时间
//
//  3. QueryVesting - 查询释放计划
//     - 查询释放计划的详细信息
//     - 查询已释放和待释放的代币数量
//
// 📚 相关文档
//
//   - [Market 模块文档](../../helpers/market/README.md)
//   - [Framework 文档](../../framework/README.md)
//   - [示例总览](../README.md)
package main

import (
	"github.com/weisyn/contract-sdk-go/helpers/market"
	"github.com/weisyn/contract-sdk-go/framework"
)

// VestingContract 分阶段释放合约
//
// 本合约使用 helpers/market 模块提供的业务语义API，
// 简化分阶段释放操作的实现，开发者只需关注业务逻辑。
//
// 分阶段释放特点：
//   - 支持线性释放（Linear Vesting）
//   - 支持阶段性释放（Cliff Vesting）
//   - 支持自定义释放时间表
type VestingContract struct {
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
//       "contract": "Vesting",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "Vesting")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// CreateVesting 创建分阶段释放计划
//
// 使用 helpers/market 模块的 Release 函数创建分阶段释放计划。
// 适用于代币分配、员工激励、投资解锁等场景。
// SDK 内部会自动处理：
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Release 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "beneficiary": "beneficiary_address",  // 受益人地址（Base58编码，必填）
//	  "token_id": "TOKEN_001",              // 代币ID（可选，nil表示原生代币）
//	  "total_amount": 1000000,              // 总释放金额（必填）
//	  "vesting_id": "vesting_001",          // 释放计划ID（必填）
//	  "start_time": 1640995200,             // 开始时间（Unix时间戳，可选）
//	  "duration": 31536000                  // 释放持续时间（秒，可选）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析受益人地址
//  3. 调用 market.Release() 创建释放计划
//     - SDK 内部自动构建交易
//  4. 发出释放计划创建事件
//  5. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - 释放时间表验证
//   - 释放条件设置
//   - 权限检查（谁可以创建释放计划）
//
// 返回：
//   - framework.SUCCESS - 创建成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Release - 释放计划创建事件（由 SDK 自动发出）
//     {
//       "from": "<创建者地址>",
//       "beneficiary": "<受益人地址>",
//       "total_amount": 1000000,
//       "vesting_id": "vesting_001"
//     }
//   - VestingCreated - 释放计划创建事件（自定义）
//     {
//       "creator": "<创建者地址>",
//       "beneficiary": "<受益人地址>",
//       "total_amount": 1000000,
//       "vesting_id": "vesting_001"
//     }
//
//export CreateVesting
func CreateVesting() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	beneficiaryStr := params.ParseJSON("beneficiary")
	tokenIDStr := params.ParseJSON("token_id")
	totalAmount := params.ParseJSONInt("total_amount")
	vestingIDStr := params.ParseJSON("vesting_id")

	if beneficiaryStr == "" || totalAmount == 0 || vestingIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析受益人地址
	beneficiary, err := framework.ParseAddressBase58(beneficiaryStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：解析代币ID（可选）
	var tokenID framework.TokenID
	if tokenIDStr != "" {
		tokenID = framework.TokenID(tokenIDStr)
	}

	// 步骤4：使用 SDK 基础能力创建释放计划
	//
	// SDK 提供的 market.Release() 会自动处理：
	//   - 交易构建
	//   - 事件发出
	//
	// ⚠️ 注意：实际应用中需要业务规则检查
	//   释放时间表、释放条件、权限检查等应在应用层实现
	caller := framework.GetCaller()
	err = market.Release(
		caller,                        // 创建者地址
		beneficiary,                   // 受益人地址
		tokenID,                       // 代币ID
		framework.Amount(totalAmount), // 总释放金额
		[]byte(vestingIDStr),          // 释放计划ID
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤5：发出释放计划创建事件（自定义事件，包含更多信息）
	startTimeStr := params.ParseJSON("start_time")
	durationStr := params.ParseJSON("duration")

	event := framework.NewEvent("VestingCreated")
	event.AddAddressField("creator", caller)
	event.AddAddressField("beneficiary", beneficiary)
	event.AddStringField("vesting_id", vestingIDStr)
	event.AddUint64Field("total_amount", uint64(totalAmount))
	if tokenIDStr != "" {
		event.AddStringField("token_id", tokenIDStr)
	}
	if startTimeStr != "" {
		event.AddStringField("start_time", startTimeStr)
	}
	if durationStr != "" {
		event.AddStringField("duration", durationStr)
	}
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// ClaimVesting 领取释放的代币
//
// 根据释放计划领取已解锁的代币。
// 实际应用中，应该检查释放条件和时间，计算可领取数量。
//
// 参数格式（JSON）:
//
//	{
//	  "vesting_id": "vesting_001",    // 释放计划ID（必填）
//	  "amount": 10000                 // 领取数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 查询释放计划信息
//  3. 检查释放条件（时间、解锁比例等）
//  4. 计算可领取数量
//  5. 转移代币给受益人
//  6. 更新释放计划状态
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中，应该：
//   - 检查释放时间是否已到
//   - 计算已解锁的代币数量
//   - 检查是否已领取完毕
//   - 更新释放计划状态（使用状态输出）
//
// 返回：
//   - framework.SUCCESS - 领取成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - VestingClaimed - 代币领取事件
//     {
//       "beneficiary": "<受益人地址>",
//       "vesting_id": "vesting_001",
//       "amount": 10000
//     }
//
//export ClaimVesting
func ClaimVesting() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	vestingIDStr := params.ParseJSON("vesting_id")
	amount := params.ParseJSONInt("amount")

	if vestingIDStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：查询释放计划信息
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从状态输出查询释放计划信息
	//   检查释放时间、计算可领取数量等

	// 步骤3：检查释放条件
	// ⚠️ 注意：实际应用中需要实现
	//   - 检查释放时间是否已到
	//   - 计算已解锁的代币数量
	//   - 检查是否已领取完毕

	// 步骤4：转移代币给受益人
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从托管账户转移代币给受益人
	//   这里简化处理，实际应该使用 token.Transfer() 从托管账户转移

	// 步骤5：发出代币领取事件
	caller := framework.GetCaller()
	event := framework.NewEvent("VestingClaimed")
	event.AddAddressField("beneficiary", caller)
	event.AddStringField("vesting_id", vestingIDStr)
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// QueryVesting 查询释放计划
//
// 查询释放计划的详细信息，包括总金额、已释放金额、待释放金额等。
//
// 参数格式（JSON）:
//
//	{
//	  "vesting_id": "vesting_001"  // 释放计划ID（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 查询释放计划信息（从状态输出）
//  3. 计算已释放和待释放的代币数量
//  4. 返回查询结果
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中，应该从状态输出查询释放计划信息
//
// 返回：
//   - framework.SUCCESS - 查询成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_NOT_FOUND - 释放计划不存在
//
//export QueryVesting
func QueryVesting() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	vestingIDStr := params.ParseJSON("vesting_id")

	if vestingIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：查询释放计划信息
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该从状态输出查询释放计划信息
	//   包括总金额、已释放金额、待释放金额、释放时间表等

	// 步骤3：返回查询结果
	// 注意：实际应用中应该返回完整的释放计划信息
	result := `{"vesting_id":"` + vestingIDStr + `","status":"active"}`
	framework.SetReturnData([]byte(result))

	return framework.SUCCESS
}

func main() {}

