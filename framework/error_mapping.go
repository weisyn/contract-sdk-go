//go:build !tinygo && !(js && wasm)

package framework

// ==================== 错误码映射到 WES 错误码 ====================
//
// 本文件提供合约错误码到 WES Problem Details 错误码的映射
// 注意：此文件仅在非合约环境中编译（用于文档和测试）
//
// 合约执行时，错误码会被区块链服务层（weisyn.git）捕获并转换为 Problem Details
// 本文件提供映射关系，帮助开发者理解错误码的含义

// ContractErrorCodeToWESCode 将合约错误码映射到 WES 错误码
func ContractErrorCodeToWESCode(contractErrorCode uint32) string {
	switch contractErrorCode {
	case SUCCESS:
		return "" // 成功，无错误码
	case ERROR_INVALID_PARAMS:
		return "COMMON_VALIDATION_ERROR"
	case ERROR_INSUFFICIENT_BALANCE:
		return "BC_INSUFFICIENT_BALANCE"
	case ERROR_UNAUTHORIZED:
		return "COMMON_VALIDATION_ERROR" // 授权错误属于验证错误
	case ERROR_NOT_FOUND:
		return "BC_CONTRACT_NOT_FOUND" // 或 BC_TX_NOT_FOUND，取决于上下文
	case ERROR_ALREADY_EXISTS:
		return "COMMON_VALIDATION_ERROR"
	case ERROR_EXECUTION_FAILED:
		return "BC_CONTRACT_INVOCATION_FAILED"
	case ERROR_INVALID_STATE:
		return "BC_CONTRACT_INVOCATION_FAILED"
	case ERROR_TIMEOUT:
		return "COMMON_TIMEOUT"
	case ERROR_NOT_IMPLEMENTED:
		return "BC_CONTRACT_INVOCATION_FAILED"
	case ERROR_PERMISSION_DENIED:
		return "COMMON_VALIDATION_ERROR"
	case ERROR_UNKNOWN:
		return "COMMON_INTERNAL_ERROR"
	default:
		return "COMMON_INTERNAL_ERROR"
	}
}

// ContractErrorCodeToUserMessage 将合约错误码映射到用户友好的消息
func ContractErrorCodeToUserMessage(contractErrorCode uint32) string {
	switch contractErrorCode {
	case SUCCESS:
		return ""
	case ERROR_INVALID_PARAMS:
		return "参数验证失败，请检查输入参数。"
	case ERROR_INSUFFICIENT_BALANCE:
		return "余额不足，无法完成交易。"
	case ERROR_UNAUTHORIZED:
		return "未授权操作，请检查权限。"
	case ERROR_NOT_FOUND:
		return "资源不存在。"
	case ERROR_ALREADY_EXISTS:
		return "资源已存在。"
	case ERROR_EXECUTION_FAILED:
		return "合约执行失败，请检查合约逻辑。"
	case ERROR_INVALID_STATE:
		return "合约状态无效，请检查合约状态。"
	case ERROR_TIMEOUT:
		return "执行超时，请稍后重试。"
	case ERROR_NOT_IMPLEMENTED:
		return "功能未实现。"
	case ERROR_PERMISSION_DENIED:
		return "权限不足，无法执行此操作。"
	case ERROR_UNKNOWN:
		return "未知错误，请稍后重试或联系管理员。"
	default:
		return "未知错误，请稍后重试或联系管理员。"
	}
}

// ContractErrorCodeToHTTPStatus 将合约错误码映射到 HTTP 状态码
func ContractErrorCodeToHTTPStatus(contractErrorCode uint32) int {
	switch contractErrorCode {
	case SUCCESS:
		return 200
	case ERROR_INVALID_PARAMS:
		return 400
	case ERROR_INSUFFICIENT_BALANCE:
		return 422
	case ERROR_UNAUTHORIZED:
		return 401
	case ERROR_NOT_FOUND:
		return 404
	case ERROR_ALREADY_EXISTS:
		return 409
	case ERROR_EXECUTION_FAILED:
		return 422
	case ERROR_INVALID_STATE:
		return 422
	case ERROR_TIMEOUT:
		return 408
	case ERROR_NOT_IMPLEMENTED:
		return 501
	case ERROR_PERMISSION_DENIED:
		return 403
	case ERROR_UNKNOWN:
		return 500
	default:
		return 500
	}
}

