//go:build tinygo || (js && wasm)

// Package main 提供门票票务NFT合约示例
//
// 📋 示例说明
//
// 本示例展示如何使用 WES Contract SDK Go 构建门票票务NFT合约。
// 通过本示例，您可以学习：
//   - 如何使用 helpers/token 模块创建NFT
//   - 如何管理NFT元数据
//   - 如何实现NFT的铸造、转移、查询等功能
//
// 🎯 核心功能
//
//  1. MintNFT - 铸造NFT
//     - 使用 token.Mint() 铸造唯一的数字艺术NFT
//     - 每个NFT都有唯一的tokenID和元数据
//
//  2. TransferNFT - 转移NFT
//     - 使用 token.Transfer() 转移NFT所有权
//     - NFT是唯一的，转移时数量为1
//
//  3. QueryNFT - 查询NFT信息
//     - 查询NFT的所有者
//     - 查询NFT的元数据
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

// DigitalArtNFTContract 数字艺术NFT合约
//
// 本合约使用 helpers/token 模块提供的业务语义API，
// 简化NFT操作的实现，开发者只需关注业务逻辑。
//
// NFT特点：
//   - 每个NFT都有唯一的tokenID
//   - 每个tokenID只能有一个所有者
//   - NFT不可分割，转移时数量为1
type TicketsNFTContract struct {
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
//       "contract": "TicketsNFT",
//       "owner": "<合约所有者地址>"
//     }
//
//export Initialize
func Initialize() uint32 {
	caller := framework.GetCaller()
	event := framework.NewEvent("ContractInitialized")
	event.AddStringField("contract", "DigitalArtNFT")
	event.AddAddressField("owner", caller)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// MintNFT 铸造数字艺术NFT
//
// 使用 helpers/token 模块的 Mint 函数铸造唯一的数字艺术NFT。
// 每个NFT都有唯一的tokenID和元数据（艺术品名称、作者、描述、图片URL等）。
// SDK 内部会自动处理：
//   - 交易构建（自动构建 UTXO 交易）
//   - 事件发出（自动发出 Mint 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "to": "receiver_address",        // 接收者地址（Base58编码，必填）
//	  "token_id": "art_001",           // NFT的tokenID（必填，唯一标识）
//	  "ticket_name": "Sunset Over Mountains", // 艺术品名称（必填）
//	  "event": "Alice",               // 艺术家名称（必填）
//	  "description": "A beautiful...", // 艺术品描述（可选）
//	  "image_url": "https://..."       // 图片URL（可选）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析接收者地址
//  3. 验证tokenID唯一性（检查是否已存在）
//  4. 调用 token.Mint() 铸造NFT
//     - SDK 内部自动构建交易
//  5. 发出NFT铸造事件（包含元数据）
//  6. 返回执行结果
//
// ⚠️ 注意：实际应用中需要业务规则检查
//   - tokenID唯一性检查（确保每个NFT唯一）
//   - 元数据格式验证
//   - 铸造权限检查（谁可以铸造NFT）
//
// 返回：
//   - framework.SUCCESS - 铸造成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_ALREADY_EXISTS - NFT已存在
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - NFTMinted - NFT铸造事件
//     {
//       "to": "<接收者地址>",
//       "token_id": "art_001",
//       "ticket_name": "Sunset Over Mountains",
//       "event": "Alice"
//     }
//
//export MintNFT
func MintNFT() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	toStr := params.ParseJSON("to")
	tokenIDStr := params.ParseJSON("token_id")
	ticket_nameStr := params.ParseJSON("ticket_name")
	eventStr := params.ParseJSON("event")

	if toStr == "" || tokenIDStr == "" || ticket_nameStr == "" || eventStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析接收者地址
	to, err := framework.ParseAddressBase58(toStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：验证tokenID唯一性
	// ⚠️ 注意：实际应用中需要检查NFT是否已存在
	// 这里简化处理，实际应该查询链上状态
	tokenID := framework.TokenID(tokenIDStr)
	balance := framework.QueryUTXOBalance(to, tokenID)
	if balance > 0 {
		return framework.ERROR_ALREADY_EXISTS
	}

	// 步骤4：使用 SDK 基础能力铸造NFT
	//
	// SDK 提供的 token.Mint() 会自动处理：
	//   - 交易构建
	//   - 事件发出（Mint事件）
	//
	// ⚠️ 注意：实际应用中需要业务规则检查
	//   tokenID唯一性、元数据格式验证、铸造权限等应在应用层实现
	err = token.Mint(to, tokenID, framework.Amount(1)) // NFT数量为1
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤5：发出NFT铸造事件（包含元数据）
	descriptionStr := params.ParseJSON("description")
	imageURLStr := params.ParseJSON("image_url")

	event := framework.NewEvent("NFTMinted")
	event.AddAddressField("to", to)
	event.AddStringField("token_id", tokenIDStr)
	event.AddStringField("ticket_name", ticket_nameStr)
	event.AddStringField("event", eventStr)
	if descriptionStr != "" {
		event.AddStringField("description", descriptionStr)
	}
	if imageURLStr != "" {
		event.AddStringField("image_url", imageURLStr)
	}
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// TransferNFT 转移NFT
//
// 使用 helpers/token 模块的 Transfer 函数转移NFT所有权。
// NFT是唯一的，转移时数量为1。
// SDK 内部会自动处理：
//   - 余额检查（确保发送者拥有该NFT）
//   - 交易构建（自动构建 UTXO 交易）
//   - 找零处理（自动处理找零 UTXO）
//   - 事件发出（自动发出 Transfer 事件）
//
// 参数格式（JSON）:
//
//	{
//	  "to": "receiver_address",  // 接收者地址（Base58编码，必填）
//	  "token_id": "art_001"      // NFT的tokenID（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 解析接收者地址
//  3. 调用 token.Transfer() 转移NFT
//     - SDK 内部自动处理余额检查
//     - SDK 内部自动构建交易
//  4. 发出NFT转移事件
//  5. 返回执行结果
//
// 返回：
//   - framework.SUCCESS - 转移成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE - 余额不足（不拥有该NFT）
//   - framework.ERROR_EXECUTION_FAILED - 执行失败
//
// 事件：
//   - Transfer - 转账事件（由 SDK 自动发出）
//     {
//       "from": "<发送者地址>",
//       "to": "<接收者地址>",
//       "token_id": "art_001",
//       "amount": 1
//     }
//   - NFTTransferred - NFT转移事件（自定义）
//     {
//       "from": "<发送者地址>",
//       "to": "<接收者地址>",
//       "token_id": "art_001"
//     }
//
//export TransferNFT
func TransferNFT() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	toStr := params.ParseJSON("to")
	tokenIDStr := params.ParseJSON("token_id")

	if toStr == "" || tokenIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：解析接收者地址
	to, err := framework.ParseAddressBase58(toStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：使用 SDK 基础能力转移NFT
	//
	// SDK 提供的 token.Transfer() 会自动处理：
	//   - 余额检查（确保发送者拥有该NFT）
	//   - 交易构建
	//   - 事件发出（Transfer事件）
	caller := framework.GetCaller()
	tokenID := framework.TokenID(tokenIDStr)
	err = token.Transfer(caller, to, tokenID, framework.Amount(1)) // NFT数量为1
	if err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤4：发出NFT转移事件（自定义事件，包含更多信息）
	event := framework.NewEvent("NFTTransferred")
	event.AddAddressField("from", caller)
	event.AddAddressField("to", to)
	event.AddStringField("token_id", tokenIDStr)
	framework.EmitEvent(event)

	return framework.SUCCESS
}

// QueryNFT 查询NFT信息
//
// 查询NFT的所有者信息。
// 注意：这是一个查询函数，不会修改链上状态。
//
// 参数格式（JSON）:
//
//	{
//	  "token_id": "art_001"  // NFT的tokenID（必填）
//	}
//
// 工作流程：
//  1. 解析参数并验证
//  2. 查询NFT余额（找到拥有该NFT的地址）
//  3. 返回查询结果
//
// ⚠️ 注意：这是一个简化实现
//   实际应用中，应该遍历所有地址查找拥有该NFT的地址
//   或者使用状态输出存储NFT所有权映射
//
// 返回：
//   - framework.SUCCESS - 查询成功
//   - framework.ERROR_INVALID_PARAMS - 参数无效
//   - framework.ERROR_NOT_FOUND - NFT不存在
//
//export QueryNFT
func QueryNFT() uint32 {
	// 步骤1：解析参数并验证
	params := framework.GetContractParams()
	tokenIDStr := params.ParseJSON("token_id")

	if tokenIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：查询NFT信息
	// ⚠️ 注意：这是一个简化实现
	//   实际应用中，应该使用状态输出存储NFT所有权映射
	//   或者遍历所有地址查找拥有该NFT的地址
	tokenID := framework.TokenID(tokenIDStr)
	caller := framework.GetCaller()

	// 简化实现：查询调用者的余额
	balance := framework.QueryUTXOBalance(caller, tokenID)
	if balance == 0 {
		return framework.ERROR_NOT_FOUND
	}

	// 步骤3：返回查询结果
	// 注意：实际应用中应该返回NFT的完整信息（元数据等）
	result := `{"token_id":"` + tokenIDStr + `","owner":"` + caller.ToString() + `","balance":1}`
	framework.SetReturnData([]byte(result))

	return framework.SUCCESS
}

func main() {}

