//go:build tinygo || (js && wasm)

//nolint // WASM环境需要使用unsafe.Pointer访问线性内存

package framework

import (
	"unsafe"
)

// ==================== WES Go合约开发框架 ====================
//
// 🌟 **设计理念**：为WES合约开发提供统一的Go语言框架
//
// 🎯 **核心特性**：
// - 基于TinyGo编译到WASM的合约开发支持
// - 统一的宿主函数绑定和封装
// - 标准化的合约接口实现辅助
// - 内置错误处理和类型转换
// - 简化的UTXO操作和事件发出
//
// 📋 **主要组件**：
// - ContractBase: 基础合约结构
// - HostFunctions: 宿主函数绑定
// - Utils: 通用辅助工具
// - Types: 标准数据类型定义
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
	ERROR_NOT_IMPLEMENTED      = 9
	ERROR_PERMISSION_DENIED    = 10
	ERROR_UNKNOWN              = 999
)

// ==================== 基础数据类型 ====================

// Address 地址类型（20字节）
type Address [20]byte

// Hash 哈希类型（32字节）
type Hash [32]byte

// TokenID 代币ID类型
type TokenID string

// Amount 金额类型
type Amount uint64

// ==================== 合约基础结构 ====================

// ContractBase 合约基础结构
// 提供所有WES合约的通用功能和接口实现
type ContractBase struct {
	// 合约元数据
	Name        string
	Symbol      string
	Version     string
	Description string
	Author      string
	License     string

	// 合约配置
	Interfaces []string
	Features   []string
}

// NewContractBase 创建新的合约基础实例
func NewContractBase(name, symbol, version string) *ContractBase {
	return &ContractBase{
		Name:       name,
		Symbol:     symbol,
		Version:    version,
		Interfaces: []string{"IContractBase"},
		Features:   []string{},
	}
}

// AddInterface 添加实现的接口
func (cb *ContractBase) AddInterface(interfaceName string) {
	cb.Interfaces = append(cb.Interfaces, interfaceName)
}

// AddFeature 添加合约特性
func (cb *ContractBase) AddFeature(feature string) {
	cb.Features = append(cb.Features, feature)
}

// ==================== 宿主函数便捷方法 ====================
// 以下方法是对全局宿主函数的便捷包装,允许通过合约实例调用

// GetCaller 获取调用者地址
func (cb *ContractBase) GetCaller() string {
	addr := GetCaller()
	return addr.String()
}

// GetContractAddress 获取当前合约地址
func (cb *ContractBase) GetContractAddress() string {
	addr := GetContractAddress()
	return addr.String()
}

// GetTimestamp 获取当前时间戳
func (cb *ContractBase) GetTimestamp() uint64 {
	return GetTimestamp()
}

// GetBlockHeight 获取当前区块高度
func (cb *ContractBase) GetBlockHeight() uint64 {
	return GetBlockHeight()
}

// SetReturnData 设置返回数据
func (cb *ContractBase) SetReturnData(data []byte) error {
	return SetReturnData(data)
}

// GetState 获取状态
func (cb *ContractBase) GetState(key string) []byte {
	data, err := GetState(key)
	if err != nil {
		return []byte{}
	}
	return data
}

// ⚠️ **已删除**：PutState
// 原因：违背WES架构原则，EUTXO模型无全局状态存储
// 替代：使用 AppendStateOutput 在交易草稿中显式记录状态

// EmitEvent 发出事件
func (cb *ContractBase) EmitEvent(name string, data []byte) error {
	event := NewEvent(name)
	event.Data["payload"] = string(data)
	return EmitEvent(event)
}

// EmitLog 发出日志(简化版,实际应使用专门的日志宿主函数)
func (cb *ContractBase) EmitLog(level, message string) error {
	event := NewEvent("Log")
	event.Data["level"] = level
	event.Data["message"] = message
	return EmitEvent(event)
}

// ==================== P1 HostABI 新增方法 ====================
// 注意：BeginTransaction、SimpleTransfer、SimpleStake 等方法已移除
// 这些方法引用了不存在的类型，且未被示例使用
// 实际开发中应使用 helpers 层的业务语义接口

// ==================== 通用辅助函数 ====================

// GetString 从内存指针构造字符串
//
// nolint // WASM环境需要使用unsafe.Pointer访问线性内存，这是必要的用法
func GetString(ptr uint32, len uint32) string {
	if ptr == 0 || len == 0 {
		return ""
	}
	return string((*[1 << 20]byte)(unsafe.Pointer(uintptr(ptr)))[:len:len]) //nolint:unsafeptr // WASM线性内存访问
}

// GetBytes 从内存指针获取字节数组
//
// nolint // WASM环境需要使用unsafe.Pointer访问线性内存，这是必要的用法
func GetBytes(ptr uint32, len uint32) []byte {
	if ptr == 0 || len == 0 {
		return nil
	}
	return (*[1 << 20]byte)(unsafe.Pointer(uintptr(ptr)))[:len:len] //nolint:unsafeptr // WASM线性内存访问
}

// AllocateString 分配字符串到WASM内存并返回指针和长度
//
// nolint // WASM环境需要使用unsafe.Pointer访问线性内存，这是必要的用法
func AllocateString(s string) (uint32, uint32) {
	if len(s) == 0 {
		return 0, 0
	}
	ptr := Malloc(uint32(len(s)))
	if ptr == 0 {
		return 0, 0
	}
	copy((*[1 << 20]byte)(unsafe.Pointer(uintptr(ptr)))[:len(s)], s) //nolint:unsafeptr // WASM线性内存访问
	return ptr, uint32(len(s))
}

// AllocateBytes 分配字节数组到WASM内存
//
// nolint // WASM环境需要使用unsafe.Pointer访问线性内存，这是必要的用法
func AllocateBytes(data []byte) (uint32, uint32) {
	if len(data) == 0 {
		return 0, 0
	}
	ptr := Malloc(uint32(len(data)))
	if ptr == 0 {
		return 0, 0
	}
	copy((*[1 << 20]byte)(unsafe.Pointer(uintptr(ptr)))[:len(data)], data) //nolint:unsafeptr // WASM线性内存访问
	return ptr, uint32(len(data))
}

// Uint64ToString 将uint64转换为字符串
func Uint64ToString(n uint64) string {
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

// ParseUint64 从字符串解析uint64
func ParseUint64(s string) uint64 {
	var result uint64
	for _, digit := range s {
		if digit >= '0' && digit <= '9' {
			result = result*10 + uint64(digit-'0')
		} else {
			break
		}
	}
	return result
}

// ==================== 地址和哈希处理 ====================

// AddressFromBytes 从字节数组创建地址
func AddressFromBytes(data []byte) Address {
	var addr Address
	if len(data) >= 20 {
		copy(addr[:], data[:20])
	}
	return addr
}

// AddressToBytes 将地址转换为字节数组
func (addr Address) ToBytes() []byte {
	return addr[:]
}

// String 将地址转换为字符串（实现 fmt.Stringer 接口）
func (addr Address) String() string {
	return addr.ToString()
}

// AddressToString 将地址转换为 Base58Check 编码字符串
//
// 🎯 **架构对齐说明**：
//   - 复用宿主 AddressManager.BytesToAddress 实现
//   - 输出标准 Base58Check 格式（符合 pb/transaction.Address 规范）
//   - 避免在 TinyGo 环境重复实现复杂编码逻辑
//
// 📋 **实现方式**：
//   - 调用宿主函数 address_bytes_to_base58
//   - 宿主侧委托给 AddressManager 进行编码
//   - 不重复造轮子，完全复用统一规范实现
func (addr Address) ToString() string {
	// 分配缓冲区（Base58Check 地址最大 34 字符）
	maxLen := uint32(64) // 预留足够空间
	buffer := malloc(maxLen)
	if buffer == 0 {
		// 内存分配失败，回退到 hex 格式
		return addr.ToHexString()
	}

	// 调用宿主函数进行 Base58Check 编码
	addrPtr, _ := AllocateBytes(addr.ToBytes())
	if addrPtr == 0 {
		return addr.ToHexString()
	}

	actualLen := addressBytesToBase58(addrPtr, buffer, maxLen)
	if actualLen == 0 {
		// 编码失败，回退到 hex 格式
		return addr.ToHexString()
	}

	// 读取 Base58 字符串
	base58Bytes := GetBytes(buffer, actualLen)
	return string(base58Bytes)
}

// ToHexString 将地址转换为十六进制字符串（调试用）
//
// 🎯 **用途**：
//   - 仅用于调试和日志输出
//   - 当 Base58Check 编码失败时的后备方案
func (addr Address) ToHexString() string {
	const hexChars = "0123456789abcdef"
	result := make([]byte, 42) // "0x" + 40 hex chars
	result[0] = '0'
	result[1] = 'x'

	for i, b := range addr {
		result[2+i*2] = hexChars[b>>4]
		result[2+i*2+1] = hexChars[b&0xf]
	}

	return string(result)
}

// HashFromBytes 从字节数组创建哈希
func HashFromBytes(data []byte) Hash {
	var hash Hash
	if len(data) >= 32 {
		copy(hash[:], data[:32])
	}
	return hash
}

// ComputeHash 计算数据的哈希值（FNV-1a算法，TinyGo WASM环境下的真实实现）
// 返回32字节的哈希值
func ComputeHash(data []byte) Hash {
	const (
		fnvOffset64 uint64 = 14695981039346656037
		fnvPrime64  uint64 = 1099511628211
	)
	hash := fnvOffset64
	for _, b := range data {
		hash ^= uint64(b)
		hash *= fnvPrime64
	}
	
	// 将64位哈希扩展到32字节（通过多次哈希和组合）
	var result Hash
	hash1 := hash
	hash2 := hash * fnvPrime64
	hash3 := hash2 * fnvPrime64
	hash4 := hash3 * fnvPrime64
	
	for i := 0; i < 8; i++ {
		result[i] = byte(hash1 >> (i * 8))
		result[i+8] = byte(hash2 >> (i * 8))
		result[i+16] = byte(hash3 >> (i * 8))
		result[i+24] = byte(hash4 >> (i * 8))
	}
	return result
}

// HashToBytes 将哈希转换为字节数组
func (hash Hash) ToBytes() []byte {
	return hash[:]
}

// ==================== JSON辅助函数 ====================

// BuildJSONField 构建JSON字段
func BuildJSONField(key, value string) string {
	return `"` + key + `":"` + value + `"`
}

// BuildJSONObject 构建JSON对象
func BuildJSONObject(fields []string) string {
	result := "{"
	for i, field := range fields {
		if i > 0 {
			result += ","
		}
		result += field
	}
	result += "}"
	return result
}

// BuildJSONArray 构建JSON数组
func BuildJSONArray(items []string) string {
	result := "["
	for i, item := range items {
		if i > 0 {
			result += ","
		}
		result += `"` + item + `"`
	}
	result += "]"
	return result
}

// ==================== 合约参数解析 ====================

// ContractParams 合约调用参数
type ContractParams struct {
	data []byte
}

// NewContractParams 创建参数解析器
func NewContractParams(data []byte) *ContractParams {
	return &ContractParams{data: data}
}

// GetRawData 获取原始数据
func (cp *ContractParams) GetRawData() []byte {
	return cp.data
}

// GetString 获取字符串参数
func (cp *ContractParams) GetString() string {
	return string(cp.data)
}

// ParseJSON 简单的JSON字段提取（字符串值）
func (cp *ContractParams) ParseJSON(key string) string {
	data := string(cp.data)
	keyPattern := `"` + key + `":"`

	startIdx := -1
	for i := 0; i <= len(data)-len(keyPattern); i++ {
		if data[i:i+len(keyPattern)] == keyPattern {
			startIdx = i + len(keyPattern)
			break
		}
	}

	if startIdx == -1 {
		return ""
	}

	endIdx := startIdx
	for endIdx < len(data) && data[endIdx] != '"' {
		endIdx++
	}

	if endIdx > startIdx {
		return data[startIdx:endIdx]
	}

	return ""
}

// MustGetString 获取必需的字符串参数（不存在则 panic）
func (cp *ContractParams) MustGetString(key string) string {
	value := cp.ParseJSON(key)
	if value == "" {
		// 在 WASM 环境中无法 panic，返回空字符串由调用方检查
		return ""
	}
	return value
}

// GetStringOr 获取字符串参数（带默认值）
func (cp *ContractParams) GetStringOr(key, defaultValue string) string {
	value := cp.ParseJSON(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// ParseJSONInt 从 JSON 中提取整数字段
func (cp *ContractParams) ParseJSONInt(key string) uint64 {
	data := string(cp.data)
	// 查找 "key": 或 "key":（数字不带引号）
	keyPattern1 := `"` + key + `":`
	keyPattern2 := `"` + key + `": `

	startIdx := -1

	// 尝试匹配 "key":
	for i := 0; i <= len(data)-len(keyPattern1); i++ {
		if data[i:i+len(keyPattern1)] == keyPattern1 {
			startIdx = i + len(keyPattern1)
			break
		}
	}

	// 尝试匹配 "key": （有空格）
	if startIdx == -1 {
		for i := 0; i <= len(data)-len(keyPattern2); i++ {
			if data[i:i+len(keyPattern2)] == keyPattern2 {
				startIdx = i + len(keyPattern2)
				break
			}
		}
	}

	if startIdx == -1 {
		return 0
	}

	// 跳过可能的空格
	for startIdx < len(data) && data[startIdx] == ' ' {
		startIdx++
	}

	// 解析数字
	var result uint64
	for i := startIdx; i < len(data); i++ {
		c := data[i]
		if c >= '0' && c <= '9' {
			result = result*10 + uint64(c-'0')
		} else {
			break
		}
	}

	return result
}

// GetIntOr 获取整数参数（带默认值）
func (cp *ContractParams) GetIntOr(key string, defaultValue uint64) uint64 {
	value := cp.ParseJSONInt(key)
	if value == 0 {
		// 注意：无法区分"不存在"和"值为0"，调用方需注意
		return defaultValue
	}
	return value
}

// IsEmpty 检查参数是否为空
func (cp *ContractParams) IsEmpty() bool {
	return len(cp.data) == 0
}

// ==================== 地址解析工具 ====================

// ParseAddressBase58 从 Base58Check 编码字符串解析地址
//
// 🎯 **架构对齐说明**：
//   - 复用宿主 AddressManager.AddressToBytes 实现
//   - 支持标准 Base58Check 格式（符合 pb/transaction.Address 规范）
//   - 避免在 TinyGo 环境重复实现复杂解码逻辑
//
// 📋 **输入格式**：
//   - "Cf1Kes6snEUeykiJJgrAtKPNPrAzPdPmSn" -> Address{20字节}
func ParseAddressBase58(base58Str string) (Address, error) {
	if base58Str == "" {
		return Address{}, NewContractError(ERROR_INVALID_PARAMS, "address string cannot be empty")
	}

	// 分配结果缓冲区（20 字节）
	resultPtr := malloc(20)
	if resultPtr == 0 {
		return Address{}, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate memory for address")
	}

	// 调用宿主函数进行 Base58Check 解码
	base58Ptr, base58Len := AllocateString(base58Str)
	if base58Ptr == 0 {
		return Address{}, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate memory for base58 string")
	}

	success := addressBase58ToBytes(base58Ptr, base58Len, resultPtr)
	if success == 0 {
		return Address{}, NewContractError(ERROR_INVALID_PARAMS, "invalid base58 address format")
	}

	// 读取 20 字节地址
	addressBytes := GetBytes(resultPtr, 20)
	return AddressFromBytes(addressBytes), nil
}

// ⚠️ **已删除**：ParseAddressFromHex 和 hexCharToNibble
// 原因：不符合统一地址规范（应使用 Base58Check）
// 替代：使用 ParseAddressBase58

// ==================== 错误处理 ====================

// ContractError 合约错误类型
type ContractError struct {
	Code    uint32
	Message string
}

// Error 实现error接口
func (ce *ContractError) Error() string {
	return ce.Message
}

// NewContractError 创建新的合约错误
func NewContractError(code uint32, message string) *ContractError {
	return &ContractError{
		Code:    code,
		Message: message,
	}
}

// WrapError 封装错误为合约错误
func WrapError(code uint32, err error) *ContractError {
	if err == nil {
		return nil
	}
	return &ContractError{
		Code:    code,
		Message: err.Error(),
	}
}

// ==================== 事件辅助 ====================

// Event 事件结构
type Event struct {
	Name string
	Data map[string]interface{}
}

// NewEvent 创建新事件
func NewEvent(name string) *Event {
	return &Event{
		Name: name,
		Data: make(map[string]interface{}),
	}
}

// AddField 添加事件字段
func (e *Event) AddField(key string, value interface{}) {
	e.Data[key] = value
}

// AddStringField 添加字符串字段
func (e *Event) AddStringField(key, value string) {
	e.Data[key] = value
}

// AddUint64Field 添加数值字段
func (e *Event) AddUint64Field(key string, value uint64) {
	e.Data[key] = value
}

// AddAddressField 添加地址字段
func (e *Event) AddAddressField(key string, addr Address) {
	e.Data[key] = addr.ToString()
}

// AddBytesField 添加字节数组字段（Base64编码）
func (e *Event) AddBytesField(key string, value []byte) {
	// 将字节数组转换为Base64编码的字符串
	// 简化实现：使用十六进制编码
	const hexChars = "0123456789abcdef"
	result := ""
	for _, b := range value {
		result += string(hexChars[b>>4])
		result += string(hexChars[b&0x0F])
	}
	e.Data[key] = "0x" + result
}

// AddBoolField 添加布尔字段
func (e *Event) AddBoolField(key string, value bool) {
	if value {
		e.Data[key] = "true"
	} else {
		e.Data[key] = "false"
	}
}

// ToJSON 转换为JSON字符串（简化实现）
func (e *Event) ToJSON() string {
	fields := []string{
		BuildJSONField("event", e.Name),
		BuildJSONField("timestamp", Uint64ToString(GetTimestamp())),
	}

	// 添加数据字段（简化实现）
	dataFields := []string{}
	for key, value := range e.Data {
		switch v := value.(type) {
		case string:
			dataFields = append(dataFields, BuildJSONField(key, v))
		case uint64:
			dataFields = append(dataFields, BuildJSONField(key, Uint64ToString(v)))
		}
	}

	if len(dataFields) > 0 {
		fields = append(fields, `"data":`+BuildJSONObject(dataFields))
	}

	return BuildJSONObject(fields)
}

// ==================== 元数据辅助 ====================

// BuildMetadataJSON 构建合约元数据JSON
func (cb *ContractBase) BuildMetadataJSON() string {
	fields := []string{
		BuildJSONField("name", cb.Name),
		BuildJSONField("symbol", cb.Symbol),
		BuildJSONField("version", cb.Version),
		BuildJSONField("description", cb.Description),
		BuildJSONField("author", cb.Author),
		BuildJSONField("license", cb.License),
	}

	if len(cb.Interfaces) > 0 {
		fields = append(fields, `"interfaces":`+BuildJSONArray(cb.Interfaces))
	}

	if len(cb.Features) > 0 {
		fields = append(fields, `"features":`+BuildJSONArray(cb.Features))
	}

	return BuildJSONObject(fields)
}
