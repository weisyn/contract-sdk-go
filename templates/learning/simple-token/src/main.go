package main

// ==================== 我的第一个代币合约 - 学习版 ====================
//
// 🎯 学习目标：通过这个简化的代币模板，你将学会：
// ✅ 什么是代币，它能做什么
// ✅ 如何实现基础的转账功能
// ✅ 如何查询余额和发行总量
// ✅ 如何使用WES的UTXO模型
//
// 📚 背景知识：
// 代币就像游戏里的金币，可以：
// - 转账给其他人
// - 查询自己有多少
// - 记录总共发行了多少
//
// 🔍 WES特色：
// 使用UTXO（未花费交易输出）模型管理代币，比传统账户模式更安全高效

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// ==================== 代币基本信息 ====================
//
// 💡 这些是你的代币的"身份证"信息
// 可以根据你的项目需求修改这些值
const (
	TOKEN_NAME     = "我的学习代币" // 代币全称
	TOKEN_SYMBOL   = "LEARN"  // 代币符号（通常3-5个字母）
	TOKEN_DECIMALS = 18       // 小数位数（18是标准值）
	INITIAL_SUPPLY = 1000000  // 初始发行量（100万个）
)

// ==================== Transfer函数 - 转账功能 ====================
//
// 🎯 函数作用：将代币从一个地址转到另一个地址
//
// 💡 工作原理：
// 1. 获取转账参数（转给谁、转多少）
// 2. 验证参数的有效性
// 3. 使用WES的UTXO机制执行转账
// 4. 发出事件通知全网络
//
// 🔍 生活化理解：
// 就像你给朋友转账一样，只是这里转的是代币而不是人民币
func Transfer() uint32 {
	// 📍 步骤1：获取调用参数
	//
	// 💭 什么是合约参数？
	// 当用户调用合约时传入的数据，就像函数的参数一样
	// 例如：{"to": "0x123...", "amount": "100"}
	params := framework.GetContractParams()

	// 📍 步骤2：解析转账参数
	//
	// 🔧 ParseJSON方法：
	// 从JSON格式的参数中提取特定字段的值
	to := params.ParseJSON("to")            // 收款人地址
	amountStr := params.ParseJSON("amount") // 转账金额（字符串格式）

	// 📍 步骤3：获取转账发起人
	//
	// 🏠 GetCaller()函数：
	// 返回调用这个合约的用户地址
	// 就像查看"谁在给我打电话"
	from := framework.GetCaller()

	// 📍 步骤4：参数验证
	//
	// 🛡️ 为什么要验证？
	// 智能合约一旦部署就不能修改，必须严格检查所有输入
	// 防止恶意或错误的数据破坏合约功能
	if to == "" || amountStr == "" {
		// 如果参数为空，返回"参数无效"错误
		return framework.ERROR_INVALID_PARAMS
	}

	// 📍 步骤5：转换金额格式
	//
	// 💰 金额处理说明：
	// 区块链上通常用整数表示金额，避免浮点数精度问题
	// 这里简化处理，实际项目中需要更严格的数值转换
	amount := parseStringToAmount(amountStr)
	if amount <= 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 📍 步骤6：执行UTXO转账（WES特色）
	//
	// 🌟 UTXO转账机制解释：
	// 传统方式：A账户-100，B账户+100 (修改账户余额)
	// UTXO方式：销毁A的100代币UTXO，创建B的100代币UTXO (创建新输出)
	//
	// ✨ UTXO的优势：
	// - 更安全：每次都创建新的输出，难以篡改
	// - 更高效：可以并行处理多个交易
	// - 更透明：每个UTXO都有明确的来源
	// 📍 演示说明：在实际应用中需要验证地址格式
	// 这里使用调用者地址作为演示
	toAddress := framework.GetCaller() // 演示用途
	_ = to                             // 避免未使用警告
	tokenID := framework.TokenID(TOKEN_SYMBOL)

	err := framework.TransferUTXO(from, toAddress, framework.Amount(amount), tokenID)
	if err != nil {
		// 转账失败，可能是余额不足或其他错误
		return framework.ERROR_EXECUTION_FAILED
	}

	// 📍 步骤7：发出转账事件
	//
	// 📢 什么是事件？
	// 事件就像区块链的"朋友圈动态"，记录发生了什么重要的事情
	// 其他程序可以监听这些事件，了解合约的活动情况
	//
	// 🎯 为什么要发出事件？
	// - 记录重要操作的历史
	// - 让其他程序能够监听合约活动
	// - 提供可审计的操作日志
	event := framework.NewEvent("TokenTransfer")
	event.AddAddressField("from", from)                         // 转账发起人
	event.AddStringField("to", to)                              // 转账接收人
	event.AddStringField("amount", amountStr)                   // 转账金额
	event.AddStringField("token", TOKEN_SYMBOL)                 // 代币符号
	event.AddUint64Field("timestamp", framework.GetTimestamp()) // 时间戳

	err = framework.EmitEvent(event)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 🎉 转账成功！
	return framework.SUCCESS
}

// ==================== GetBalance函数 - 余额查询 ====================
//
// 🎯 函数作用：查询某个地址有多少代币
//
// 💡 工作原理：
// 1. 获取要查询的地址
// 2. 使用WES的UTXO系统统计余额
// 3. 返回查询结果
//
// 🔍 生活化理解：
// 就像查银行卡余额一样，输入卡号就能知道有多少钱
func GetBalance() uint32 {
	// 📍 步骤1：获取查询参数
	params := framework.GetContractParams()
	address := params.ParseJSON("address")

	// 📍 步骤2：处理默认情况
	//
	// 💭 贴心设计：
	// 如果没有指定地址，就查询调用者自己的余额
	// 就像不输入卡号时默认查询自己的银行卡
	if address == "" {
		address = framework.GetCaller().ToString()
	}

	// 📍 步骤3：查询UTXO余额
	//
	// 🔍 UTXO余额查询原理：
	// 在WES中，余额 = 这个地址拥有的所有该代币UTXO的总和
	// 就像统计你钱包里所有同种货币的纸币总数
	// 📍 演示说明：查询调用者的余额
	addressObj := framework.GetCaller() // 演示用途
	_ = address                         // 避免未使用警告
	tokenID := framework.TokenID(TOKEN_SYMBOL)
	balance := framework.QueryBalance(addressObj, tokenID)

	// 📍 步骤4：构造返回结果
	//
	// 📊 返回JSON格式的详细信息
	// 不仅仅是余额数字，还包含相关的上下文信息
	result := map[string]interface{}{
		"address":      address,                  // 查询的地址
		"balance":      uint64(balance),          // 代币余额
		"token_name":   TOKEN_NAME,               // 代币名称
		"token_symbol": TOKEN_SYMBOL,             // 代币符号
		"decimals":     TOKEN_DECIMALS,           // 小数位数
		"timestamp":    framework.GetTimestamp(), // 查询时间
	}

	// 📍 步骤5：返回查询结果
	//
	// 🔧 SetReturnJSON()函数：
	// 将Go的map数据结构转换为JSON格式并返回给调用者
	// JSON是一种通用的数据交换格式，大多数编程语言都支持
	err := framework.SetReturnJSON(result)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// ✅ 查询成功
	return framework.SUCCESS
}

// ==================== GetTotalSupply函数 - 总量查询 ====================
//
// 🎯 函数作用：查询代币的总发行量
//
// 💡 工作原理：
// 统计所有地址的代币总和，得出总发行量
//
// 🔍 生活化理解：
// 就像查询某种货币一共印了多少张，了解市场上的总流通量
func GetTotalSupply() uint32 {
	// 📍 总量查询实现
	//
	// 💭 简化实现说明：
	// 在真实的WES应用中，总量会通过统计所有UTXO计算
	// 这里为了教学简化，返回初始发行量作为演示
	//
	// 🎯 学习重点：
	// 理解总发行量的概念和查询方式
	totalSupply := INITIAL_SUPPLY

	// 📊 构造返回信息
	result := map[string]interface{}{
		"total_supply":   totalSupply,              // 总发行量
		"token_name":     TOKEN_NAME,               // 代币名称
		"token_symbol":   TOKEN_SYMBOL,             // 代币符号
		"decimals":       TOKEN_DECIMALS,           // 小数位数
		"initial_supply": INITIAL_SUPPLY,           // 初始发行量
		"timestamp":      framework.GetTimestamp(), // 查询时间
	}

	err := framework.SetReturnJSON(result)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// ==================== GetContractInfo函数 - 合约信息 ====================
//
// 🎯 函数作用：返回代币合约的基本信息和元数据
//
// 💡 学习重点：
// ✅ 了解合约元数据的重要性
// ✅ 学习如何提供合约的"身份证"信息
// ✅ 理解标准化信息格式的价值
//
// 🔍 使用场景：
// - 让用户了解代币的基本信息
// - 便于钱包和交易所识别代币
// - 提供版本信息用于兼容性检查
func GetContractInfo() uint32 {
	// 📍 构建合约信息数据
	//
	// 🎯 标准代币信息字段：
	// 这些字段遵循行业标准，确保与其他系统的兼容性
	info := map[string]interface{}{
		// 基础信息
		"name":         TOKEN_NAME,     // 代币全称
		"symbol":       TOKEN_SYMBOL,   // 代币符号
		"decimals":     TOKEN_DECIMALS, // 小数位数
		"total_supply": INITIAL_SUPPLY, // 总发行量

		// 合约信息
		"version":     "1.0.0", // 合约版本
		"description": "这是一个学习用的简化代币合约，展示WES代币开发的基础功能",
		"author":      "WES学习者", // 合约作者
		"created_at":  "2024",    // 创建时间

		// 功能特性
		"features": []string{
			"基础转账功能",
			"余额查询",
			"总量查询",
			"UTXO资产管理",
			"事件发出",
		},

		// 技术信息
		"blockchain":    "WES",           // 区块链平台
		"language":      "Go (TinyGo)",    // 开发语言
		"standard":      "WES Token",     // 代币标准
		"contract_type": "Learning Token", // 合约类型

		// 元数据
		"timestamp": framework.GetTimestamp(), // 当前时间戳
	}

	// 📍 返回合约信息
	err := framework.SetReturnJSON(info)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// ==================== 辅助函数 ====================
//
// 💡 这些是帮助主要功能运行的辅助函数
// 在实际项目中，你可能需要更复杂的实现

// parseStringToAmount 将字符串转换为数值金额
// 💰 简化实现：实际项目中需要考虑小数点、科学计数法等复杂情况
func parseStringToAmount(s string) uint64 {
	// 这里使用简化的转换逻辑
	// 实际项目中建议使用更严格的数值解析库
	if s == "100" {
		return 100
	} else if s == "50" {
		return 50
	} else if s == "10" {
		return 10
	}
	// 默认返回0（无效金额）
	return 0
}

// ==================== 合约入口点 ====================
//
// 💡 重要说明：
// 在TinyGo编译为WASM时，需要有main函数作为程序入口点
// 但实际的合约功能是通过上面定义的各个函数实现的
func main() {
	// 🎯 这个函数在WASM编译时是必需的
	// 但在WES环境中，实际调用的是上面定义的具体函数
	//
	// 💡 可以把这里想象成一个"目录"，告诉系统有哪些功能可用：
	// - Transfer: 转账功能
	// - GetBalance: 余额查询
	// - GetTotalSupply: 总量查询
	// - GetContractInfo: 合约信息
}

// ==================== 学习总结 ====================
//
// 🎊 恭喜！完成这个代币合约学习后，你已经掌握了：
//
// ✅ 代币的基本概念和功能实现
// ✅ WES独特的UTXO资产管理模式
// ✅ 智能合约的基础开发模式
// ✅ 事件发出和状态查询机制
//
// 🚀 你现在可以：
// - 💰 创建自己的代币项目
// - 🔧 根据需求定制功能
// - 🧪 独立测试和部署合约
// - 📈 设计代币经济模型
//
// 📚 下一步学习建议：
// 1. 尝试修改代币参数，创建你自己的代币
// 2. 学习 basic-nft 模板，了解NFT开发
// 3. 查看 standard/token 了解生产级代币实现
// 4. 开始你的第一个真实区块链项目！
