//go:build tinygo || (js && wasm)

package internal

// ==================== 公共辅助函数（内部包）====================
//
// ⚠️ **内部包**：此包仅供 helpers 层使用，外部开发者不应导入

// base64Encode base64编码（SDK内部使用）
func base64Encode(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	result := ""

	for i := 0; i < len(data); i += 3 {
		b1 := data[i]
		b2 := byte(0)
		b3 := byte(0)

		if i+1 < len(data) {
			b2 = data[i+1]
		}
		if i+2 < len(data) {
			b3 = data[i+2]
		}

		result += string(base64Table[(b1>>2)&0x3F])
		result += string(base64Table[((b1&0x03)<<4)|((b2>>4)&0x0F)])

		if i+1 < len(data) {
			result += string(base64Table[((b2&0x0F)<<2)|((b3>>6)&0x03)])
		} else {
			result += "="
		}

		if i+2 < len(data) {
			result += string(base64Table[b3&0x3F])
		} else {
			result += "="
		}
	}

	return result
}

// uint64ToStr uint64转字符串（SDK内部使用）
func uint64ToStr(n uint64) string {
	if n == 0 {
		return "0"
	}

	var buf [20]byte
	i := len(buf)

	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}

	return string(buf[i:])
}

