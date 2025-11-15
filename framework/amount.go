//go:build tinygo || (js && wasm)

package framework

import "strconv"

const weiPerUnit uint64 = 1_000_000_000

// FormatWeiToDecimal 将最小单位（wei）格式化为 WES 字符串表示（最多 9 位小数）。
func FormatWeiToDecimal(amountWei uint64) string {
	integerPart := amountWei / weiPerUnit
	fractionPart := amountWei % weiPerUnit

	if fractionPart == 0 {
		return strconv.FormatUint(integerPart, 10)
	}

	// 补齐 9 位小数，再去掉尾随的 0
	fractionWithPadding := strconv.FormatUint(fractionPart+weiPerUnit, 10)[1:]
	trimmedLen := len(fractionWithPadding)
	for trimmedLen > 0 && fractionWithPadding[trimmedLen-1] == '0' {
		trimmedLen--
	}

	if trimmedLen == 0 {
		return strconv.FormatUint(integerPart, 10)
	}

	return strconv.FormatUint(integerPart, 10) + "." + fractionWithPadding[:trimmedLen]
}

// BuildBalanceResult 生成标准余额返回结构，包含原始 wei 值与格式化字符串。
func BuildBalanceResult(address string, tokenID string, balanceWei uint64) map[string]interface{} {
	return map[string]interface{}{
		"address":     address,
		"token_id":    tokenID,
        "balance":     FormatWeiToDecimal(balanceWei),
		"balance_wei": balanceWei,
	}
}

