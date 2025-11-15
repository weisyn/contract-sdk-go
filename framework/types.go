//go:build tinygo || (js && wasm)

package framework

// ==================== 类型定义 ====================
//
// 🌟 **设计理念**：提供类型安全的类型系统
//
// 🎯 **核心特性**：
// - 类型安全：编译期类型检查
// - 零运行时成本：类型别名，无额外开销
// - 语义清晰：类型名称表达业务含义

// ==================== 类型定义扩展 ====================
//
// 注意：基础类型（Address、Hash、TokenID、Amount）定义在contract_base.go中
// 本文件提供交易和资源相关的扩展类型定义

// ==================== 交易相关类型 ====================

// OutPoint UTXO引用点
//
// **用途**：唯一标识一个UTXO
//
// **组成**：
//   - TxHash: 交易哈希（32字节）
//   - Index: 输出索引（uint32）
type OutPoint struct {
	TxHash []byte // 32字节交易哈希
	Index  uint32 // 输出索引
}

// TxOutput 交易输出（简化版）
//
// **用途**：表示交易输出的基本信息
//
// **类型**：
//   - "asset": 资产输出
//   - "resource": 资源输出
//   - "state": 状态输出
type TxOutput struct {
	Type      string  // "asset" | "resource" | "state"
	Recipient Address // 接收者地址（仅asset类型）
	Amount    Amount  // 金额（仅asset类型）
	TokenID   TokenID // 代币ID（仅asset类型）
	Data      []byte  // 其他数据
}

// UTXO 未花费交易输出
//
// **用途**：表示一个UTXO的完整信息
type UTXO struct {
	OutPoint OutPoint
	Output   TxOutput
}

// ==================== 锁定相关类型 ====================

// LockingCondition 锁定条件
//
// **用途**：定义UTXO的解锁条件
//
// **类型**：
//   - "singleKey": 单密钥锁定
//   - "timeLock": 时间锁
//   - "heightLock": 高度锁
//   - "contractLock": 合约锁定
//   - "multiKey": 多密钥锁定
//   - "thresholdLock": 阈值锁定
//
// **格式**：Condition字段为JSON编码的锁定条件
type LockingCondition struct {
	Type      string // 锁定类型
	Condition []byte // 条件数据（JSON编码）
}

// UnlockingProof 解锁证明
//
// **用途**：提供解锁UTXO的证明
//
// **类型**：
//   - "signature": 签名证明
//   - "contract": 合约证明
//   - "multiSig": 多重签名证明
type UnlockingProof struct {
	Type      string // 证明类型
	ProofData []byte // 证明数据
}

// ==================== 资源相关类型 ====================

// Resource 资源元数据
//
// **用途**：表示资源的元数据信息
type Resource struct {
	ContentHash []byte // 32字节内容哈希
	Category    string // "static" | "executable"
	MimeType    string // MIME类型
	Size        uint64 // 资源大小（字节）
}

// ==================== 受控外部交互相关类型（ISPC创新）====================

// ExternalStateClaim 外部状态声明
//
// **用途**：声明外部数据的预期状态
//
// **ISPC创新**：通过"受控声明+佐证+验证"机制，替代传统预言机
type ExternalStateClaim struct {
	ClaimType      string                 // "api_response" | "database_query" | "file_content"
	Source         string                 // API端点/数据库标识/文件标识
	QueryParams    map[string]interface{} // 查询参数
	Timestamp      uint64                 // 时间戳
	ExpectedResponse map[string]interface{} // 预期响应数据
	ClaimID        []byte                 // 声明ID（由系统生成）
}

// Evidence 验证佐证
//
// **用途**：提供可密码学验证的佐证数据
type Evidence struct {
	ClaimID        []byte // 关联的声明ID
	APISignature   []byte // API数字签名
	ResponseHash   []byte // 响应数据哈希
	TimestampProof []byte // 时间戳证明
	DataIntegrity  []byte // 数据完整性证明（如Merkle证明）
	Attestation    []byte // 第三方验证者签名
}

// ResourceCategory 资源类别
const (
	ResourceCategoryStatic     = "static"     // 静态资源（文件/数据）
	ResourceCategoryExecutable = "executable" // 可执行资源（WASM/ONNX）
)

// ==================== 事件相关类型 ====================
//
// 注意：Event类型定义在contract_base.go中
// 此处仅提供类型说明文档

// ==================== 类型转换方法 ====================

// ==================== Amount类型扩展方法 ====================
//
// 注意：Address和Hash的基础方法（ToBytes、String等）定义在contract_base.go中
// 此处提供Amount类型的扩展方法

// String 将金额转换为字符串
func (amount Amount) String() string {
	return Uint64ToString(uint64(amount))
}

// Add 金额相加
func (amount Amount) Add(other Amount) Amount {
	return amount + other
}

// Sub 金额相减
func (amount Amount) Sub(other Amount) Amount {
	if amount < other {
		return 0
	}
	return amount - other
}

// Mul 金额相乘
func (amount Amount) Mul(multiplier uint64) Amount {
	return Amount(uint64(amount) * multiplier)
}

// Div 金额相除
func (amount Amount) Div(divisor uint64) Amount {
	if divisor == 0 {
		return 0
	}
	return Amount(uint64(amount) / divisor)
}

// Cmp 金额比较
//
// **返回**：
//   - -1: amount < other
//   - 0: amount == other
//   - 1: amount > other
func (amount Amount) Cmp(other Amount) int {
	if amount < other {
		return -1
	}
	if amount > other {
		return 1
	}
	return 0
}

