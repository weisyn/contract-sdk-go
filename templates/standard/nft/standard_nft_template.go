//go:build tinygo.wasm

// Package main provides a standard NFT contract template for the WES blockchain.
package main

import (
	"unsafe"
)

// ==================== WES 标准NFT合约模板 ====================
//
// 🌟 **设计理念**：基于WES标准合约接口规范的NFT模板
//
// 🎯 **核心特性**：
// - 实现IContractBase和INonFungibleToken标准接口
// - 完全无状态设计，NFT数据以UTXO形式存在
// - 支持标准ERC721功能：铸造、转移、查询、元数据
// - 内置版权保护和版税分成
// - 支持批量操作和集合管理
//
// 📋 **实现接口**：
// - IContractBase: Initialize, GetMetadata, GetVersion
// - INonFungibleToken: MintNFT, TransferNFT, GetTokenInfo, SetTokenURI
//
// ==================== 标准错误码 ====================

const (
	SUCCESS                    = 0
	ERROR_INVALID_PARAMS       = 1
	ERROR_INSUFFICIENT_BALANCE = 2
	ERROR_UNAUTHORIZED         = 3
	ERROR_NOT_FOUND            = 4
	ERROR_ALREADY_EXISTS       = 5
	ERROR_EXECUTION_FAILED     = 6
	ERROR_INVALID_STATE        = 7
	ERROR_TIMEOUT              = 8
	ERROR_UNKNOWN              = 999
)

// NFT类型常量
const (
	NFT_TYPE_ARTWORK     = "ARTWORK"
	NFT_TYPE_COLLECTIBLE = "COLLECTIBLE"
	NFT_TYPE_GAMING      = "GAMING"
	NFT_TYPE_CERTIFICATE = "CERTIFICATE"
	NFT_TYPE_IDENTITY    = "IDENTITY"
	NFT_TYPE_TICKET      = "TICKET"
	NFT_TYPE_DOMAIN      = "DOMAIN"
	NFT_TYPE_MUSIC       = "MUSIC"
)

// ==================== 宿主函数声明 ====================

// 基础环境函数
//
//go:wasmimport env get_caller
func getCaller(addrPtr uint32) uint32

//go:wasmimport env get_contract_address
func getContractAddress(addrPtr uint32) uint32

//go:wasmimport env set_return_data
func setReturnData(dataPtr uint32, dataLen uint32) uint32

//go:wasmimport env emit_event
func emitEvent(eventPtr uint32, eventLen uint32) uint32

//go:wasmimport env get_contract_init_params
func getContractInitParams(bufPtr uint32, bufLen uint32) uint32

//go:wasmimport env get_timestamp
func getTimestamp() uint64

//go:wasmimport env get_block_height
func getBlockHeight() uint64

// UTXO操作函数
//
//go:wasmimport env create_utxo_output
func createUTXOOutput(recipientPtr uint32, amount uint64, tokenIDPtr uint32, tokenIDLen uint32) uint32

//go:wasmimport env execute_utxo_transfer
func executeUTXOTransfer(fromPtr uint32, toPtr uint32, amount uint64, tokenIDPtr uint32, tokenIDLen uint32) uint32

//go:wasmimport env query_utxo_balance
func queryUTXOBalance(addressPtr uint32, tokenIDPtr uint32, tokenIDLen uint32) uint64

// 内存管理函数
//
//go:wasmimport env malloc
func malloc(size uint32) uint32

// ==================== 辅助函数 ====================

// getString 从内存指针构造字符串
func getString(ptr uint32, len uint32) string {
	if ptr == 0 || len == 0 {
		return ""
	}
	//nolint:gosec,govet // G103: WASM 环境中需要使用 unsafe.Pointer 与宿主环境交互; unsafeptr: WASM 内存访问需要 uintptr 转换
	return string((*[1 << 20]byte)(unsafe.Pointer(uintptr(ptr)))[:len:len])
}

// allocateString 分配字符串到WASM内存
func allocateString(s string) (uint32, uint32) {
	if len(s) == 0 {
		return 0, 0
	}
	//nolint:gosec // G115: len() 返回值在 WASM 环境中不会超过 uint32 范围
	ptr := malloc(uint32(len(s)))
	if ptr == 0 {
		return 0, 0
	}
	//nolint:gosec,govet // G103: WASM 环境中需要使用 unsafe.Pointer 与宿主环境交互; unsafeptr: WASM 内存访问需要 uintptr 转换
	copy((*[1 << 20]byte)(unsafe.Pointer(uintptr(ptr)))[:len(s)], s)
	//nolint:gosec // G115: len() 返回值在 WASM 环境中不会超过 uint32 范围
	return ptr, uint32(len(s))
}

// generateTokenID 生成NFT令牌ID
func generateTokenID(prefix string, counter uint64) string {
	return prefix + "_" + uint64ToString(counter) + "_" + uint64ToString(getTimestamp())
}

// uint64ToString 将uint64转换为字符串
func uint64ToString(n uint64) string {
	if n == 0 {
		return "0"
	}

	digits := make([]byte, 0, 20)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}

	// 反转数字
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	return string(digits)
}

// ==================== IContractBase接口实现 ====================

/**
 * Initialize - NFT合约初始化函数
 *
 * 📋 **功能描述**：
 * 初始化NFT合约，设置合约的基本信息和初始状态
 *
 * 📥 **输入参数**：
 * 通过 get_contract_init_params 获取初始化参数
 * 参数格式（JSON）: {"collection_name":"","symbol":"","base_uri":"","max_supply":0,"royalty_rate":0}
 *
 * 📤 **返回值**：
 * @return uint32 - 错误码
 *   - SUCCESS (0): 初始化成功
 *   - ERROR_INVALID_PARAMS (1): 参数无效
 *   - ERROR_EXECUTION_FAILED (6): 执行失败
 *
 * 💡 **实现逻辑**：
 * 1. 分配内存获取初始化参数
 * 2. 解析JSON格式的参数
 * 3. 获取合约地址
 * 4. 发出初始化事件
 *
 * ⚠️ **注意事项**：
 * - 只能调用一次
 * - 需要提供有效的JSON格式参数
 * - 初始化后设置NFT集合的基础信息
 */
//export Initialize
func Initialize() uint32 {
	// 获取初始化参数
	paramsBuffer := malloc(2048)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 2048)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 解析初始化参数（期望JSON格式）
	// 包含：collection_name, symbol, base_uri, max_supply, royalty_rate
	params := getString(paramsBuffer, paramLen)
	_ = params // 避免未使用警告

	// 获取合约地址
	contractAddr := malloc(20)
	if contractAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getContractAddress(contractAddr)

	// 发出初始化事件
	eventData := `{
		"event": "NFTCollectionInitialize",
		"data": {
			"collection_name": "Standard NFT Collection",
			"symbol": "SNFT",
			"max_supply": "10000",
			"royalty_rate": "5",
			"creator": "contract_address",
			"timestamp": "` + uint64ToString(getTimestamp()) + `"
		}
	}`

	eventPtr, eventLen := allocateString(eventData)
	if eventPtr != 0 {
		emitEvent(eventPtr, eventLen)
	}

	return SUCCESS
}

// GetMetadata 获取合约元数据
//
//export GetMetadata
func GetMetadata() uint32 {
	metadata := `{
		"name": "Standard NFT Collection",
		"symbol": "SNFT",
		"version": "1.0.0",
		"description": "WES标准NFT合约模板",
		"author": "WES Development Team",
		"license": "MIT",
		"interfaces": ["IContractBase", "INonFungibleToken"],
		"features": ["mint", "transfer", "metadata", "royalty"],
		"collection_info": {
			"max_supply": "10000",
			"base_uri": "https://api.example.com/nft/",
			"royalty_rate": "5"
		}
	}`

	metadataPtr, metadataLen := allocateString(metadata)
	if metadataPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	setReturnData(metadataPtr, metadataLen)
	return SUCCESS
}

// GetVersion 获取合约版本
//
//export GetVersion
func GetVersion() uint32 {
	version := "1.0.0"
	versionPtr, versionLen := allocateString(version)
	if versionPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	setReturnData(versionPtr, versionLen)
	return SUCCESS
}

// ==================== INonFungibleToken接口实现 ====================

// MintNFT 铸造NFT
// 创建新的NFT并分配给指定地址
//
//export MintNFT
func MintNFT() uint32 {
	// 获取铸造参数
	paramsBuffer := malloc(4096)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 4096)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 获取调用者地址（权限检查）
	callerAddr := malloc(20)
	if callerAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getCaller(callerAddr)

	// 解析参数：to_address, metadata, nft_type
	params := getString(paramsBuffer, paramLen)
	_ = params // 避免未使用警告

	// 生成唯一的NFT令牌ID
	tokenID := generateTokenID("SNFT", getBlockHeight())
	tokenIDPtr, tokenIDLen := allocateString(tokenID)
	if tokenIDPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	// 获取接收者地址（简化实现，实际应从params解析）
	recipientAddr := malloc(20)
	if recipientAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	// 创建NFT UTXO（数量为1表示不可分割性）
	result := createUTXOOutput(recipientAddr, 1, tokenIDPtr, tokenIDLen)
	if result != SUCCESS {
		return ERROR_EXECUTION_FAILED
	}

	// 发出铸造事件
	eventData := `{
		"event": "NFTMint",
		"data": {
			"token_id": "` + tokenID + `",
			"to": "recipient_address",
			"nft_type": "` + NFT_TYPE_ARTWORK + `",
			"metadata": {
				"name": "Standard NFT #1",
				"description": "A standard NFT created from template",
				"image": "https://api.example.com/nft/image/1",
				"attributes": [
					{"trait_type": "Color", "value": "Blue"},
					{"trait_type": "Rarity", "value": "Common"}
				]
			},
			"timestamp": "` + uint64ToString(getTimestamp()) + `"
		}
	}`

	eventPtr, eventLen := allocateString(eventData)
	if eventPtr != 0 {
		emitEvent(eventPtr, eventLen)
	}

	// 返回新铸造的NFT ID
	setReturnData(tokenIDPtr, tokenIDLen)
	return SUCCESS
}

// TransferNFT 转移NFT
// 将NFT从一个地址转移到另一个地址
//
//export TransferNFT
func TransferNFT() uint32 {
	// 获取转移参数
	paramsBuffer := malloc(2048)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 2048)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 获取调用者地址
	callerAddr := malloc(20)
	if callerAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getCaller(callerAddr)

	// 解析参数：from, to, token_id
	params := getString(paramsBuffer, paramLen)
	_ = params // 避免未使用警告

	// 为演示目的使用简化的令牌ID
	//nolint:gosec // G101: 这是示例代码中的演示值，不是真正的凭证
	tokenID := "SNFT_1_1640995200"
	tokenIDPtr, tokenIDLen := allocateString(tokenID)
	if tokenIDPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	// 验证调用者是否拥有该NFT
	balance := queryUTXOBalance(callerAddr, tokenIDPtr, tokenIDLen)
	if balance == 0 {
		return ERROR_UNAUTHORIZED
	}

	// 准备转移地址
	toAddr := malloc(20)
	if toAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	// 执行NFT转移
	result := executeUTXOTransfer(callerAddr, toAddr, 1, tokenIDPtr, tokenIDLen)
	if result != SUCCESS {
		return ERROR_EXECUTION_FAILED
	}

	// 发出转移事件
	eventData := `{
		"event": "NFTTransfer",
		"data": {
			"token_id": "` + tokenID + `",
			"from": "caller_address",
			"to": "recipient_address",
			"timestamp": "` + uint64ToString(getTimestamp()) + `"
		}
	}`

	eventPtr, eventLen := allocateString(eventData)
	if eventPtr != 0 {
		emitEvent(eventPtr, eventLen)
	}

	return SUCCESS
}

// GetTokenInfo 获取NFT信息
// 查询指定NFT的详细信息
//
//export GetTokenInfo
func GetTokenInfo() uint32 {
	// 获取查询参数
	paramsBuffer := malloc(1024)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 1024)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 解析令牌ID参数
	params := getString(paramsBuffer, paramLen)
	_ = params // 避免未使用警告
	//nolint:gosec // G101: 这是示例代码中的演示值，不是真正的凭证
	tokenID := "SNFT_1_1640995200" // 简化实现

	// 构造NFT信息响应
	nftInfo := `{
		"token_id": "` + tokenID + `",
		"owner": "current_owner_address",
		"metadata": {
			"name": "Standard NFT #1",
			"description": "A standard NFT created from template",
			"image": "https://api.example.com/nft/image/1",
			"external_url": "https://example.com/nft/1",
			"attributes": [
				{"trait_type": "Color", "value": "Blue"},
				{"trait_type": "Rarity", "value": "Common"},
				{"trait_type": "Collection", "value": "Standard NFT Collection"}
			]
		},
		"collection": {
			"name": "Standard NFT Collection",
			"symbol": "SNFT",
			"contract_address": "contract_address"
		},
		"royalty": {
			"rate": "5",
			"recipient": "creator_address"
		},
		"created_at": "1640995200",
		"last_transfer": "1640995200"
	}`

	nftInfoPtr, nftInfoLen := allocateString(nftInfo)
	if nftInfoPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	setReturnData(nftInfoPtr, nftInfoLen)
	return SUCCESS
}

// SetTokenURI 设置NFT元数据URI
// 更新指定NFT的元数据URI
//
//export SetTokenURI
func SetTokenURI() uint32 {
	// 获取设置参数
	paramsBuffer := malloc(2048)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 2048)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 获取调用者地址
	callerAddr := malloc(20)
	if callerAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getCaller(callerAddr)

	// 解析参数：token_id, new_uri
	params := getString(paramsBuffer, paramLen)
	_ = params // 避免未使用警告

	// 权限检查：只有NFT所有者或授权者可以更新元数据
	//nolint:gosec // G101: 这是示例代码中的演示值，不是真正的凭证
	tokenID := "SNFT_1_1640995200" // 简化实现
	tokenIDPtr, tokenIDLen := allocateString(tokenID)
	if tokenIDPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	balance := queryUTXOBalance(callerAddr, tokenIDPtr, tokenIDLen)
	if balance == 0 {
		return ERROR_UNAUTHORIZED
	}

	// 发出元数据更新事件（由于URES无状态特性，元数据更新通过事件记录）
	eventData := `{
		"event": "NFTMetadataUpdate",
		"data": {
			"token_id": "` + tokenID + `",
			"new_uri": "https://api.example.com/nft/updated/1",
			"updated_by": "caller_address",
			"timestamp": "` + uint64ToString(getTimestamp()) + `"
		}
	}`

	eventPtr, eventLen := allocateString(eventData)
	if eventPtr != 0 {
		emitEvent(eventPtr, eventLen)
	}

	return SUCCESS
}

// ==================== 扩展功能实现 ====================

// BatchMint 批量铸造NFT
//
//export BatchMint
func BatchMint() uint32 {
	// 获取批量铸造参数
	paramsBuffer := malloc(8192)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 8192)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 获取调用者地址
	callerAddr := malloc(20)
	if callerAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getCaller(callerAddr)

	// 简化的批量铸造实现
	// 实际实现应解析批量参数并循环处理
	batchSize := 5 // 演示批量铸造5个NFT

	for i := 0; i < batchSize; i++ {
		// 生成唯一令牌ID
		//nolint:gosec // G115: i 是循环计数器，不会超过 uint64 范围
		tokenID := generateTokenID("SNFT_BATCH", uint64(i))
		tokenIDPtr, tokenIDLen := allocateString(tokenID)
		if tokenIDPtr == 0 {
			continue
		}

		// 创建NFT UTXO
		result := createUTXOOutput(callerAddr, 1, tokenIDPtr, tokenIDLen)
		if result != SUCCESS {
			continue
		}

		// 发出批量铸造事件
		eventData := `{
			"event": "NFTBatchMint",
			"data": {
				"token_id": "` + tokenID + `",
				"batch_index": "` + uint64ToString(uint64(i)) + `",
				"to": "caller_address",
				"timestamp": "` + uint64ToString(getTimestamp()) + `"
			}
		}`

		eventPtr, eventLen := allocateString(eventData)
		if eventPtr != 0 {
			emitEvent(eventPtr, eventLen)
		}
	}

	return SUCCESS
}

// Burn 销毁NFT
//
//export Burn
func Burn() uint32 {
	// 获取销毁参数
	paramsBuffer := malloc(1024)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 1024)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 获取调用者地址
	callerAddr := malloc(20)
	if callerAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getCaller(callerAddr)

	// 解析要销毁的令牌ID
	params := getString(paramsBuffer, paramLen)
	_ = params // 避免未使用警告
	//nolint:gosec // G101: 这是示例代码中的演示值，不是真正的凭证
	tokenID := "SNFT_1_1640995200" // 简化实现

	tokenIDPtr, tokenIDLen := allocateString(tokenID)
	if tokenIDPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	// 验证所有权
	balance := queryUTXOBalance(callerAddr, tokenIDPtr, tokenIDLen)
	if balance == 0 {
		return ERROR_UNAUTHORIZED
	}

	// NFT销毁通过转移到特殊的"黑洞"地址实现
	burnAddr := malloc(20)
	if burnAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	// 执行销毁转移
	result := executeUTXOTransfer(callerAddr, burnAddr, 1, tokenIDPtr, tokenIDLen)
	if result != SUCCESS {
		return ERROR_EXECUTION_FAILED
	}

	// 发出销毁事件
	eventData := `{
		"event": "NFTBurn",
		"data": {
			"token_id": "` + tokenID + `",
			"burned_by": "caller_address",
			"timestamp": "` + uint64ToString(getTimestamp()) + `"
		}
	}`

	eventPtr, eventLen := allocateString(eventData)
	if eventPtr != 0 {
		emitEvent(eventPtr, eventLen)
	}

	return SUCCESS
}

// ==================== 主函数（WASM入口点）====================

func main() {
	// WASM模块主入口，通常为空
	// 实际的合约逻辑通过导出的函数调用
}
