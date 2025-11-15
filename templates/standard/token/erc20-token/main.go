//go:build tinygo || (js && wasm)

// Package main 提供 ERC-20 兼容代币合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建 ERC-20 兼容的代币合约。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/token 模块进行代币操作
//   - 如何使用业务语义API简化合约开发
//   - 如何实现完整的代币功能（Transfer、Mint、Burn、Approve、Freeze、Airdrop）
//
// 🎯 核心功能
//
//  1. Transfer - 转账
//     - 使用 token.Transfer() 进行代币转账
//     - SDK 内部自动处理余额检查、交易构建、事件发出
//
//  2. Mint - 铸造
//     - 使用 token.Mint() 铸造新代币
//     - 支持向指定地址铸造指定数量代币
//
//  3. Burn - 销毁
//     - 使用 token.Burn() 销毁代币
//     - 从调用者地址销毁指定数量代币
//
//  4. Approve - 授权
//     - 使用 token.Approve() 授权其他地址使用代币
//     - 支持 ERC-20 风格的授权机制
//
//  5. Freeze - 冻结
//     - 使用 token.Freeze() 冻结指定地址的代币
//     - 适用于合规、风控等场景
//
//  6. Airdrop - 空投
//     - 使用 token.Airdrop() 批量空投代币
//     - 支持一次性向多个地址空投不同数量代币
//
// 📚 相关文档
//
//   - [Token 模块文档](../../helpers/token/README.md)
//   - [Framework 文档](../../framework/README.md)
//   - [示例总览](../README.md)
package main

import (
	"github.com/weisyn/contract-sdk-go/helpers/token"
	"github.com/weisyn/contract-sdk-go/framework"
)

// TokenContract ERC-20 兼容代币合约
//
// 本合约使用 helpers/token 模块提供的业务语义API，
// 简化代币操作的实现，开发者只需关注业务逻辑。
type TokenContract struct {
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
//       "contract": "Token",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "Token")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// Transfer 转账代币
//
// 使用 helpers/token 模块的 Transfer 函数进行代币转账。
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
//	  "amount": 100                 // 转账数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析接收者地址
//  3. 调用 token.Transfer() 进行转账
//     - SDK 内部自动处理余额检查
//     - SDK 内部自动构建交易
//     - SDK 内部自动处理找零
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
	// 获取参数
	params := framework.GetContractParams()
	toStr := params.ParseJSON("to")
	amount := params.ParseJSONInt("amount")

	if toStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 解析地址
	to, err := framework.ParseAddressBase58(toStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 获取调用者
	caller := framework.GetCaller()

	// 使用helpers进行转账
	err = token.Transfer(caller, to, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		// 检查错误类型
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Mint 铸造代币
//
// 使用 helpers/token 模块的 Mint 函数铸造新代币。
// SDK 内部会自动处理：
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Mint 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "to": "receiver_address",    // 接收者地址（Base58编码，必填）
//	  "amount": 1000                // 铸造数量（必填）
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
	// 获取参数
	params := framework.GetContractParams()
	toStr := params.ParseJSON("to")
	amount := params.ParseJSONInt("amount")

	if toStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 解析地址
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

// Burn 销毁代币
//
// 使用 helpers/token 模块的 Burn 函数销毁代币。
// SDK 内部会自动处理：
//   - 余额检查（确保调用者余额充足）
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Burn 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "amount": 500                // 销毁数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 调用 token.Burn() 进行销毁
//     - SDK 内部自动处理余额检查
//     - SDK 内部自动构建交易
//  3. 返回执行结果
//
// 返回：
//   - framework.SUCCESS - 销毁成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Burn - 销毁事件（由 SDK 自动发出）
//     {
//       "from": "<销毁者地址>",
//       "amount": 500
//     }
//
//export Burn
func Burn() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	amount := params.ParseJSONInt("amount")

	if amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 获取调用者
	caller := framework.GetCaller()

	// 使用helpers进行销毁
	err := token.Burn(caller, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Approve 授权代币
//
// 使用 helpers/token 模块的 Approve 函数授权其他地址使用代币。
// 支持 ERC-20 风格的授权机制，允许授权地址代表所有者进行转账。
// SDK 内部会自动处理：
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Approve 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "spender": "spender_address", // 被授权地址（Base58编码，必填）
//	  "amount": 1000                // 授权数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析被授权地址
//  3. 调用 token.Approve() 进行授权
//     - SDK 内部自动构建交易
//  4. 返回执行结果
//
// 返回：
//   - framework.SUCCESS - 授权成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Approve - 授权事件（由 SDK 自动发出）
//     {
//       "owner": "<所有者地址>",
//       "spender": "<被授权地址>",
//       "amount": 1000
//     }
//
//export Approve
func Approve() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	spenderStr := params.ParseJSON("spender")
	amount := params.ParseJSONInt("amount")

	if spenderStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 解析地址
	spender, err := framework.ParseAddressBase58(spenderStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 获取调用者
	caller := framework.GetCaller()

	// 使用helpers进行授权
	err = token.Approve(caller, spender, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Airdrop 批量空投代币
//
// 使用 helpers/token 模块的 Airdrop 函数批量空投代币。
// 支持一次性向多个地址空投不同数量的代币，适用于：
//   - 代币分发活动
//   - 社区奖励发放
//   - 批量转账场景
//
// SDK 内部会自动处理：
//   - 批量交易构建（自动构建多个 UTXO 交易）
//   - 事件发出（自动发出多个 Transfer 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "recipients": ["addr1", "addr2"],  // 接收者地址数组（Base58编码，必填）
//	  "amounts": [100, 200]             // 对应数量数组（必填，长度需与recipients一致）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析接收者地址数组和数量数组
//  3. 验证数组长度一致
//  4. 构建 AirdropRecipient 列表
//  5. 调用 token.Airdrop() 进行批量空投
//     - SDK 内部自动构建多个交易
//  6. 返回执行结果
//
// ⚠️ 注意：
//   - 本示例使用简化的 JSON 解析，实际应用中应使用完整的 JSON 解析库
//   - 批量空投可能涉及大量交易，需要注意 Gas 费用
//
// 返回：
//   - framework.SUCCESS - 空投成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Transfer - 转账事件（由 SDK 自动发出，每个接收者一个事件）
//     {
//       "from": "<发送者地址>",
//       "to": "<接收者地址>",
//       "amount": 100
//     }
//
//export Airdrop
func Airdrop() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	recipientsStr := params.ParseJSON("recipients")
	amountsStr := params.ParseJSON("amounts")

	if recipientsStr == "" || amountsStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 解析接收者地址数组
	recipientAddrs := parseJSONArray(recipientsStr)
	if len(recipientAddrs) == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 解析金额数组
	amounts := parseJSONIntArray(amountsStr)
	if len(amounts) == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 验证数组长度一致
	if len(recipientAddrs) != len(amounts) {
		return framework.ERROR_INVALID_PARAMS
	}

	// 获取调用者
	caller := framework.GetCaller()

	// 构建接收者列表
	recipients := make([]token.AirdropRecipient, len(recipientAddrs))
	for i := 0; i < len(recipientAddrs); i++ {
		addr, err := framework.ParseAddressBase58(recipientAddrs[i])
		if err != nil {
			return framework.ERROR_INVALID_PARAMS
		}
		recipients[i] = token.AirdropRecipient{
			Address: addr,
			Amount:  framework.Amount(amounts[i]),
		}
	}

	// 使用helpers进行空投
	err := token.Airdrop(caller, recipients, framework.TokenID(""))
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Freeze 冻结代币
//
// 使用 helpers/token 模块的 Freeze 函数冻结指定地址的代币。
// 适用于合规、风控等场景，可以临时冻结特定地址的代币，防止转移。
// SDK 内部会自动处理：
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Freeze 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "target": "target_address",  // 目标地址（Base58编码，必填）
//	  "amount": 1000               // 冻结数量（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析目标地址
//  3. 调用 token.Freeze() 进行冻结
//     - SDK 内部自动构建交易
//  4. 返回执行结果
//
// ⚠️ 注意：实际应用中需要权限检查
//   - 只有授权地址才能调用 Freeze
//   - 权限检查逻辑应在应用层实现
//
// 返回：
//   - framework.SUCCESS - 冻结成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Freeze - 冻结事件（由 SDK 自动发出）
//     {
//       "target": "<目标地址>",
//       "amount": 1000
//     }
//
//export Freeze
func Freeze() uint32 {
	// 获取参数
	params := framework.GetContractParams()
	targetStr := params.ParseJSON("target")
	amount := params.ParseJSONInt("amount")

	if targetStr == "" || amount == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 解析地址
	target, err := framework.ParseAddressBase58(targetStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 使用helpers进行冻结
	err = token.Freeze(target, framework.TokenID(""), framework.Amount(amount))
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// parseJSONArray 解析JSON字符串数组
func parseJSONArray(jsonStr string) []string {
	jsonStr = trimSpace(jsonStr)
	if len(jsonStr) < 2 || jsonStr[0] != '[' || jsonStr[len(jsonStr)-1] != ']' {
		return nil
	}

	jsonStr = jsonStr[1 : len(jsonStr)-1]
	if len(jsonStr) == 0 {
		return []string{}
	}

	var result []string
	start := 0
	inString := false

	for i := 0; i < len(jsonStr); i++ {
		c := jsonStr[i]
		if c == '"' {
			inString = !inString
		} else if c == ',' && !inString {
			item := jsonStr[start:i]
			item = trimSpace(item)
			if len(item) >= 2 && item[0] == '"' && item[len(item)-1] == '"' {
				item = item[1 : len(item)-1]
			}
			if len(item) > 0 {
				result = append(result, item)
			}
			start = i + 1
		}
	}

	if start < len(jsonStr) {
		item := jsonStr[start:]
		item = trimSpace(item)
		if len(item) >= 2 && item[0] == '"' && item[len(item)-1] == '"' {
			item = item[1 : len(item)-1]
		}
		if len(item) > 0 {
			result = append(result, item)
		}
	}

	return result
}

// parseJSONIntArray 解析JSON整数数组
func parseJSONIntArray(jsonStr string) []uint64 {
	jsonStr = trimSpace(jsonStr)
	if len(jsonStr) < 2 || jsonStr[0] != '[' || jsonStr[len(jsonStr)-1] != ']' {
		return nil
	}

	jsonStr = jsonStr[1 : len(jsonStr)-1]
	if len(jsonStr) == 0 {
		return []uint64{}
	}

	var result []uint64
	start := 0

	for i := 0; i < len(jsonStr); i++ {
		c := jsonStr[i]
		if c == ',' {
			item := trimSpace(jsonStr[start:i])
			if len(item) > 0 {
				amount := parseUint64(item)
				result = append(result, amount)
			}
			start = i + 1
		}
	}

	if start < len(jsonStr) {
		item := trimSpace(jsonStr[start:])
		if len(item) > 0 {
			amount := parseUint64(item)
			result = append(result, amount)
		}
	}

	return result
}

// trimSpace 移除字符串两端的空格
func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}

	return s[start:end]
}

// parseUint64 解析无符号整数
func parseUint64(s string) uint64 {
	var result uint64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			result = result*10 + uint64(c-'0')
		} else {
			break
		}
	}
	return result
}

func main() {}

