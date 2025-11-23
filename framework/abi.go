/**
 * ABI JSON 解析和 JSON Payload 构建工具
 *
 * 提供 ABI JSON 解析和合约调用参数编码功能
 * 用于将合约调用参数转换为 WES 节点期望的 JSON payload 格式
 *
 * 与 contract-sdk-js 的 abi.ts 对齐，确保跨语言一致性
 */

package framework

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ABIParameter ABI 方法参数信息
type ABIParameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
	// StructFields 结构体字段定义（当 type 为 struct 或 object 时）
	// 用于嵌套结构体的编码
	StructFields []ABIParameter `json:"structFields,omitempty"`
}

// ABIMethod ABI 方法信息
type ABIMethod struct {
	Name            string         `json:"name"`
	Type            string         `json:"type"` // "read" or "write"
	Parameters      []ABIParameter `json:"parameters"`
	ReturnType      string         `json:"returnType,omitempty"`
	IsReferenceOnly bool           `json:"isReferenceOnly,omitempty"`
}

// ABI ABI JSON 格式
type ABI struct {
	Methods []ABIMethod `json:"methods"`
	Version string      `json:"version,omitempty"`
}

// BuildPayloadOptions JSON Payload 构建选项
type BuildPayloadOptions struct {
	// IncludeFrom 是否包含调用者地址（from）
	IncludeFrom bool
	// From 调用者地址（如果 IncludeFrom 为 true）
	From interface{} // string or []byte
	// IncludeTo 是否包含接收者地址（to）
	IncludeTo bool
	// To 接收者地址（如果 IncludeTo 为 true）
	To interface{} // string or []byte
	// IncludeAmount 是否包含金额（amount）
	IncludeAmount bool
	// Amount 金额值（如果 IncludeAmount 为 true）
	Amount interface{} // string, int, int64, or *big.Int
	// IncludeTokenID 是否包含代币 ID（token_id）
	IncludeTokenID bool
	// TokenID 代币 ID（如果 IncludeTokenID 为 true）
	TokenID interface{} // string or []byte
}

// ParseABI 解析 ABI JSON 字符串
func ParseABI(abiJSON string) (*ABI, error) {
	var abi ABI
	if err := json.Unmarshal([]byte(abiJSON), &abi); err != nil {
		return nil, fmt.Errorf("failed to parse ABI JSON: %w", err)
	}
	return NormalizeABI(&abi), nil
}

// NormalizeABI 规范化 ABI 格式
func NormalizeABI(data interface{}) *ABI {
	// 如果已经是 ABI 格式
	if abi, ok := data.(*ABI); ok {
		if abi.Version == "" {
			abi.Version = "1.0.0"
		}
		return abi
	}

	// 尝试从 map 转换
	if m, ok := data.(map[string]interface{}); ok {
		abi := &ABI{
			Version: "1.0.0",
		}
		if methods, ok := m["methods"].([]interface{}); ok {
			abi.Methods = make([]ABIMethod, 0, len(methods))
			for _, method := range methods {
				if methodMap, ok := method.(map[string]interface{}); ok {
					abiMethod := convertMapToABIMethod(methodMap)
					abi.Methods = append(abi.Methods, abiMethod)
				}
			}
		}
		return abi
	}

	// 默认返回空 ABI
	return &ABI{
		Methods: []ABIMethod{},
		Version: "1.0.0",
	}
}

// convertMapToABIMethod 将 map 转换为 ABIMethod
func convertMapToABIMethod(m map[string]interface{}) ABIMethod {
	method := ABIMethod{
		Type: "write",
	}

	if name, ok := m["name"].(string); ok {
		method.Name = name
	}
	if typ, ok := m["type"].(string); ok {
		method.Type = typ
	}
	if returnType, ok := m["returnType"].(string); ok {
		method.ReturnType = returnType
	}
	if isRef, ok := m["isReferenceOnly"].(bool); ok {
		method.IsReferenceOnly = isRef
	}

	if params, ok := m["parameters"].([]interface{}); ok {
		method.Parameters = make([]ABIParameter, 0, len(params))
		for _, param := range params {
			if paramMap, ok := param.(map[string]interface{}); ok {
				abiParam := ABIParameter{}
				if name, ok := paramMap["name"].(string); ok {
					abiParam.Name = name
				}
				if typ, ok := paramMap["type"].(string); ok {
					abiParam.Type = typ
				}
				if required, ok := paramMap["required"].(bool); ok {
					abiParam.Required = required
				}
				if desc, ok := paramMap["description"].(string); ok {
					abiParam.Description = desc
				}
				// 解析结构体字段定义（如果存在）
				if structFields, ok := paramMap["structFields"].([]interface{}); ok {
					abiParam.StructFields = parseStructFields(structFields)
				}
				method.Parameters = append(method.Parameters, abiParam)
			}
		}
	}

	return method
}

// parseStructFields 解析结构体字段定义（递归支持嵌套结构体）
func parseStructFields(fields []interface{}) []ABIParameter {
	result := make([]ABIParameter, 0, len(fields))
	for _, field := range fields {
		if fieldMap, ok := field.(map[string]interface{}); ok {
			structField := ABIParameter{}
			if name, ok := fieldMap["name"].(string); ok {
				structField.Name = name
			}
			if typ, ok := fieldMap["type"].(string); ok {
				structField.Type = typ
			}
			if required, ok := fieldMap["required"].(bool); ok {
				structField.Required = required
			}
			if desc, ok := fieldMap["description"].(string); ok {
				structField.Description = desc
			}
			// 递归解析嵌套结构体字段
			if nestedFields, ok := fieldMap["structFields"].([]interface{}); ok {
				structField.StructFields = parseStructFields(nestedFields)
			}
			result = append(result, structField)
		}
	}
	return result
}

// FindMethod 从 ABI 中查找方法信息
func FindMethod(abi *ABI, methodName string) *ABIMethod {
	for i := range abi.Methods {
		if abi.Methods[i].Name == methodName {
			return &abi.Methods[i]
		}
	}
	return nil
}

// BuildJSONPayload 构建 JSON Payload（用于 WES 合约调用）
//
// 根据方法签名和参数类型，将参数值转换为 JSON 对象
// WES 节点期望 JSON payload（Base64 编码），而非 Ethereum ABI 编码
func BuildJSONPayload(methodInfo *ABIMethod, args []interface{}, options *BuildPayloadOptions) (map[string]interface{}, error) {
	payload := make(map[string]interface{})

	// 添加选项中的字段
	if options != nil {
		if options.IncludeFrom && options.From != nil {
			addr, err := convertAddress(options.From)
			if err != nil {
				return nil, fmt.Errorf("invalid from address: %w", err)
			}
			payload["from"] = addr
		}

		if options.IncludeTo && options.To != nil {
			addr, err := convertAddress(options.To)
			if err != nil {
				return nil, fmt.Errorf("invalid to address: %w", err)
			}
			payload["to"] = addr
		}

		if options.IncludeAmount && options.Amount != nil {
			amount, err := convertAmount(options.Amount)
			if err != nil {
				return nil, fmt.Errorf("invalid amount: %w", err)
			}
			payload["amount"] = amount
		}

		if options.IncludeTokenID && options.TokenID != nil {
			tokenID, err := convertTokenID(options.TokenID)
			if err != nil {
				return nil, fmt.Errorf("invalid token ID: %w", err)
			}
			payload["token_id"] = tokenID
		}
	}

	// 添加方法参数
	if methodInfo != nil && len(methodInfo.Parameters) > 0 {
		// 如果有方法信息，使用参数类型进行转换
		for i, param := range methodInfo.Parameters {
			if i < len(args) {
				// 传递结构体字段定义（用于嵌套结构体）
				value, err := convertValueToJSONType(args[i], param.Type, param.StructFields)
				if err != nil {
					return nil, fmt.Errorf("failed to convert parameter %s: %w", param.Name, err)
				}
				payload[param.Name] = value
			}
		}
	} else {
		// 如果没有方法信息，使用类型推断
		for i, arg := range args {
			value := inferJSONType(arg)
			payload[fmt.Sprintf("arg%d", i)] = value
		}
	}

	return payload, nil
}

// convertValueToJSONType 将值转换为 JSON 类型（根据参数类型）
func convertValueToJSONType(value interface{}, typ string, structFields []ABIParameter) (interface{}, error) {
	normalizedType := strings.ToLower(typ)

	// 检查是否为结构体类型
	if isStructType(normalizedType) {
		return convertStructToJSON(value, structFields)
	}

	// 检查是否为数组类型
	if strings.HasSuffix(normalizedType, "[]") {
		return convertArrayToJSON(value, normalizedType, structFields)
	}

	switch normalizedType {
	case "address":
		return convertAddress(value)
	case "number", "uint256", "int256", "u32", "u64", "i32", "i64":
		return convertAmount(value)
	case "bool", "boolean":
		return convertBool(value), nil
	case "bytes", "bytes32":
		return convertBytes(value)
	case "string":
		return convertString(value), nil
	default:
		return convertString(value), nil
	}
}

// isStructType 判断是否为结构体类型
func isStructType(typ string) bool {
	return typ == "struct" ||
		typ == "object" ||
		strings.HasPrefix(typ, "struct:") ||
		strings.HasPrefix(typ, "object:")
}

// convertStructToJSON 将结构体转换为 JSON 对象
func convertStructToJSON(value interface{}, structFields []ABIParameter) (map[string]interface{}, error) {
	// 如果值已经是 map，直接使用
	if valueMap, ok := value.(map[string]interface{}); ok {
		// 如果有结构体字段定义，递归转换每个字段
		if len(structFields) > 0 {
			result := make(map[string]interface{})
			for _, field := range structFields {
				if fieldValue, exists := valueMap[field.Name]; exists {
					// 递归转换字段值
					converted, err := convertValueToJSONType(fieldValue, field.Type, field.StructFields)
					if err != nil {
						return nil, fmt.Errorf("failed to convert field %s: %w", field.Name, err)
					}
					result[field.Name] = converted
				} else if field.Required {
					// 必需字段缺失，使用 nil（调用方应处理）
					result[field.Name] = nil
				}
			}
			return result, nil
		}
		// 没有字段定义，直接返回 map
		return valueMap, nil
	}

	// 如果是字符串，尝试解析为 JSON
	if strValue, ok := value.(string); ok {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(strValue), &parsed); err == nil {
			return convertStructToJSON(parsed, structFields)
		}
		// 解析失败，返回空 map
	}

	// 如果值不是 map，返回空 map（或抛出错误）
	if len(structFields) > 0 {
		// 有字段定义但值不是 map，返回空 map（调用方应处理）
		return make(map[string]interface{}), nil
	}

	// 默认返回空 map
	return make(map[string]interface{}), nil
}

// convertArrayToJSON 将数组转换为 JSON 数组
func convertArrayToJSON(value interface{}, arrayType string, structFields []ABIParameter) ([]interface{}, error) {
	// 确保值是数组
	valueSlice, ok := value.([]interface{})
	if !ok {
		// 如果不是 []interface{}，尝试转换
		if strValue, ok := value.(string); ok {
			var parsed []interface{}
			if err := json.Unmarshal([]byte(strValue), &parsed); err == nil {
				return convertArrayToJSON(parsed, arrayType, structFields)
			}
		}
		// 无法转换，返回空数组
		return []interface{}{}, nil
	}

	// 提取元素类型（移除 []）
	elementType := strings.TrimSuffix(arrayType, "[]")

	// 转换每个元素
	result := make([]interface{}, 0, len(valueSlice))
	for _, item := range valueSlice {
		converted, err := convertValueToJSONType(item, elementType, structFields)
		if err != nil {
			return nil, fmt.Errorf("failed to convert array element: %w", err)
		}
		result = append(result, converted)
	}

	return result, nil
}

// inferJSONType 推断 JSON 类型（当没有方法信息时）
func inferJSONType(value interface{}) interface{} {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case bool:
		return v
	case []byte:
		return hex.EncodeToString(v)
	case string:
		return v
	case map[string]interface{}:
		// 如果是 map（可能是结构体），直接返回
		return v
	case []interface{}:
		// 如果是数组，递归推断每个元素
		result := make([]interface{}, 0, len(v))
		for _, item := range v {
			result = append(result, inferJSONType(item))
		}
		return result
	default:
		return fmt.Sprintf("%v", v)
	}
}

// convertAddress 转换地址为字符串格式
func convertAddress(address interface{}) (string, error) {
	switch v := address.(type) {
	case string:
		return v, nil
	case []byte:
		return hex.EncodeToString(v), nil
	default:
		return "", fmt.Errorf("invalid address type: %T", address)
	}
}

// convertAmount 转换金额为字符串格式
//nolint:unparam // error 返回值保留用于未来扩展
func convertAmount(amount interface{}) (string, error) {
	switch v := amount.(type) {
	case string:
		return v, nil
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v), nil
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32, float64:
		return fmt.Sprintf("%g", v), nil
	default:
		// 尝试转换为字符串
		return fmt.Sprintf("%v", v), nil
	}
}

// convertTokenID 转换代币 ID 为字符串格式
func convertTokenID(tokenID interface{}) (string, error) {
	switch v := tokenID.(type) {
	case string:
		return v, nil
	case []byte:
		return hex.EncodeToString(v), nil
	default:
		return "", fmt.Errorf("invalid token ID type: %T", tokenID)
	}
}

// convertBytes 转换字节数组为十六进制字符串
func convertBytes(bytes interface{}) (string, error) {
	switch v := bytes.(type) {
	case string:
		// 如果已经是字符串，移除 0x 前缀（如果有）
		if strings.HasPrefix(v, "0x") {
			return v[2:], nil
		}
		return v, nil
	case []byte:
		return hex.EncodeToString(v), nil
	default:
		return "", fmt.Errorf("invalid bytes type: %T", bytes)
	}
}

// convertBool 转换布尔值
func convertBool(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false
		}
		return b
	case int, int8, int16, int32, int64:
		return v != 0
	case uint, uint8, uint16, uint32, uint64:
		return v != 0
	default:
		return false
	}
}

// convertString 转换字符串
func convertString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

// EncodePayload 将 JSON Payload 编码为 Base64 字符串
func EncodePayload(payload map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// BuildAndEncodePayload 构建并编码 JSON Payload（一步完成）
func BuildAndEncodePayload(methodInfo *ABIMethod, args []interface{}, options *BuildPayloadOptions) (string, error) {
	payload, err := BuildJSONPayload(methodInfo, args, options)
	if err != nil {
		return "", err
	}
	return EncodePayload(payload)
}
