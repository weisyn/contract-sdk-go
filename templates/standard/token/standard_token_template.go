//go:build tinygo.wasm

// Package main provides a standard token contract template.
package main

import (
	"unsafe"
)

// ==================== WES 标准代币合约模板 ====================
//
// 🌟 **设计理念**：基于WES标准合约接口规范的代币模板
//
// 🎯 **核心特性**：
// - 实现IContractBase和ITokenStandard标准接口
// - 完全无状态设计，基于UTXO的资产管理
// - 支持标准ERC20功能：转账、授权、查询
// - 内置安全检查和错误处理
// - 事件发出和元数据管理
//
// 📋 **实现接口**：
// - IContractBase: Initialize, GetMetadata, GetVersion
// - ITokenStandard: Transfer, GetBalance, GetTotalSupply, Approve
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

// getBytes 从内存指针获取字节数组
func getBytes(ptr uint32, len uint32) []byte {
	if ptr == 0 || len == 0 {
		return nil
	}
	//nolint:gosec,govet // G103: WASM 环境中需要使用 unsafe.Pointer 与宿主环境交互; unsafeptr: WASM 内存访问需要 uintptr 转换
	return (*[1 << 20]byte)(unsafe.Pointer(uintptr(ptr)))[:len:len]
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

// ==================== IContractBase接口实现 ====================

// Initialize 合约初始化
// 创建初始代币供应并分配给合约部署者
//
//export Initialize
func Initialize() uint32 {
	// 获取初始化参数
	paramsBuffer := malloc(1024)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 1024)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 解析初始化参数（简化实现，实际可使用JSON解析）
	// 期望格式: "name,symbol,decimals,totalSupply"
	params := getString(paramsBuffer, paramLen)
	_ = params // 避免未使用警告

	// 获取合约地址作为初始代币接收者
	contractAddr := malloc(20)
	if contractAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	getContractAddress(contractAddr)

	// 创建初始代币供应的UTXO
	// 这里使用默认值，实际应从params解析
	tokenIDPtr, tokenIDLen := allocateString("STANDARD_TOKEN")
	if tokenIDPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	// 创建1000000个代币的初始供应 (避免uint64溢出)
	initialSupply := uint64(1000000000000) // 1M tokens (简化为12位精度)
	result := createUTXOOutput(contractAddr, initialSupply, tokenIDPtr, tokenIDLen)
	if result != SUCCESS {
		return ERROR_EXECUTION_FAILED
	}

	// 发出初始化事件
	eventData := `{"event":"Initialize","data":{"name":"Standard Token","symbol":"STD","totalSupply":"1000000000000000000000000"}}`
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
		"name": "Standard Token",
		"symbol": "STD",
		"version": "1.0.0",
		"description": "WES标准代币合约模板",
		"author": "WES Development Team",
		"license": "MIT",
		"interfaces": ["IContractBase", "ITokenStandard"],
		"features": ["transfer", "approve", "balance_query"],
		"decimals": 18,
		"totalSupply": "1000000000000000000000000"
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

// ==================== ITokenStandard接口实现 ====================

// Transfer 转账代币
// 通过UTXO转移实现代币转账
//
//export Transfer
func Transfer() uint32 {
	// 获取调用参数（简化实现，实际应解析复杂参数格式）
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

	// 解析参数：to,amount (简化格式)
	params := getString(paramsBuffer, paramLen)
	_ = params // 避免未使用警告
	// 实际实现中应进行完整的参数解析和验证

	// 为演示目的，假设转账参数
	toAddr := malloc(20)
	if toAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	// 查询发送者余额
	tokenIDPtr, tokenIDLen := allocateString("STANDARD_TOKEN")
	if tokenIDPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	senderBalance := queryUTXOBalance(callerAddr, tokenIDPtr, tokenIDLen)
	if senderBalance == 0 {
		return ERROR_INSUFFICIENT_BALANCE
	}

	// 执行UTXO转移（简化实现）
	transferAmount := uint64(1000000000000000000) // 1 token

	result := executeUTXOTransfer(callerAddr, toAddr, transferAmount, tokenIDPtr, tokenIDLen)
	if result != SUCCESS {
		return ERROR_EXECUTION_FAILED
	}

	// 发出转账事件
	eventData := `{"event":"Transfer","data":{"from":"sender","to":"recipient","amount":"1000000000000000000"}}`
	eventPtr, eventLen := allocateString(eventData)
	if eventPtr != 0 {
		emitEvent(eventPtr, eventLen)
	}

	return SUCCESS
}

// GetBalance 查询余额
//
//export GetBalance
func GetBalance() uint32 {
	// 获取查询参数
	paramsBuffer := malloc(1024)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 1024)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 解析地址参数（简化实现）
	queryAddr := malloc(20)
	if queryAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	// 查询指定地址的代币余额
	tokenIDPtr, tokenIDLen := allocateString("STANDARD_TOKEN")
	if tokenIDPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	balance := queryUTXOBalance(queryAddr, tokenIDPtr, tokenIDLen)

	// 返回余额信息
	balanceData := `{"balance":"` + uint64ToString(balance) + `","token":"STANDARD_TOKEN"}`
	balancePtr, balanceLen := allocateString(balanceData)
	if balancePtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	setReturnData(balancePtr, balanceLen)
	return SUCCESS
}

// GetTotalSupply 获取总供应量
//
//export GetTotalSupply
func GetTotalSupply() uint32 {
	// 返回代币总供应量信息
	supplyData := `{"totalSupply":"1000000000000000000000000","token":"STANDARD_TOKEN"}`
	supplyPtr, supplyLen := allocateString(supplyData)
	if supplyPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	setReturnData(supplyPtr, supplyLen)
	return SUCCESS
}

// Approve 授权代币使用权
//
//export Approve
func Approve() uint32 {
	// 获取授权参数
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

	// 授权逻辑（简化实现）
	// 实际实现中需要维护授权关系，由于URES无状态设计，
	// 可通过特殊的UTXO类型或事件记录授权信息

	// 发出授权事件
	eventData := `{"event":"Approval","data":{"owner":"caller","spender":"spender","amount":"1000000000000000000"}}`
	eventPtr, eventLen := allocateString(eventData)
	if eventPtr != 0 {
		emitEvent(eventPtr, eventLen)
	}

	return SUCCESS
}

// ==================== 辅助工具函数 ====================

// uint64ToString 将uint64转换为字符串（简化实现）
func uint64ToString(n uint64) string {
	if n == 0 {
		return "0"
	}

	// 简化的数字转字符串实现
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

// ==================== 主函数（WASM入口点）====================

func main() {
	// WASM模块主入口，通常为空
	// 实际的合约逻辑通过导出的函数调用
}
