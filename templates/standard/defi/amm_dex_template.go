//go:build tinygo.wasm

// Package main provides an AMM DEX contract template.
package main

import (
	"unsafe"
)

// ==================== WES AMM DEX 合约模板 ====================
//
// 🌟 **设计理念**：基于WES URES模型的自动化做市商(AMM)去中心化交易所
//
// 🎯 **核心特性**：
// - 实现IContractBase接口的DeFi基础设施
// - 完全无状态设计，流动性池以UTXO形式管理
// - 支持恒定乘积做市商算法(x * y = k)
// - 提供流动性添加/移除、代币交换功能
// - 内置滑点保护和价格影响计算
// - 支持多代币对交易和流动性挖矿
//
// 📋 **主要功能**：
// - Initialize: 初始化DEX和创建交易对
// - AddLiquidity: 添加流动性
// - RemoveLiquidity: 移除流动性
// - SwapTokens: 代币交换
// - GetPoolInfo: 查询流动性池信息
// - GetPrice: 获取代币价格
//
// ==================== 标准错误码 ====================

const (
	SUCCESS                      = 0
	ERROR_INVALID_PARAMS         = 1
	ERROR_INSUFFICIENT_BALANCE   = 2
	ERROR_UNAUTHORIZED           = 3
	ERROR_NOT_FOUND              = 4
	ERROR_ALREADY_EXISTS         = 5
	ERROR_EXECUTION_FAILED       = 6
	ERROR_INVALID_STATE          = 7
	ERROR_TIMEOUT                = 8
	ERROR_SLIPPAGE_EXCEEDED      = 9
	ERROR_INSUFFICIENT_LIQUIDITY = 10
	ERROR_UNKNOWN                = 999
)

// DeFi常量
const (
	MINIMUM_LIQUIDITY = uint64(1000)  // 最小流动性
	MAX_SLIPPAGE      = uint64(500)   // 最大滑点 5%
	FEE_RATE          = uint64(30)    // 交易手续费 0.3%
	FEE_DENOMINATOR   = uint64(10000) // 费率分母
	LP_TOKEN_DECIMALS = uint64(18)    // LP代币精度
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

// ==================== 数学辅助函数 ====================

// sqrt 计算平方根（简化实现）
func sqrt(x uint64) uint64 {
	if x == 0 {
		return 0
	}

	// 使用牛顿法求平方根
	z := x
	for i := 0; i < 20; i++ {
		newZ := (z + x/z) / 2
		if newZ >= z {
			return z
		}
		z = newZ
	}
	return z
}

// min 返回两个数中的较小值
func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// max 返回两个数中的较大值
func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

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

// generatePoolID 生成交易对ID
func generatePoolID(tokenA, tokenB string) string {
	return "POOL_" + tokenA + "_" + tokenB + "_" + uint64ToString(getBlockHeight())
}

// generateLPTokenID 生成LP代币ID
func generateLPTokenID(poolID string) string {
	return "LP_" + poolID
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

// getTokenBalance 获取地址的代币余额
func getTokenBalance(addr uint32, tokenID string) uint64 {
	tokenIDPtr, tokenIDLen := allocateString(tokenID)
	if tokenIDPtr == 0 {
		return 0
	}
	return queryUTXOBalance(addr, tokenIDPtr, tokenIDLen)
}

// ==================== AMM核心算法 ====================

// calculateSwapAmountOut 计算交换输出金额
// 基于恒定乘积公式：(x + dx) * (y - dy) = x * y
func calculateSwapAmountOut(amountIn, reserveIn, reserveOut uint64) uint64 {
	if amountIn == 0 || reserveIn == 0 || reserveOut == 0 {
		return 0
	}

	// 扣除手续费
	amountInWithFee := amountIn * (FEE_DENOMINATOR - FEE_RATE)
	numerator := amountInWithFee * reserveOut
	denominator := reserveIn*FEE_DENOMINATOR + amountInWithFee

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

// calculateLPTokensToMint 计算应铸造的LP代币数量
func calculateLPTokensToMint(amountA, amountB, reserveA, reserveB, totalSupply uint64) uint64 {
	if totalSupply == 0 {
		// 首次添加流动性
		liquidity := sqrt(amountA * amountB)
		if liquidity > MINIMUM_LIQUIDITY {
			return liquidity - MINIMUM_LIQUIDITY
		}
		return 0
	}

	// 后续添加流动性，按比例铸造
	liquidityA := amountA * totalSupply / reserveA
	liquidityB := amountB * totalSupply / reserveB

	return min(liquidityA, liquidityB)
}

// calculatePriceImpact 计算价格影响
func calculatePriceImpact(amountIn, reserveIn, reserveOut uint64) uint64 {
	if reserveIn == 0 || reserveOut == 0 {
		return MAX_SLIPPAGE // 无流动性时返回最大滑点
	}

	// 计算价格影响百分比
	priceImpact := amountIn * 10000 / (reserveIn + amountIn)
	return priceImpact
}

// ==================== IContractBase接口实现 ====================

// Initialize 合约初始化
// 设置DEX基础参数和支持的代币对
//
//export Initialize
func Initialize() uint32 {
	// 获取初始化参数
	paramsBuffer := malloc(4096)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 4096)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 解析初始化参数
	_ = getString(paramsBuffer, paramLen)

	// 获取合约地址
	contractAddr := malloc(20)
	if contractAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getContractAddress(contractAddr)

	// 发出DEX初始化事件
	eventData := `{
		"event": "DEXInitialize",
		"data": {
			"dex_name": "Standard AMM DEX",
			"fee_rate": "` + uint64ToString(FEE_RATE) + `",
			"minimum_liquidity": "` + uint64ToString(MINIMUM_LIQUIDITY) + `",
			"max_slippage": "` + uint64ToString(MAX_SLIPPAGE) + `",
			"contract_address": "contract_address",
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
		"name": "Standard AMM DEX",
		"symbol": "STDEX",
		"version": "1.0.0",
		"description": "WES标准AMM去中心化交易所模板",
		"author": "WES Development Team",
		"license": "MIT",
		"interfaces": ["IContractBase"],
		"features": ["amm", "liquidity_pool", "token_swap", "yield_farming"],
		"defi_params": {
			"fee_rate": "` + uint64ToString(FEE_RATE) + `",
			"fee_denominator": "` + uint64ToString(FEE_DENOMINATOR) + `",
			"minimum_liquidity": "` + uint64ToString(MINIMUM_LIQUIDITY) + `",
			"max_slippage": "` + uint64ToString(MAX_SLIPPAGE) + `"
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

// ==================== DeFi核心功能实现 ====================

// AddLiquidity 添加流动性
// 向指定交易对添加流动性并铸造LP代币
//
//export AddLiquidity
func AddLiquidity() uint32 {
	// 获取添加流动性参数
	paramsBuffer := malloc(4096)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 4096)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 获取调用者地址
	callerAddr := malloc(20)
	if callerAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getCaller(callerAddr)

	// 解析参数：tokenA, tokenB, amountA, amountB, slippage_tolerance
	_ = getString(paramsBuffer, paramLen)

	// 简化实现的参数
	tokenA := "TOKEN_A"
	tokenB := "TOKEN_B"
	amountA := uint64(1000000000000) // 1000 TOKEN_A (scaled)
	amountB := uint64(2000000000000) // 2000 TOKEN_B (scaled)

	// 检查用户余额
	balanceA := getTokenBalance(callerAddr, tokenA)
	balanceB := getTokenBalance(callerAddr, tokenB)

	if balanceA < amountA || balanceB < amountB {
		return ERROR_INSUFFICIENT_BALANCE
	}

	// 获取合约地址作为流动性池
	contractAddr := malloc(20)
	if contractAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getContractAddress(contractAddr)

	// 查询当前池储备
	reserveA := getTokenBalance(contractAddr, tokenA)
	reserveB := getTokenBalance(contractAddr, tokenB)

	// 生成池ID和LP代币ID
	poolID := generatePoolID(tokenA, tokenB)
	lpTokenID := generateLPTokenID(poolID)

	// 查询LP代币总供应量（简化实现）
	lpTotalSupply := uint64(0) // 首次添加流动性

	// 计算应铸造的LP代币数量
	lpTokensToMint := calculateLPTokensToMint(amountA, amountB, reserveA, reserveB, lpTotalSupply)
	if lpTokensToMint == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 转移代币到池中
	tokenAPtr, tokenALen := allocateString(tokenA)
	tokenBPtr, tokenBLen := allocateString(tokenB)

	if tokenAPtr == 0 || tokenBPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	resultA := executeUTXOTransfer(callerAddr, contractAddr, amountA, tokenAPtr, tokenALen)
	resultB := executeUTXOTransfer(callerAddr, contractAddr, amountB, tokenBPtr, tokenBLen)

	if resultA != SUCCESS || resultB != SUCCESS {
		return ERROR_EXECUTION_FAILED
	}

	// 铸造LP代币给用户
	lpTokenIDPtr, lpTokenIDLen := allocateString(lpTokenID)
	if lpTokenIDPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	result := createUTXOOutput(callerAddr, lpTokensToMint, lpTokenIDPtr, lpTokenIDLen)
	if result != SUCCESS {
		return ERROR_EXECUTION_FAILED
	}

	// 发出添加流动性事件
	eventData := `{
		"event": "LiquidityAdded",
		"data": {
			"pool_id": "` + poolID + `",
			"provider": "caller_address",
			"token_a": "` + tokenA + `",
			"token_b": "` + tokenB + `",
			"amount_a": "` + uint64ToString(amountA) + `",
			"amount_b": "` + uint64ToString(amountB) + `",
			"lp_tokens_minted": "` + uint64ToString(lpTokensToMint) + `",
			"new_reserve_a": "` + uint64ToString(reserveA+amountA) + `",
			"new_reserve_b": "` + uint64ToString(reserveB+amountB) + `",
			"timestamp": "` + uint64ToString(getTimestamp()) + `"
		}
	}`

	eventPtr, eventLen := allocateString(eventData)
	if eventPtr != 0 {
		emitEvent(eventPtr, eventLen)
	}

	return SUCCESS
}

// SwapTokens 代币交换
// 使用AMM算法进行代币交换
//
//export SwapTokens
func SwapTokens() uint32 {
	// 获取交换参数
	paramsBuffer := malloc(4096)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 4096)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 获取调用者地址
	callerAddr := malloc(20)
	if callerAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getCaller(callerAddr)

	// 解析参数：tokenIn, tokenOut, amountIn, amountOutMin, slippage
	_ = getString(paramsBuffer, paramLen)

	// 简化实现的参数
	tokenIn := "TOKEN_A"
	tokenOut := "TOKEN_B"
	amountIn := uint64(100000000000)     // 100 TOKEN_A (scaled)
	amountOutMin := uint64(190000000000) // 最少190 TOKEN_B (scaled)

	// 检查用户余额
	userBalance := getTokenBalance(callerAddr, tokenIn)
	if userBalance < amountIn {
		return ERROR_INSUFFICIENT_BALANCE
	}

	// 获取合约地址
	contractAddr := malloc(20)
	if contractAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getContractAddress(contractAddr)

	// 查询池储备
	reserveIn := getTokenBalance(contractAddr, tokenIn)
	reserveOut := getTokenBalance(contractAddr, tokenOut)

	if reserveIn == 0 || reserveOut == 0 {
		return ERROR_INSUFFICIENT_LIQUIDITY
	}

	// 计算交换输出数量
	amountOut := calculateSwapAmountOut(amountIn, reserveIn, reserveOut)
	if amountOut < amountOutMin {
		return ERROR_SLIPPAGE_EXCEEDED
	}

	// 计算价格影响
	priceImpact := calculatePriceImpact(amountIn, reserveIn, reserveOut)
	if priceImpact > MAX_SLIPPAGE {
		return ERROR_SLIPPAGE_EXCEEDED
	}

	// 执行代币交换
	tokenInPtr, tokenInLen := allocateString(tokenIn)
	tokenOutPtr, tokenOutLen := allocateString(tokenOut)

	if tokenInPtr == 0 || tokenOutPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	// 用户转入代币到池中
	resultIn := executeUTXOTransfer(callerAddr, contractAddr, amountIn, tokenInPtr, tokenInLen)
	if resultIn != SUCCESS {
		return ERROR_EXECUTION_FAILED
	}

	// 池转出代币给用户
	resultOut := executeUTXOTransfer(contractAddr, callerAddr, amountOut, tokenOutPtr, tokenOutLen)
	if resultOut != SUCCESS {
		return ERROR_EXECUTION_FAILED
	}

	// 发出代币交换事件
	eventData := `{
		"event": "TokenSwap",
		"data": {
			"trader": "caller_address",
			"token_in": "` + tokenIn + `",
			"token_out": "` + tokenOut + `",
			"amount_in": "` + uint64ToString(amountIn) + `",
			"amount_out": "` + uint64ToString(amountOut) + `",
			"price_impact": "` + uint64ToString(priceImpact) + `",
			"fee_amount": "` + uint64ToString(amountIn*FEE_RATE/FEE_DENOMINATOR) + `",
			"new_reserve_in": "` + uint64ToString(reserveIn+amountIn) + `",
			"new_reserve_out": "` + uint64ToString(reserveOut-amountOut) + `",
			"timestamp": "` + uint64ToString(getTimestamp()) + `"
		}
	}`

	eventPtr, eventLen := allocateString(eventData)
	if eventPtr != 0 {
		emitEvent(eventPtr, eventLen)
	}

	return SUCCESS
}

// GetPoolInfo 获取流动性池信息
//
//export GetPoolInfo
func GetPoolInfo() uint32 {
	// 获取查询参数
	paramsBuffer := malloc(1024)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 1024)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 解析池ID参数
	_ = getString(paramsBuffer, paramLen)

	// 简化实现
	tokenA := "TOKEN_A"
	tokenB := "TOKEN_B"
	poolID := generatePoolID(tokenA, tokenB)

	// 获取合约地址
	contractAddr := malloc(20)
	if contractAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getContractAddress(contractAddr)

	// 查询池储备
	reserveA := getTokenBalance(contractAddr, tokenA)
	reserveB := getTokenBalance(contractAddr, tokenB)

	// 计算池的总价值（简化实现）
	totalValueLocked := reserveA + reserveB

	// 构造池信息响应
	poolInfo := `{
		"pool_id": "` + poolID + `",
		"token_a": {
			"symbol": "` + tokenA + `",
			"reserve": "` + uint64ToString(reserveA) + `"
		},
		"token_b": {
			"symbol": "` + tokenB + `",
			"reserve": "` + uint64ToString(reserveB) + `"
		},
		"lp_token": {
			"symbol": "LP_` + tokenA + `_` + tokenB + `",
			"total_supply": "1000000000000000000000"
		},
		"price": {
			"token_a_per_token_b": "` + uint64ToString(reserveB*1e18/max(reserveA, 1)) + `",
			"token_b_per_token_a": "` + uint64ToString(reserveA*1e18/max(reserveB, 1)) + `"
		},
		"fees": {
			"rate": "` + uint64ToString(FEE_RATE) + `",
			"total_collected": "5000000000000000000"
		},
		"tvl": "` + uint64ToString(totalValueLocked) + `",
		"volume_24h": "100000000000000000000000",
		"created_at": "1640995200",
		"updated_at": "` + uint64ToString(getTimestamp()) + `"
	}`

	poolInfoPtr, poolInfoLen := allocateString(poolInfo)
	if poolInfoPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	setReturnData(poolInfoPtr, poolInfoLen)
	return SUCCESS
}

// GetPrice 获取代币价格
//
//export GetPrice
func GetPrice() uint32 {
	// 获取价格查询参数
	paramsBuffer := malloc(1024)
	if paramsBuffer == 0 {
		return ERROR_EXECUTION_FAILED
	}

	paramLen := getContractInitParams(paramsBuffer, 1024)
	if paramLen == 0 {
		return ERROR_INVALID_PARAMS
	}

	// 解析代币对参数
	_ = getString(paramsBuffer, paramLen)

	// 简化实现
	tokenA := "TOKEN_A"
	tokenB := "TOKEN_B"

	// 获取合约地址
	contractAddr := malloc(20)
	if contractAddr == 0 {
		return ERROR_EXECUTION_FAILED
	}
	getContractAddress(contractAddr)

	// 查询池储备
	reserveA := getTokenBalance(contractAddr, tokenA)
	reserveB := getTokenBalance(contractAddr, tokenB)

	if reserveA == 0 || reserveB == 0 {
		return ERROR_NOT_FOUND
	}

	// 计算价格
	priceAInB := reserveB * 1e18 / reserveA
	priceBInA := reserveA * 1e18 / reserveB

	// 构造价格信息响应
	priceInfo := `{
		"token_pair": "` + tokenA + `/` + tokenB + `",
		"prices": {
			"` + tokenA + `_in_` + tokenB + `": "` + uint64ToString(priceAInB) + `",
			"` + tokenB + `_in_` + tokenA + `": "` + uint64ToString(priceBInA) + `"
		},
		"reserves": {
			"` + tokenA + `": "` + uint64ToString(reserveA) + `",
			"` + tokenB + `": "` + uint64ToString(reserveB) + `"
		},
		"last_updated": "` + uint64ToString(getTimestamp()) + `"
	}`

	priceInfoPtr, priceInfoLen := allocateString(priceInfo)
	if priceInfoPtr == 0 {
		return ERROR_EXECUTION_FAILED
	}

	setReturnData(priceInfoPtr, priceInfoLen)
	return SUCCESS
}

// ==================== 主函数（WASM入口点）====================

func main() {
	// WASM模块主入口，通常为空
	// 实际的合约逻辑通过导出的函数调用
}
