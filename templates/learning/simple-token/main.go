//go:build tinygo || (js && wasm)

// Package main 提供最简单的代币合约示例 - Simple Token
//
// 📋 示例说明
//
// 本示例展示如何构建一个最简单的代币合约，实现基本的代币功能。
// 通过本示例，您可以学习：
//   - 如何创建代币合约
//   - 如何实现代币的铸造、转账、查询功能
//   - 如何使用状态存储代币余额和总供应量
//   - 如何处理参数解析和错误处理
//   - 如何发出代币相关事件
//
// 🎯 核心功能
//
//  1. Initialize - 初始化合约
//     - 设置初始代币供应量
//     - 将初始供应量分配给部署者
//
//  2. Mint - 铸造代币
//     - 向调用者铸造固定数量的代币
//     - 更新总供应量
//
//  3. Transfer - 转账
//     - 从调用者向指定地址转账
//     - 检查余额是否充足
//     - 更新发送者和接收者的余额
//
//  4. BalanceOf - 查询余额
//     - 查询指定地址的代币余额
//     - 返回 JSON 格式的结果
//
//  5. TotalSupply - 查询总供应量
//     - 查询代币的总供应量
//     - 返回 JSON 格式的结果
//
// 📚 学习要点
//
//   - **UTXO模型**：使用 SDK 的 `token.Mint()` 和 `token.Transfer()` 创建和管理代币UTXO
//   - **余额查询**：使用 `framework.QueryUTXOBalance()` 从UTXO集合查询余额
//   - **参数解析**：使用 `framework.GetContractParams()` 获取 JSON 格式的参数
//   - **错误处理**：检查参数有效性、余额是否充足等
//   - **事件发出**：SDK 自动发出代币相关事件
//
// ✅ 最佳实践
//
//   - 本示例使用 `helpers/token` 层的业务语义接口
//   - SDK 自动处理：交易构建、UTXO创建、事件发出
//   - 合约只需关注业务逻辑：参数解析、权限检查
//   - 参考示例：`examples/token/erc20-token/` - 展示了更多代币功能
//
// 📚 相关文档
//
//   - [Token 模块文档](../../helpers/token/README.md) - Token 模块详细说明
//   - [Framework 文档](../../framework/README.md) - Framework 层详细说明
//   - [示例总览](../README.md) - 所有示例索引
//   - [Simple Token 示例 README](./README.md) - 本示例详细文档
package main

import (
	"encoding/json"
	"strconv"

	"github.com/weisyn/contract-sdk-go/framework"
	"github.com/weisyn/contract-sdk-go/helpers/token"
)

// SimpleToken 最小代币合约
//
// 本合约展示了如何实现一个最简单的代币合约，包含基本的代币功能。
// 所有 WES 合约都需要嵌入 `framework.ContractBase`，以获得基础功能。
//
// 设计理念：
//   - 简单易懂：使用最简单的功能展示代币合约的基本结构
//   - 完整示例：包含初始化、铸造、转账、查询等基本操作
//   - 学习友好：适合初学者理解代币合约的基本概念
//
// ✅ **最佳实践**：
//   - 本示例使用 helpers 层的业务语义接口
//   - token.Mint() - 铸造代币，创建AssetOutput（UTXO输出）
//   - token.Transfer() - 转账代币，自动处理余额检查和交易构建
//   - framework.QueryUTXOBalance() - 查询UTXO余额
//
// 参考示例：
//   - `examples/token/erc20-token/` - 展示了更多代币功能（授权、冻结、空投等）
type SimpleToken struct {
	framework.ContractBase
}

// InitParams 初始化参数结构体
//
// 用于解析合约初始化时的 JSON 参数。
// 实际项目中，初始化参数可以通过 `framework.GetContractParams()` 获取。
type InitParams struct {
	InitialSupply string `json:"initialSupply"` // 初始供应量（字符串格式）
}

// TransferParams 转账参数结构体
//
// 用于解析转账函数的 JSON 参数。
type TransferParams struct {
	To     string `json:"to"`     // 接收者地址
	Amount string `json:"amount"` // 转账金额（字符串格式）
}

// BalanceQuery 余额查询参数结构体
//
// 用于解析余额查询函数的 JSON 参数。
type BalanceQuery struct {
	Address string `json:"address"` // 要查询的地址
}

// Initialize 初始化合约
//
// 🎯 **用途**：在合约部署时自动调用，初始化代币合约
//
// **调用时机**：
//   - 合约部署时自动调用一次
//   - 只会被调用一次，用于设置初始状态
//
// **工作流程**：
//  1. 检查 ABI 版本兼容性
//  2. 使用 SDK 的 token.Mint() 向部署者铸造初始代币
//     - SDK 内部自动创建 AssetOutput（UTXO输出）
//     - SDK 内部自动发出铸造事件
//  3. 发出初始化事件
//
// **参数**：无（使用固定初始供应量 1000000）
//
// **返回**：
//   - framework.SUCCESS (0) - 初始化成功
//   - framework.ERROR_INVALID_PARAMS (1) - ABI 版本不兼容
//   - framework.ERROR_EXECUTION_FAILED (6) - 执行失败
//
// **事件**：
//   - Initialized - 合约初始化事件
//     {
//     "owner": "<部署者地址>",
//     "initialSupply": "1000000"
//     }
//   - Mint - 代币铸造事件（由 SDK 自动发出）
//
// **状态变化**：
//   - 创建 AssetOutput（UTXO输出），向部署者分配初始代币
//
// **示例**：
//
//	合约部署时自动调用，无需手动调用
//
//export Initialize
func Initialize() uint32 {
	contract := &SimpleToken{}

	// 步骤1：检查 ABI 版本兼容性
	if err := framework.CheckABICompatibility(0x00010000); err != nil {
		contract.EmitLog("error", "ABI version mismatch")
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤2：获取部署者地址
	// 注意：framework.GetCaller() 返回 Address 类型，contract.GetCaller() 返回 string 类型
	owner := framework.GetCaller() // 返回 framework.Address 类型
	initialSupply := uint64(1000000)

	// 步骤3：使用 SDK 的 token.Mint() 向部署者铸造初始代币
	// ✅ token.Mint() 会自动处理：
	//   - 创建 AssetOutput（UTXO输出）
	//   - 交易构建
	//   - 事件发出
	// 注意：
	//   - TokenID 映射到 proto 的 fungible_class_id（代币类别标识符）
	//   - contract_address 由 CORE 自动设置（从 ExecutionContext 获取），SDK 无法篡改
	//   - 使用 "default" 作为代币标识符（如果合约只发行一种代币）
	tokenID := framework.TokenID("default")
	err := token.Mint(owner, tokenID, framework.Amount(initialSupply))
	if err != nil {
		contract.EmitLog("error", "Failed to mint initial supply: "+err.Error())
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤4：发出初始化事件
	eventData, _ := json.Marshal(map[string]string{
		"owner":         owner.String(), // Address 类型需要调用 String() 方法
		"initialSupply": strconv.FormatUint(initialSupply, 10),
	})
	contract.EmitEvent("Initialized", eventData)

	// 返回成功
	return framework.SUCCESS
}

// Mint 铸造代币
//
// 🎯 **用途**：向指定地址铸造代币
//
// **调用时机**：
//   - 任何用户都可以调用此函数
//   - 支持向指定地址铸造指定数量的代币
//
// **工作流程**：
//  1. 解析参数（to, amount）
//  2. 调用 SDK 的 token.Mint() 进行铸造
//     - SDK 内部自动创建 AssetOutput（UTXO输出）
//     - SDK 内部自动发出铸造事件
//
// **参数格式（JSON）**：
//
//	{
//	  "to": "接收者地址（Base58编码，可选，默认给调用者）",
//	  "amount": "100"  // 铸造数量（字符串格式，可选，默认100）
//	}
//
// **返回**：
//   - framework.SUCCESS (0) - 铸造成功
//   - framework.ERROR_INVALID_PARAMS (1) - 参数无效
//   - framework.ERROR_EXECUTION_FAILED (6) - 执行失败
//
// **事件**：
//   - Mint - 代币铸造事件（由 SDK 自动发出）
//     {
//     "to": "<接收者地址>",
//     "token_id": "",
//     "amount": 100,
//     "minter": "<调用者地址>"
//     }
//
// **状态变化**：
//   - 创建 AssetOutput（UTXO输出），可在任意节点查询余额
//
// **示例**：
//
//	调用 Mint()，参数为 {"to": "receiver_address", "amount": "100"}
//	会向 receiver_address 铸造 100 个代币
//
//export Mint
func Mint() uint32 {
	contract := &SimpleToken{}
	// 注意：framework.GetCaller() 返回 Address 类型，contract.GetCaller() 返回 string 类型
	caller := framework.GetCaller() // 获取调用者地址（Address 类型）

	// 步骤1：解析参数
	params := framework.GetContractParams()
	var toStr string
	var amountStr string
	if params != nil {
		toStr = params.ParseJSON("to")
		amountStr = params.ParseJSON("amount")
	}

	// 默认值：如果没有指定接收者，使用调用者地址
	if toStr == "" {
		toStr = caller.String() // Address 类型需要调用 String() 方法
	}

	// 默认值：如果没有指定数量，使用100
	amount := uint64(100)
	if amountStr != "" {
		var err error
		amount, err = strconv.ParseUint(amountStr, 10, 64)
		if err != nil || amount == 0 {
			contract.EmitLog("error", "Invalid amount")
			return framework.ERROR_INVALID_PARAMS
		}
	}

	// 步骤2：解析接收者地址
	to, err := framework.ParseAddressBase58(toStr)
	if err != nil {
		contract.EmitLog("error", "Invalid recipient address")
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤3：设置代币标识符
	// 注意：
	//   - TokenID 映射到 proto 的 fungible_class_id（代币类别标识符）
	//   - contract_address 由 CORE 自动设置（从 ExecutionContext 获取），SDK 无法篡改
	//   - 使用 "default" 作为代币标识符（如果合约只发行一种代币）
	tokenID := framework.TokenID("default")

	// 步骤4：使用 SDK 提供的铸造方法
	// ✅ token.Mint() 会自动处理：
	//   - 创建 AssetOutput（UTXO输出）
	//   - 交易构建
	//   - 事件发出
	err = token.Mint(to, tokenID, framework.Amount(amount))
	if err != nil {
		contract.EmitLog("error", "Mint failed: "+err.Error())
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 返回成功
	return framework.SUCCESS
}

// Transfer 转账
//
// 🎯 **用途**：从调用者向指定地址转账代币
//
// **调用时机**：
//   - 任何用户都可以调用此函数
//   - 调用者必须有足够的余额（通过UTXO查询）
//
// **工作流程**：
//  1. 解析参数（接收者地址、转账金额）
//  2. 验证参数有效性
//  3. 调用 SDK 的 token.Transfer() 进行转账
//     - SDK 内部自动检查余额（通过UTXO查询）
//     - SDK 内部自动构建交易
//     - SDK 内部自动发出转账事件
//
// **参数格式（JSON）**：
//
//	{
//	  "to": "receiver_address",  // 接收者地址（Base58编码，必填）
//	  "amount": "50"             // 转账金额（必填，字符串格式）
//	}
//
// **返回**：
//   - framework.SUCCESS (0) - 转账成功
//   - framework.ERROR_INVALID_PARAMS (1) - 参数无效
//   - framework.ERROR_INSUFFICIENT_BALANCE (4) - 余额不足
//   - framework.ERROR_EXECUTION_FAILED (6) - 执行失败
//
// **事件**：
//   - Transfer - 转账事件（由 SDK 自动发出）
//     {
//     "from": "<发送者地址>",
//     "to": "<接收者地址>",
//     "token_id": "",
//     "amount": 50
//     }
//
// **状态变化**：
//   - 消费发送者的UTXO，创建接收者的UTXO
//
// **示例**：
//
//	调用 Transfer()，参数为 {"to": "receiver_address", "amount": "50"}
//	如果调用者有足够余额，会从调用者扣除 50 个代币，向接收者增加 50 个代币
//
//export Transfer
func Transfer() uint32 {
	contract := &SimpleToken{}
	// 注意：framework.GetCaller() 返回 Address 类型，contract.GetCaller() 返回 string 类型
	caller := framework.GetCaller() // 获取调用者地址（Address 类型）

	// 步骤2：获取并解析参数
	params := framework.GetContractParams()
	var toStr string
	var amountStr string
	if params != nil {
		toStr = params.ParseJSON("to")         // 解析接收者地址
		amountStr = params.ParseJSON("amount") // 解析转账金额
	}

	// 步骤3：验证参数有效性
	if toStr == "" || amountStr == "" {
		contract.EmitLog("error", "Invalid parameters: to and amount are required")
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤4：解析金额
	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil || amount == 0 {
		contract.EmitLog("error", "Invalid amount")
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤5：解析接收者地址
	to, err := framework.ParseAddressBase58(toStr)
	if err != nil {
		contract.EmitLog("error", "Invalid recipient address")
		return framework.ERROR_INVALID_PARAMS
	}

	// 步骤6：设置代币标识符
	// 注意：
	//   - TokenID 映射到 proto 的 fungible_class_id（代币类别标识符）
	//   - contract_address 由 CORE 自动设置（从 ExecutionContext 获取），SDK 无法篡改
	//   - 使用 "default" 作为代币标识符（如果合约只发行一种代币）
	tokenID := framework.TokenID("default")

	// 步骤7：使用 SDK 提供的转账方法
	// ✅ token.Transfer() 会自动处理：
	//   - 余额检查（通过UTXO查询）
	//   - 交易构建
	//   - 事件发出
	// caller 是 framework.Address 类型，可以直接使用
	err = token.Transfer(caller, to, tokenID, framework.Amount(amount))
	if err != nil {
		contract.EmitLog("error", "Transfer failed: "+err.Error())
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 返回成功
	return framework.SUCCESS
}

// BalanceOf 查询余额
//
// 🎯 **用途**：查询指定地址的代币余额（只读函数）
//
// **调用时机**：
//   - 任何用户都可以调用此函数
//   - 这是一个只读函数，不会修改状态
//
// **工作流程**：
//  1. 解析参数（可选，默认查询调用者余额）
//  2. 使用 SDK 的 QueryUTXOBalance() 查询UTXO余额
//  3. 构造返回数据（JSON 格式）
//  4. 返回结果
//
// **参数格式（JSON）**（可选）：
//
//	{
//	  "address": "查询地址（Base58编码，可选，默认查询调用者）"
//	}
//
// **返回**：
//   - framework.SUCCESS (0) - 查询成功
//   - framework.ERROR_EXECUTION_FAILED (6) - 执行失败
//   - 返回数据（JSON 格式）：
//     {
//     "address": "<查询地址>",
//     "balance": 1000
//     }
//
// **状态变化**：无（只读函数）
//
// **示例**：
//
//	调用 BalanceOf()，返回 {"address": "<调用者地址>", "balance": 1000}
//
//export BalanceOf
func BalanceOf() uint32 {
	contract := &SimpleToken{}

	// 步骤1：解析参数（可选）
	params := framework.GetContractParams()
	addressStr := ""
	tokenIDStr := ""
	if params != nil {
		addressStr = params.ParseJSON("address")
		tokenIDStr = params.ParseJSON("token_id")
	}

	// 默认查询调用者余额
	// 注意：framework.GetCaller() 返回 Address 类型，contract.GetCaller() 返回 string 类型
	address := framework.GetCaller() // 获取调用者地址（Address 类型）
	if addressStr != "" {
		parsedAddr, err := framework.ParseAddressBase58(addressStr)
		if err == nil {
			address = parsedAddr // parsedAddr 是 Address 类型，可以直接赋值
		}
	}

	// 步骤2：设置代币标识符
	// 注意：
	//   - TokenID 映射到 proto 的 fungible_class_id（代币类别标识符）
	//   - contract_address 由 CORE 自动设置（从 ExecutionContext 获取），SDK 无法篡改
	//   - 默认使用 "default"，支持外部传入自定义 token_id
	tokenID := framework.TokenID("default")
	if tokenIDStr != "" {
		tokenID = framework.TokenID(tokenIDStr)
	}

	// 步骤3：使用 SDK 查询UTXO余额
	// ✅ QueryUTXOBalance() 从UTXO集合查询余额（返回最小单位）
	balance := framework.QueryUTXOBalance(address, tokenID)
	displayTokenID := "default"
	if tokenIDStr != "" {
		displayTokenID = tokenIDStr
	}
	result := framework.BuildBalanceResult(address.String(), displayTokenID, uint64(balance))

	// 将 map 序列化为 JSON 字符串
	resultJSON, err := json.Marshal(result)
	if err != nil {
		contract.EmitLog("error", "Failed to marshal result")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 步骤4：返回结果
	if err := contract.SetReturnData(resultJSON); err != nil {
		contract.EmitLog("error", "Failed to set return data")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 返回成功
	return framework.SUCCESS
}

// TotalSupply 查询总供应量
//
// 🎯 **用途**：查询代币的总供应量（只读函数）
//
// **调用时机**：
//   - 任何用户都可以调用此函数
//   - 这是一个只读函数，不会修改状态
//
// **说明**：
//   - 在UTXO模型中，总供应量可以通过查询所有UTXO的总和来计算
//   - 本示例简化实现：返回固定值（实际应用中可能需要遍历所有UTXO）
//   - 如果需要精确的总供应量，建议使用状态存储模式或事件索引
//
// **参数**：无
//
// **返回**：
//   - framework.SUCCESS (0) - 查询成功
//   - framework.ERROR_EXECUTION_FAILED (6) - 执行失败
//   - 返回数据（JSON 格式）：
//     {
//     "totalSupply": 1000000,
//     "note": "在UTXO模型中，总供应量等于所有UTXO的总和"
//     }
//
// **状态变化**：无（只读函数）
//
// **示例**：
//
//	调用 TotalSupply() 后，返回总供应量信息
//
//export TotalSupply
func TotalSupply() uint32 {
	contract := &SimpleToken{}

	// 注意：在UTXO模型中，总供应量应该通过查询所有UTXO的总和来计算
	// 本示例简化实现：返回固定值
	// 实际应用中，可以通过以下方式获取：
	// 1. 使用状态存储记录总供应量（StateOutput）
	// 2. 遍历所有UTXO并求和（需要节点支持）
	// 3. 通过事件索引统计所有Mint事件的总和

	totalSupply := uint64(1000000) // 简化实现：固定值

	// 构造返回数据（JSON 格式）
	result := map[string]interface{}{
		"totalSupply": totalSupply,
		"note":        "在UTXO模型中，总供应量等于所有UTXO的总和",
	}

	// 将 map 序列化为 JSON 字符串
	resultJSON, err := json.Marshal(result)
	if err != nil {
		contract.EmitLog("error", "Failed to marshal result")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 返回结果
	if err := contract.SetReturnData(resultJSON); err != nil {
		contract.EmitLog("error", "Failed to set return data")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 返回成功
	return framework.SUCCESS
}

// main 函数（TinyGo 编译 WASM 时需要的入口函数）
//
// 🎯 **用途**：TinyGo 编译 WASM 时需要的入口函数
//
// **说明**：
//   - WASM 合约必须有 main 函数，但实际运行时不会被调用
//   - 合约的入口是使用 `//export` 标记的函数（如 Initialize、Transfer 等）
func main() {}
