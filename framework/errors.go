//go:build tinygo || (js && wasm)

package framework

// ==================== 错误定义 ====================
//
// 🌟 **设计理念**：统一的错误码和错误处理
//
// 🎯 **核心特性**：
// - 标准错误码：与HostABI错误码对齐
// - 类型安全：ContractError类型
// - 错误信息：提供详细的错误消息

// ==================== 错误定义扩展 ====================
//
// 注意：基础错误码和ContractError类型定义在contract_base.go中
// 本文件提供错误处理的扩展功能

// ==================== 错误检查辅助函数 ====================

// IsSuccess 检查错误码是否为成功
func IsSuccess(code uint32) bool {
	return code == SUCCESS
}

// IsError 检查错误码是否为错误
func IsError(code uint32) bool {
	return code != SUCCESS
}

// ErrorCodeToString 将错误码转换为字符串
func ErrorCodeToString(code uint32) string {
	switch code {
	case SUCCESS:
		return "SUCCESS"
	case ERROR_INVALID_PARAMS:
		return "ERROR_INVALID_PARAMS"
	case ERROR_INSUFFICIENT_BALANCE:
		return "ERROR_INSUFFICIENT_BALANCE"
	case ERROR_UNAUTHORIZED:
		return "ERROR_UNAUTHORIZED"
	case ERROR_NOT_FOUND:
		return "ERROR_NOT_FOUND"
	case ERROR_ALREADY_EXISTS:
		return "ERROR_ALREADY_EXISTS"
	case ERROR_EXECUTION_FAILED:
		return "ERROR_EXECUTION_FAILED"
	case ERROR_INVALID_STATE:
		return "ERROR_INVALID_STATE"
	case ERROR_TIMEOUT:
		return "ERROR_TIMEOUT"
	case ERROR_NOT_IMPLEMENTED:
		return "ERROR_NOT_IMPLEMENTED"
	case ERROR_PERMISSION_DENIED:
		return "ERROR_PERMISSION_DENIED"
	case ERROR_UNKNOWN:
		return "ERROR_UNKNOWN"
	default:
		return "UNKNOWN_ERROR"
	}
}

