//go:build tinygo || (js && wasm)

// Package main 提供最简单的 WES 智能合约示例 - Hello World
//
// 📋 示例说明
//
// 本示例是 WES Contract SDK Go 的入门示例，展示如何构建一个最简单的智能合约。
// 通过本示例，您可以学习：
//   - 如何创建一个基本的智能合约
//   - 如何使用 Framework 层的基础功能
//   - 如何初始化合约并保存状态
//   - 如何发出事件和日志
//   - 如何查询合约状态
//
// 🎯 核心功能
//
//  1. Initialize - 初始化合约
//     - 在合约部署时自动调用
//     - 记录部署者地址和部署时间
//     - 初始化问候计数器
//
//  2. SayHello - 问候功能
//     - 向调用者发送问候消息
//     - 记录并增加问候计数
//     - 发出问候事件
//
//  3. GetGreetingCount - 查询问候次数
//     - 只读函数，查询总问候次数
//     - 返回字符串格式的结果
//
//  4. GetDeployerInfo - 查询部署者信息
//     - 只读函数，查询部署者地址和部署时间
//     - 返回字符串格式的结果
//
// 📚 学习要点
//
//  - **合约结构**：所有 WES 合约都嵌入 `framework.ContractBase`
//  - **导出函数**：使用 `//export` 注释标记可被外部调用的函数
//  - **状态存储**：使用 `appendStateKV()` 保存状态（底层使用 AppendStateOutput）
//  - **状态查询**：使用 `GetState()` 读取状态
//  - **事件发出**：使用 `EmitEvent()` 发出链上事件
//  - **日志记录**：使用 `EmitLog()` 记录调试日志
//
// ⚠️ 注意事项
//
//  - 本示例直接使用底层接口（`framework.AppendStateOutputSimple`），仅用于演示目的
//  - ✅ **推荐**：在实际开发中，应使用 `helpers` 层的业务语义接口
//  - 例如：使用 `token.Transfer()` 而不是直接操作状态输出
//
// 📚 相关文档
//
//   - [Framework 文档](../../framework/README.md) - Framework 层详细说明
//   - [示例总览](../README.md) - 所有示例索引
//   - [Hello World 示例 README](./README.md) - 本示例详细文档
package main

import (
	framework "github.com/weisyn/contract-sdk-go/framework"
)

// HelloContract 最简单的 WES 智能合约示例
//
// 本合约展示了 WES 智能合约的基本结构和使用方式。
// 所有 WES 合约都需要嵌入 `framework.ContractBase`，以获得基础功能。
//
// 设计理念：
//   - 简单易懂：使用最简单的功能展示合约的基本结构
//   - 完整示例：包含初始化、状态读写、事件发出等基本操作
//   - 学习友好：适合初学者理解 WES 合约的基本概念
type HelloContract struct {
	framework.ContractBase
}

// Initialize 合约初始化函数
//
// 🎯 **用途**：在合约部署时自动调用，用于初始化合约状态
//
// **调用时机**：
//   - 合约部署时自动调用一次
//   - 只会被调用一次，用于设置初始状态
//
// **工作流程**：
//   1. 获取部署者地址（调用者地址）
//   2. 获取当前时间戳
//   3. 保存部署信息到状态：
//      - deployer: 部署者地址
//      - deployed_at: 部署时间戳
//      - greeting_count: 问候计数器（初始化为 0）
//   4. 发出部署事件
//   5. 记录日志
//
// **参数**：无
//
// **返回**：
//   - 0 (SUCCESS) - 初始化成功
//
// **事件**：
//   - ContractDeployed - 合约部署事件
//     {
//       "deployer": "<部署者地址>"
//     }
//
// **状态变化**：
//   - 创建状态键 "deployer"，值为部署者地址
//   - 创建状态键 "deployed_at"，值为部署时间戳
//   - 创建状态键 "greeting_count"，值为 0
//
// **示例**：
//   合约部署时自动调用，无需手动调用
//
//export Initialize
func Initialize() uint32 {
	contract := &HelloContract{}

	// 步骤1：获取部署者地址
	// GetCaller() 返回调用当前函数的地址（在部署时就是部署者地址）
	deployer := contract.GetCaller()

	// 步骤2：获取当前时间戳
	// GetTimestamp() 返回当前区块的时间戳（Unix 时间戳，秒）
	timestamp := contract.GetTimestamp()

	// 步骤3：保存部署信息到状态
	// appendStateKV() 是一个辅助函数，用于保存键值对到状态
	// 底层使用 AppendStateOutput 机制，状态会被记录到交易输出中
	appendStateKV("deployer", []byte(deployer))
	appendStateKV("deployed_at", uint64ToBytes(timestamp))
	appendStateKV("greeting_count", uint64ToBytes(0))

	// 步骤4：发出部署事件
	// EmitEvent() 发出链上事件，可以被外部系统监听
	// 事件会被记录到区块链上，用于追踪合约状态变化
	contract.EmitEvent("ContractDeployed", []byte(deployer))

	// 步骤5：记录日志
	// EmitLog() 记录调试日志，用于开发和调试
	// 日志不会上链，仅用于调试目的
	contract.EmitLog("INFO", "Hello World 合约已部署,部署者: "+deployer)

	// 返回成功
	return 0 // SUCCESS
}

// SayHello 问候函数
//
// 🎯 **用途**：向调用者发送问候消息，并记录访问次数
//
// **调用时机**：
//   - 任何用户都可以调用此函数
//   - 每次调用都会增加问候计数
//
// **工作流程**：
//   1. 获取调用者地址
//   2. 获取当前区块高度
//   3. 读取当前问候计数
//   4. 增加问候计数并保存
//   5. 构造问候消息（包含调用者、问候次数、区块高度）
//   6. 发出问候事件
//   7. 记录日志
//   8. 返回问候消息
//
// **参数**：无
//
// **返回**：
//   - 0 (SUCCESS) - 执行成功
//   - 返回数据：问候消息（字符串格式）
//
// **事件**：
//   - Greeting - 问候事件
//     {
//       "message": "Hello, <调用者地址>! This is greeting #<次数> at block <区块高度>"
//     }
//
// **状态变化**：
//   - 更新状态键 "greeting_count"，值增加 1
//
// **示例**：
//   用户调用 SayHello() 后，会收到问候消息，问候计数增加
//
//export SayHello
func SayHello() uint32 {
	contract := &HelloContract{}

	// 步骤1：获取调用者信息
	// GetCaller() 返回调用当前函数的地址
	caller := contract.GetCaller()

	// 步骤2：获取当前区块高度
	// GetBlockHeight() 返回当前区块的高度（区块编号）
	blockHeight := contract.GetBlockHeight()

	// 步骤3：读取当前问候计数
	// GetState() 读取状态值，返回字节数组
	countData := contract.GetState("greeting_count")
	// bytesToUint64() 将字节数组转换为 uint64 数字
	count := bytesToUint64(countData)

	// 步骤4：增加问候计数
	count++
	// 保存更新后的计数到状态
	appendStateKV("greeting_count", uint64ToBytes(count))

	// 步骤5：构造问候消息
	// 消息包含：调用者地址、问候次数、区块高度
	message := "Hello, " + caller + "! This is greeting #" + uint64ToString(count) +
		" at block " + uint64ToString(blockHeight)

	// 步骤6：发出问候事件
	// 事件会被记录到区块链上，外部系统可以监听此事件
	contract.EmitEvent("Greeting", []byte(message))

	// 步骤7：记录日志
	// 日志用于调试，不会上链
	contract.EmitLog("INFO", "问候: "+message)

	// 步骤8：返回问候消息
	// SetReturnData() 设置函数返回值，外部调用者可以获取此值
	contract.SetReturnData([]byte(message))

	// 返回成功
	return 0 // SUCCESS
}

// GetGreetingCount 查询问候次数
//
// 🎯 **用途**：查询总问候次数（只读函数）
//
// **调用时机**：
//   - 任何用户都可以调用此函数
//   - 这是一个只读函数，不会修改状态
//
// **工作流程**：
//   1. 读取问候计数状态
//   2. 构造返回消息
//   3. 返回结果
//
// **参数**：无
//
// **返回**：
//   - 0 (SUCCESS) - 查询成功
//   - 返回数据：总问候次数（字符串格式，如 "Total greetings: 5"）
//
// **状态变化**：无（只读函数）
//
// **示例**：
//   调用 GetGreetingCount() 后，返回 "Total greetings: 10"
//
//export GetGreetingCount
func GetGreetingCount() uint32 {
	contract := &HelloContract{}

	// 步骤1：读取问候计数
	// GetState() 读取状态值，返回字节数组
	countData := contract.GetState("greeting_count")
	// bytesToUint64() 将字节数组转换为 uint64 数字
	count := bytesToUint64(countData)

	// 步骤2：构造返回消息
	// 将数字转换为字符串并格式化
	result := "Total greetings: " + uint64ToString(count)

	// 步骤3：返回结果
	// SetReturnData() 设置函数返回值
	contract.SetReturnData([]byte(result))

	// 返回成功
	return 0 // SUCCESS
}

// GetDeployerInfo 获取部署者信息
//
// 🎯 **用途**：查询合约部署者地址和部署时间（只读函数）
//
// **调用时机**：
//   - 任何用户都可以调用此函数
//   - 这是一个只读函数，不会修改状态
//
// **工作流程**：
//   1. 读取部署者地址
//   2. 读取部署时间戳
//   3. 构造返回信息
//   4. 返回结果
//
// **参数**：无
//
// **返回**：
//   - 0 (SUCCESS) - 查询成功
//   - 返回数据：部署者信息（字符串格式，如 "Deployer: <地址>, Deployed at: <时间戳>"）
//
// **状态变化**：无（只读函数）
//
// **示例**：
//   调用 GetDeployerInfo() 后，返回 "Deployer: <地址>, Deployed at: 1640995200"
//
//export GetDeployerInfo
func GetDeployerInfo() uint32 {
	contract := &HelloContract{}

	// 步骤1：读取部署信息
	// GetState() 读取状态值，返回字节数组
	deployer := string(contract.GetState("deployer"))
	deployedAtData := contract.GetState("deployed_at")
	deployedAt := bytesToUint64(deployedAtData)

	// 步骤2：构造返回信息
	// 将部署者地址和时间戳格式化为字符串
	info := "Deployer: " + deployer + ", Deployed at: " + uint64ToString(deployedAt)

	// 步骤3：返回结果
	// SetReturnData() 设置函数返回值
	contract.SetReturnData([]byte(info))

	// 返回成功
	return 0 // SUCCESS
}

// ============================================================================
// 辅助函数
// ============================================================================

// uint64ToBytes 将 uint64 转换为字节数组（小端序）
//
// 🎯 **用途**：将数字转换为字节数组，用于状态存储
//
// **参数**：
//   - n: 要转换的 uint64 数字
//
// **返回**：
//   - []byte: 8 字节的字节数组（小端序）
//
// **说明**：
//   - 小端序：最低有效字节在前
//   - 例如：0x1234567890ABCDEF 转换为 [0xEF, 0xCD, 0xAB, 0x90, 0x78, 0x56, 0x34, 0x12]
//
// **示例**：
//   bytes := uint64ToBytes(1000)
//   // 结果: [0xE8, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00]
func uint64ToBytes(n uint64) []byte {
	return []byte{
		byte(n),         // 最低有效字节（第 0 字节）
		byte(n >> 8),    // 第 1 字节
		byte(n >> 16),   // 第 2 字节
		byte(n >> 24),   // 第 3 字节
		byte(n >> 32),   // 第 4 字节
		byte(n >> 40),   // 第 5 字节
		byte(n >> 48),   // 第 6 字节
		byte(n >> 56),   // 最高有效字节（第 7 字节）
	}
}

// bytesToUint64 将字节数组转换为 uint64（小端序）
//
// 🎯 **用途**：将字节数组转换为数字，用于状态读取
//
// **参数**：
//   - b: 字节数组（至少 1 字节，不足 8 字节会自动填充）
//
// **返回**：
//   - uint64: 转换后的数字
//
// **说明**：
//   - 如果字节数组为空，返回 0
//   - 如果字节数组长度小于 8，会自动填充到 8 字节
//   - 小端序：最低有效字节在前
//
// **示例**：
//   n := bytesToUint64([]byte{0xE8, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
//   // 结果: 1000
func bytesToUint64(b []byte) uint64 {
	// 如果字节数组为空，返回 0
	if len(b) == 0 {
		return 0
	}

	// 如果字节数组长度小于 8，填充到 8 字节
	if len(b) < 8 {
		padded := make([]byte, 8)
		copy(padded, b)
		b = padded
	}

	// 小端序解码：最低有效字节在前
	return uint64(b[0]) |
		uint64(b[1])<<8 |
		uint64(b[2])<<16 |
		uint64(b[3])<<24 |
		uint64(b[4])<<32 |
		uint64(b[5])<<40 |
		uint64(b[6])<<48 |
		uint64(b[7])<<56
}

// uint64ToString 将 uint64 转换为字符串
//
// 🎯 **用途**：将数字转换为字符串，用于消息格式化
//
// **参数**：
//   - n: 要转换的 uint64 数字
//
// **返回**：
//   - string: 数字的字符串表示
//
// **说明**：
//   - 简化实现，仅用于示例
//   - 实际项目中可以使用标准库的 strconv.FormatUint()
//
// **示例**：
//   s := uint64ToString(1000)
//   // 结果: "1000"
func uint64ToString(n uint64) string {
	// 特殊情况：0
	if n == 0 {
		return "0"
	}

	// 从低位到高位提取每一位数字
	digits := make([]byte, 0, 20) // 最多 20 位数字
	for n > 0 {
		digits = append(digits, byte('0'+n%10)) // 取最低位数字
		n /= 10                                  // 右移一位
	}

	// 反转数字数组（因为是从低位到高位提取的）
	for i := 0; i < len(digits)/2; i++ {
		digits[i], digits[len(digits)-1-i] = digits[len(digits)-1-i], digits[i]
	}

	return string(digits)
}

// main 函数（WASM 合约必须有 main，但不会被调用）
//
// 🎯 **用途**：TinyGo 编译 WASM 时需要的入口函数
//
// **说明**：
//   - WASM 合约必须有 main 函数，但实际运行时不会被调用
//   - 合约的入口是使用 `//export` 标记的函数（如 Initialize、SayHello 等）
//
func main() {}

// =========================================================================
// 状态锚定辅助函数（AppendStateOutput 封装）
// =========================================================================

// appendStateKV 将 key/value 以显式状态输出方式锚定
//
// 🎯 **用途**：保存键值对到合约状态
//
// **参数**：
//   - key: 状态键（字符串）
//   - value: 状态值（字节数组）
//
// **说明**：
//   - 使用 AppendStateOutputSimple（推荐）：合约仅提供业务锚定信息
//   - ZK 证明由系统层在交易固化时自动生成并附加
//   - 保证业务逻辑与证明生成完全解耦
//
// **约定**：
//   - stateID = key 的字节
//   - execution_result_hash = value 原文（演示用，生产可改为哈希）
//   - version = 1（演示场景固定版本）
//   - parentHash = nil（演示场景不构建状态链）
//
// ⚠️ **注意**：
//   - 本函数直接使用底层接口（`framework.AppendStateOutputSimple`），仅用于演示目的
//   - ✅ **推荐**：在实际开发中，应使用 `helpers` 层的业务语义接口
//   - 例如：使用 `token.Transfer()` 而不是直接操作状态输出
//
// **示例**：
//   appendStateKV("greeting_count", uint64ToBytes(10))
//   // 保存键 "greeting_count"，值为 10 的字节表示
func appendStateKV(key string, value []byte) {
	// 参数验证：如果 key 为空，直接返回
	if len(key) == 0 {
		return
	}

	// 构建状态输出参数
	stateID := []byte(key)        // 状态ID = key 的字节
	version := uint64(1)          // 版本号（固定为 1）
	execHash := value             // 执行结果哈希 = value 原文（演示用）
	var parentHash []byte         // 父状态哈希（nil，不构建状态链）

	// ⚠️ 不推荐：直接使用底层接口（仅用于演示目的）
	// ✅ 推荐：使用 helpers 层的业务语义接口
	// 使用 framework 包的状态输出接口
	_, _ = framework.AppendStateOutputSimple(stateID, version, execHash, parentHash)
}

// =========================================================================
// 演示：批量输出（v1.1 新增）
// =========================================================================

// DemoBatchOutputs 演示批量创建资产输出
//
// 🎯 **用途**：演示如何批量创建多个资产输出
//
// **适用场景**：
//   - 空投：向多个地址批量发送代币
//   - 批量转账：一次性向多个地址转账
//   - 批量分配：将资产分配给多个受益人
//
// **说明**：
//   - 本函数仅用于演示，实际项目中应该作为合约函数导出
//   - 批量操作可以提高效率，减少交易次数
//
// ⚠️ **注意**：
//   - 本函数直接使用底层接口（`framework.BatchCreateOutputsSimple`），仅用于演示目的
//   - ✅ **推荐**：在实际开发中，应使用 `helpers` 层的业务语义接口
//   - 例如：使用 `token.Airdrop()` 而不是直接批量创建输出
//
// **示例**：
//   DemoBatchOutputs()
//   // 批量创建 3 个输出，分别向 3 个地址发送 100、110、120 个代币
func DemoBatchOutputs() {
	// 步骤1：准备接收者地址列表
	// 示例：3 个接收者地址
	recipients := [][]byte{
		[]byte("recipient_address_1_"),
		[]byte("recipient_address_2_"),
		[]byte("recipient_address_3_"),
	}

	// 步骤2：构建批量输出项
	// 每个输出项包含：接收者地址、金额、代币ID
	items := make([]struct {
		Recipient []byte // 接收者地址
		Amount    uint64 // 金额
		TokenID   []byte // 代币ID（nil 表示原生币）
	}, len(recipients))

	// 步骤3：填充输出项
	for i, addr := range recipients {
		items[i].Recipient = addr
		items[i].Amount = uint64(100 + i*10) // 金额：100, 110, 120
		items[i].TokenID = nil               // 原生币（nil 表示原生币）
	}

	// 步骤4：批量创建输出
	// BatchCreateOutputsSimple() 批量创建多个资产输出
	count, err := framework.BatchCreateOutputsSimple(items)
	if err != nil {
		// 错误处理：记录错误并返回
		framework.SetReturnData([]byte("批量输出失败: " + err.Error()))
		return
	}

	// 步骤5：记录结果
	// 将批量输出的数量保存到状态
	appendStateKV("batch_output_count", []byte(uint64ToString(uint64(count))))
}
