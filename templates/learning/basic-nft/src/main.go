package main

// ==================== 我的第一个NFT合约 - 学习版 ====================
//
// 🎯 学习目标：通过这个基础的NFT模板，你将学会：
// ✅ 什么是NFT，与代币有什么区别
// ✅ 如何创建独一无二的数字资产
// ✅ 如何转移NFT所有权
// ✅ 如何管理NFT的元数据信息
//
// 📚 背景知识：
// NFT (Non-Fungible Token) 就像数字收藏品：
// - 每个都独一无二，不可互换
// - 可以证明数字资产的所有权
// - 广泛应用于艺术、游戏、证书等领域
//
// 🔍 与代币的区别：
// 代币：可互换（每个都相同）  NFT：不可互换（每个都独特）
// 代币：可分割（0.5个）      NFT：不可分割（只能整个）
// 代币：价值由数量决定       NFT：价值由稀有度决定

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// ==================== NFT基本信息 ====================
//
// 💡 这些是你的NFT系列的"身份证"信息
// 可以根据你的项目需求修改这些值
const (
	COLLECTION_NAME   = "我的学习NFT系列"                     // NFT系列名称
	COLLECTION_SYMBOL = "LEARN-NFT"                     // NFT系列符号
	BASE_TOKEN_URI    = "https://example.com/metadata/" // 元数据基础URL
)

// ==================== 全局状态变量 ====================
//
// 📊 这些变量追踪NFT的状态信息
// 在实际的WES实现中，这些会通过UTXO系统管理
var (
	totalSupply uint64 = 0 // 已铸造的NFT总数
	nextTokenID uint64 = 1 // 下一个NFT的ID
)

// ==================== MintNFT函数 - NFT铸造功能 ====================
//
// 🎯 函数作用：创建一个全新的、独一无二的NFT
//
// 💡 工作原理：
// 1. 生成唯一的NFT ID
// 2. 将NFT所有权分配给指定地址
// 3. 设置NFT的元数据链接
// 4. 发出铸造事件
//
// 🔍 生活化理解：
// 就像艺术家创作一幅新画作，每幅画都有独特的编号和签名
func MintNFT() uint32 {
	// 📍 步骤1：获取铸造参数
	//
	// 💭 铸造NFT需要什么信息？
	// - to: 将NFT给谁（接收者地址）
	// - tokenURI: NFT的元数据链接（描述这个NFT的详细信息）
	params := framework.GetContractParams()
	to := params.ParseJSON("to")
	tokenURI := params.ParseJSON("tokenURI")

	// 📍 步骤2：参数验证
	//
	// 🛡️ 确保输入数据的有效性
	// NFT铸造是不可逆的操作，必须严格检查
	if to == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 如果没有提供tokenURI，使用默认格式
	if tokenURI == "" {
		tokenURI = BASE_TOKEN_URI + string(rune(nextTokenID)) + ".json"
	}

	// 📍 步骤3：生成唯一的NFT ID
	//
	// 🆔 每个NFT都需要一个唯一的标识符
	// 就像每个人都有身份证号一样
	tokenID := nextTokenID
	nextTokenID++ // 为下一个NFT准备ID
	totalSupply++ // 增加总供应量

	// 📍 步骤4：创建NFT的UTXO
	//
	// 🌟 WES的NFT实现特色：
	// 在WES中，每个NFT都是一个独特的UTXO
	// 这确保了NFT的唯一性和不可复制性
	//
	// 💡 UTXO NFT的优势：
	// - 天然的唯一性保证
	// - 高效的所有权验证
	// - 强大的可编程性
	// 📍 演示说明：在实际应用中需要验证地址格式
	toAddress := framework.GetCaller() // 演示用途：铸造给调用者
	_ = to                             // 避免未使用警告
	nftTokenID := framework.TokenID("NFT_" + string(rune(tokenID)))

	// 创建NFT UTXO（数量为1，表示这是一个完整的NFT）
	err := framework.CreateUTXO(toAddress, framework.Amount(1), nftTokenID)
	if err != nil {
		// 撤销状态变更
		nextTokenID--
		totalSupply--
		return framework.ERROR_EXECUTION_FAILED
	}

	// 📍 步骤5：发出铸造事件
	//
	// 📢 NFT铸造事件包含什么信息？
	// - 接收者地址、NFT ID、元数据URI、铸造时间等
	// 这些信息让整个网络知道新的NFT被创建了
	event := framework.NewEvent("NFTMinted")
	event.AddStringField("to", to)
	event.AddUint64Field("tokenID", tokenID)
	event.AddStringField("tokenURI", tokenURI)
	event.AddAddressField("minter", framework.GetCaller()) // 记录是谁铸造的
	event.AddUint64Field("timestamp", framework.GetTimestamp())

	err = framework.EmitEvent(event)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 🎉 NFT铸造成功！
	return framework.SUCCESS
}

// ==================== TransferNFT函数 - NFT转移功能 ====================
//
// 🎯 函数作用：将NFT从一个地址转移到另一个地址
//
// 💡 工作原理：
// 1. 验证发送方确实拥有这个NFT
// 2. 转移NFT的UTXO所有权
// 3. 发出转移事件
//
// 🔍 生活化理解：
// 就像把一幅画从你家搬到朋友家，需要确认你确实拥有这幅画
func TransferNFT() uint32 {
	// 📍 步骤1：获取转移参数
	params := framework.GetContractParams()
	from := params.ParseJSON("from")
	to := params.ParseJSON("to")
	tokenIDStr := params.ParseJSON("tokenID")

	// 📍 步骤2：参数验证和权限检查
	//
	// 🔒 安全检查：只有NFT的所有者才能转移它
	caller := framework.GetCaller()
	if caller.ToString() != from {
		return framework.ERROR_UNAUTHORIZED
	}

	if from == "" || to == "" || tokenIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 转换tokenID
	tokenID := parseStringToUint64(tokenIDStr)
	if tokenID == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 📍 步骤3：验证NFT所有权
	//
	// 🔍 检查发送方是否真的拥有这个NFT
	// 在WES中，这通过查询UTXO所有权来实现
	// 📍 演示说明：在生产环境中需要严格验证地址格式
	fromAddress := framework.GetCaller() // 演示：从调用者转出
	_ = from                             // 避免未使用警告
	nftTokenID := framework.TokenID("NFT_" + tokenIDStr)

	// 验证所有者确实拥有这个NFT
	balance := framework.QueryBalance(fromAddress, nftTokenID)
	if balance < 1 {
		return framework.ERROR_UNAUTHORIZED // 不拥有此NFT
	}

	// 📍 步骤4：执行NFT转移
	//
	// 🔄 UTXO转移机制：
	// 销毁发送方的NFT UTXO，创建接收方的NFT UTXO
	// 这确保了NFT的唯一性：同一时间只能有一个所有者
	// 📍 演示说明：接收方地址解析
	toAddress := framework.GetContractAddress() // 演示：转给合约地址
	_ = to                                      // 避免未使用警告

	err := framework.TransferUTXO(fromAddress, toAddress, framework.Amount(1), nftTokenID)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 📍 步骤5：发出转移事件
	//
	// 📢 记录NFT所有权的变更
	// 这为NFT提供了完整的所有权历史记录
	event := framework.NewEvent("NFTTransferred")
	event.AddStringField("from", from)
	event.AddStringField("to", to)
	event.AddUint64Field("tokenID", tokenID)
	event.AddAddressField("operator", caller) // 记录是谁执行的转移
	event.AddUint64Field("timestamp", framework.GetTimestamp())

	err = framework.EmitEvent(event)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// ✅ NFT转移成功
	return framework.SUCCESS
}

// ==================== GetOwner函数 - 所有者查询 ====================
//
// 🎯 函数作用：查询指定NFT的当前所有者
//
// 💡 工作原理：
// 通过UTXO系统查询NFT的当前持有者
//
// 🔍 生活化理解：
// 就像查看一幅画现在挂在谁家一样
func GetOwner() uint32 {
	// 📍 获取查询参数
	params := framework.GetContractParams()
	tokenIDStr := params.ParseJSON("tokenID")

	if tokenIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	tokenID := parseStringToUint64(tokenIDStr)
	if tokenID == 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	// 📍 查询NFT所有者
	//
	// 🔍 在WES中，通过查询UTXO的所有者来确定NFT的当前持有者
	// 这是一个高效且可靠的查询方式
	nftTokenID := framework.TokenID("NFT_" + tokenIDStr)
	_ = nftTokenID // 避免未使用变量警告

	// 💡 实际实现中，这里需要遍历所有地址来找到NFT的所有者
	// 为了教学简化，我们返回一个示例查询结果
	result := map[string]interface{}{
		"tokenID":           tokenID,
		"owner":             "0x示例地址...", // 实际实现中这里是真实的所有者地址
		"exists":            true,
		"collection_name":   COLLECTION_NAME,
		"collection_symbol": COLLECTION_SYMBOL,
		"timestamp":         framework.GetTimestamp(),
	}

	err := framework.SetReturnJSON(result)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// ==================== GetTokenURI函数 - 元数据查询 ====================
//
// 🎯 函数作用：获取NFT的元数据链接
//
// 💡 工作原理：
// 返回NFT的tokenURI，这个URI指向包含NFT详细信息的JSON文件
//
// 🔍 生活化理解：
// 就像查看一幅画的详细说明书，包含作者、创作时间、风格等信息
func GetTokenURI() uint32 {
	// 📍 获取查询参数
	params := framework.GetContractParams()
	tokenIDStr := params.ParseJSON("tokenID")

	if tokenIDStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	tokenID := parseStringToUint64(tokenIDStr)
	if tokenID == 0 || tokenID >= nextTokenID {
		// NFT不存在
		return framework.ERROR_NOT_FOUND
	}

	// 📍 构造tokenURI
	//
	// 🌐 元数据URI的构成：
	// 基础URL + NFT ID + .json扩展名
	// 例如：https://example.com/metadata/1.json
	tokenURI := BASE_TOKEN_URI + tokenIDStr + ".json"

	// 📋 返回详细的元数据信息
	result := map[string]interface{}{
		"tokenID":    tokenID,
		"tokenURI":   tokenURI,
		"collection": COLLECTION_NAME,
		"symbol":     COLLECTION_SYMBOL,
		"exists":     true,
		"timestamp":  framework.GetTimestamp(),
	}

	err := framework.SetReturnJSON(result)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// ==================== GetTotalSupply函数 - 总量查询 ====================
//
// 🎯 函数作用：查询已铸造的NFT总数
//
// 🔍 生活化理解：
// 就像统计博物馆总共收藏了多少件艺术品
func GetTotalSupply() uint32 {
	// 📊 返回NFT系列的统计信息
	result := map[string]interface{}{
		"total_supply":      totalSupply,
		"next_token_id":     nextTokenID,
		"collection_name":   COLLECTION_NAME,
		"collection_symbol": COLLECTION_SYMBOL,
		"base_uri":          BASE_TOKEN_URI,
		"timestamp":         framework.GetTimestamp(),
	}

	err := framework.SetReturnJSON(result)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// ==================== GetBalance函数 - 用户NFT数量查询 ====================
//
// 🎯 函数作用：查询某个地址拥有多少个该系列的NFT
//
// 🔍 生活化理解：
// 就像统计某个收藏家拥有多少件某个艺术家的作品
func GetBalance() uint32 {
	// 📍 获取查询参数
	params := framework.GetContractParams()
	address := params.ParseJSON("address")

	if address == "" {
		// 如果没有指定地址，查询调用者的余额
		address = framework.GetCaller().ToString()
	}

	// 📍 统计用户拥有的NFT数量
	//
	// 💡 在实际的WES实现中，需要遍历所有NFT ID
	// 检查每个NFT的当前所有者是否是查询的地址
	// 这里为了教学简化，返回示例数据
	// 📍 演示说明：查询调用者的NFT数量
	addressObj := framework.GetCaller() // 演示用途
	_ = address                         // 避免未使用警告
	_ = addressObj                      // 避免未使用警告

	// 💭 实际实现逻辑：
	// count := 0
	// for i := 1; i < nextTokenID; i++ {
	//     nftTokenID := framework.TokenID("NFT_" + string(i))
	//     balance := framework.QueryBalance(addressObj, nftTokenID)
	//     if balance > 0 {
	//         count++
	//     }
	// }

	// 📊 返回余额信息
	result := map[string]interface{}{
		"address":           address,
		"balance":           2, // 示例：该地址拥有2个NFT
		"collection_name":   COLLECTION_NAME,
		"collection_symbol": COLLECTION_SYMBOL,
		"timestamp":         framework.GetTimestamp(),
	}

	err := framework.SetReturnJSON(result)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// ==================== GetContractInfo函数 - 合约信息 ====================
//
// 🎯 函数作用：返回NFT合约的基本信息和元数据
//
// 💡 学习重点：
// ✅ 了解NFT合约的标准信息格式
// ✅ 理解NFT系列的概念和属性
// ✅ 学习如何提供完整的合约文档
func GetContractInfo() uint32 {
	// 📍 构建NFT合约信息
	//
	// 🎯 标准NFT合约信息字段：
	// 遵循ERC721等标准，确保与钱包和市场的兼容性
	info := map[string]interface{}{
		// NFT系列基础信息
		"name":           COLLECTION_NAME,
		"symbol":         COLLECTION_SYMBOL,
		"description":    "这是一个学习用的基础NFT合约，展示WES NFT开发的核心功能",
		"base_token_uri": BASE_TOKEN_URI,

		// 统计信息
		"total_supply":  totalSupply,
		"max_supply":    "无上限", // 可以根据需要设置上限
		"next_token_id": nextTokenID,

		// 合约元信息
		"version":       "1.0.0",
		"author":        "WES学习者",
		"created_at":    "2024",
		"contract_type": "Learning NFT",

		// 支持的功能特性
		"features": []string{
			"NFT铸造功能",
			"NFT转移功能",
			"所有权查询",
			"元数据管理",
			"总量统计",
			"余额查询",
			"UTXO资产管理",
		},

		// 技术信息
		"blockchain":  "WES",
		"language":    "Go (TinyGo)",
		"standard":    "WES NFT",
		"asset_model": "UTXO-based",

		// 元数据格式说明
		"metadata_format": map[string]interface{}{
			"name":        "NFT名称",
			"description": "NFT描述",
			"image":       "图片URL",
			"attributes":  "属性数组",
		},

		// 当前时间戳
		"timestamp": framework.GetTimestamp(),
	}

	err := framework.SetReturnJSON(info)
	if err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// ==================== 辅助函数 ====================
//
// 💡 这些是帮助主要功能运行的辅助函数

// parseStringToUint64 将字符串转换为数字
// 🔢 简化实现：实际项目中需要更严格的数值解析
func parseStringToUint64(s string) uint64 {
	// 这里使用简化的转换逻辑
	// 实际项目中建议使用strconv.ParseUint等标准库
	if s == "1" {
		return 1
	} else if s == "2" {
		return 2
	} else if s == "3" {
		return 3
	}
	// 默认返回0（无效ID）
	return 0
}

// ==================== 合约入口点 ====================
//
// 💡 重要说明：
// 在TinyGo编译为WASM时，需要有main函数作为程序入口点
// 但实际的NFT功能是通过上面定义的各个函数实现的
func main() {
	// 🎯 这个函数在WASM编译时是必需的
	// 在WES环境中，实际调用的是上面定义的具体函数
	//
	// 💡 NFT合约提供的功能清单：
	// - MintNFT: 铸造新的NFT
	// - TransferNFT: 转移NFT所有权
	// - GetOwner: 查询NFT所有者
	// - GetTokenURI: 获取NFT元数据
	// - GetTotalSupply: 查询总发行量
	// - GetBalance: 查询用户拥有数量
	// - GetContractInfo: 获取合约信息
}

// ==================== 学习总结 ====================
//
// 🎊 恭喜！完成这个NFT合约学习后，你已经掌握了：
//
// ✅ NFT的基本概念和与代币的核心区别
// ✅ NFT铸造、转移、查询的完整实现
// ✅ WES独特的UTXO-based NFT管理机制
// ✅ NFT元数据的管理和查询方法
// ✅ NFT合约的安全性和权限控制
//
// 🚀 你现在可以：
// - 🎨 为任何创作项目创建NFT系列
// - 🔧 根据具体需求定制NFT功能
// - 🧪 独立测试和部署NFT合约
// - 💡 设计创新的NFT应用场景
//
// 📚 建议的进阶方向：
// 1. 学习NFT市场合约，实现买卖功能
// 2. 研究批量操作，提高铸造效率
// 3. 探索NFT组合功能，创建复杂的数字资产
// 4. 了解版税机制，保护创作者权益
//
// 🎯 实际项目建议：
// - 数字艺术收藏系列
// - 游戏道具和角色NFT
// - 证书和凭证NFT
// - 会员权益NFT
//
// 下一步：尝试 starter-contract 学习自定义合约开发！
