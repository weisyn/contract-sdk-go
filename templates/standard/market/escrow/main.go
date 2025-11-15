//go:build tinygo || (js && wasm)

// Package main 提供市场托管合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建市场交易相关的智能合约。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/market 模块进行市场操作
//   - 如何使用业务语义API简化市场合约开发
//   - 如何实现完整的市场功能（Escrow、Release）
//
// 🎯 核心功能
//
//  1. Escrow - 托管
//     - 使用 market.Escrow() 创建代币托管
//     - SDK 内部自动处理余额检查、交易构建、事件发出
//
//  2. Release - 分阶段释放
//     - 使用 market.Release() 创建分阶段释放计划
//     - SDK 内部自动处理交易构建、事件发出
//
// ⚠️ 注意：本模块仅提供原子操作，不包含组合场景（如Swap、Liquidity等）
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

// MarketContract 市场托管合约
//
// 本合约使用 helpers/market 模块提供的业务语义API，
// 简化市场操作的实现，开发者只需关注业务逻辑。
type MarketContract struct {
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
//       "contract": "Market",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "Market")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Escrow 创建托管
//
// 使用 helpers/market 模块的 Escrow 函数创建代币托管。
// 适用于交易托管、质押托管等场景。
// SDK 内部会自动处理：
//   - 余额检查（确保买方余额充足）
//   - 交易构建（自动构建 UTXO 交易）
//   - 资产锁定（自动锁定托管资产）
//   - 事件发出（自动发出 Escrow 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "buyer": "buyer_address",      // 买方地址（Base58编码，必填）
//	  "seller": "seller_address",    // 卖方地址（Base58编码，必填）
//	  "amount": 10000,               // 托管数量（必填）
//	  "escrow_id": "escrow_123"      // 托管ID（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析买方和卖方地址
//  3. 调用 market.Escrow() 创建托管
//     - SDK 内部自动处理余额检查
//     - SDK 内部自动构建交易
//     - SDK 内部自动锁定资产
//  4. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - 托管条件验证（交易条件、质押条件等）
//   - 托管金额限制
//   - 托管ID唯一性检查
//
// 返回：
//   - framework.SUCCESS - 托管创建成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Escrow - 托管事件（由 SDK 自动发出）
//     {
//       "buyer": "<买方地址>",
//       "seller": "<卖方地址>",
//       "amount": 10000,
//       "escrow_id": "escrow_123"
//     }
//
//export Escrow
func Escrow() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	buyerStr := params.ParseJSON("buyer")
	sellerStr := params.ParseJSON("seller")
	amount := params.ParseJSONInt("amount")
	escrowIDStr := params.ParseJSON("escrow_id")

	if buyerStr == "" || sellerStr == "" || amount == 0 || escrowIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析买方和卖方地址
	buyer, err1 := framework.ParseAddressBase58(buyerStr)
	seller, err2 := framework.ParseAddressBase58(sellerStr)
	if err1 != nil || err2 != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：使用 SDK 基础能力创建托管
	//
	// SDK 提供的 market.Escrow() 会自动处理：
	//   - 余额检查
	//   - 交易构建
	//   - 资产锁定
	//   - 事件发出
	//
	// ⚠️ 注意：实际应用中需要业务规则检查
	//   托管条件验证、托管金额限制、托管ID唯一性等应在应用层实现
	err := market.Escrow(
		buyer,
		seller,
		framework.TokenID(""), // 原生币（空字符串表示使用原生币）
		framework.Amount(amount),
		[]byte(escrowIDStr),
	)
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Release 创建分阶段释放计划
//
// 使用 helpers/market 模块的 Release 函数创建分阶段释放计划。
// 适用于分红释放、租金分配、代币解锁等场景。
// SDK 内部会自动处理：
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Release 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "from": "from_address",          // 分配者地址（Base58编码，必填）
//	  "beneficiary": "beneficiary_address", // 受益人地址（Base58编码，必填）
//	  "total_amount": 100000,         // 总释放金额（必填）
//	  "vesting_id": "vesting_123"      // 释放计划ID（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析分配者和受益人地址
//  3. 调用 market.Release() 创建分阶段释放计划
//     - SDK 内部自动构建交易
//  4. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - 释放条件验证（时间锁、高度锁等）
//   - 释放计划ID唯一性检查
//   - 分阶段释放逻辑（线性释放、阶梯释放等）需要在合约中实现
//
// 返回：
//   - framework.SUCCESS - 释放计划创建成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Release - 释放计划事件（由 SDK 自动发出）
//     {
//       "from": "<分配者地址>",
//       "beneficiary": "<受益人地址>",
//       "total_amount": 100000,
//       "vesting_id": "vesting_123"
//     }
//
//export Release
func Release() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	fromStr := params.ParseJSON("from")
	beneficiaryStr := params.ParseJSON("beneficiary")
	totalAmount := params.ParseJSONInt("total_amount")
	vestingIDStr := params.ParseJSON("vesting_id")

	if fromStr == "" || beneficiaryStr == "" || totalAmount == 0 || vestingIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析分配者和受益人地址
	from, err1 := framework.ParseAddressBase58(fromStr)
	beneficiary, err2 := framework.ParseAddressBase58(beneficiaryStr)
	if err1 != nil || err2 != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：使用 SDK 基础能力创建分阶段释放计划
	//
	// SDK 提供的 market.Release() 会自动处理：
	//   - 交易构建
	//   - 事件发出
	//
	// ⚠️ 注意：实际应用中需要业务规则检查
	//   释放条件验证、释放计划ID唯一性、分阶段释放逻辑等应在应用层实现
	err := market.Release(
		from,
		beneficiary,
		framework.TokenID(""), // 原生币（空字符串表示使用原生币）
		framework.Amount(totalAmount),
		[]byte(vestingIDStr),
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

