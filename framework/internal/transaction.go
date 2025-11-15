//go:build tinygo || (js && wasm)

package internal

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// ==================== WES 合约交易构建器（链式API）====================
//
// ⚠️ **内部包**：此包仅供 helpers 层使用，外部开发者不应导入
//
// 🌟 **设计理念**：为合约开发提供 Rust-like 链式 API
//
// 🎯 **核心特性**：
// - 链式调用，代码简洁
// - 类型安全，编译检查
// - 确定性保证
// - 与 P1 HostABI 完整集成

// TransactionDraft 交易草稿（SDK侧）
type TransactionDraft struct {
	inputs  []InputDescriptor  // 交易输入
	outputs []OutputDescriptor // 交易输出
	intents []IntentDescriptor // 转账意图（用于账户抽象）
}

// OutputDescriptor 输出描述符
type OutputDescriptor struct {
	outputType string
	to         []byte
	tokenID    []byte
	amount     uint64
	stateID    []byte
	stateVer   uint64
	execHash   []byte
	resource   []byte // 资源序列化数据
}

// InputDescriptor 输入描述符
type InputDescriptor struct {
	outpoint        framework.OutPoint
	isReferenceOnly bool
	unlockingProof  framework.UnlockingProof
}

// IntentDescriptor 意图描述符
type IntentDescriptor struct {
	intentType string
	from       []byte
	to         []byte
	tokenID    []byte
	amount     uint64
	validator  []byte
}

// TransactionBuilder 交易构建器（链式API）
type TransactionBuilder struct {
	draft *TransactionDraft
	err   error
}

// BeginTransaction 开始交易构建
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func BeginTransaction() *TransactionBuilder {
	return &TransactionBuilder{
		draft: &TransactionDraft{
			inputs:  make([]InputDescriptor, 0),
			outputs: make([]OutputDescriptor, 0),
			intents: make([]IntentDescriptor, 0),
		},
		err: nil,
	}
}

// AddAssetOutput 添加资产输出
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func (tb *TransactionBuilder) AddAssetOutput(to framework.Address, tokenID framework.TokenID, amount framework.Amount) *TransactionBuilder {
	if tb.err != nil {
		return tb
	}

	tb.draft.outputs = append(tb.draft.outputs, OutputDescriptor{
		outputType: "asset",
		to:         to.ToBytes(),
		tokenID:    []byte(tokenID),
		amount:     uint64(amount),
	})

	return tb
}

// AddResourceOutput 添加资源输出
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func (tb *TransactionBuilder) AddResourceOutput(resource []byte) *TransactionBuilder {
	if tb.err != nil {
		return tb
	}

	tb.draft.outputs = append(tb.draft.outputs, OutputDescriptor{
		outputType: "resource",
		resource:   resource,
	})

	return tb
}

// AddStateOutput 添加状态输出
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func (tb *TransactionBuilder) AddStateOutput(stateID []byte, version uint64, execHash []byte) *TransactionBuilder {
	if tb.err != nil {
		return tb
	}

	tb.draft.outputs = append(tb.draft.outputs, OutputDescriptor{
		outputType: "state",
		stateID:    stateID,
		stateVer:   version,
		execHash:   execHash,
	})

	return tb
}

// AddInput 添加交易输入
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func (tb *TransactionBuilder) AddInput(outpoint framework.OutPoint, isReferenceOnly bool, unlockingProof framework.UnlockingProof) *TransactionBuilder {
	if tb.err != nil {
		return tb
	}

	// 验证outpoint
	if len(outpoint.TxHash) != 32 {
		tb.err = framework.NewContractError(framework.ERROR_INVALID_PARAMS, "txHash must be 32 bytes")
		return tb
	}

	tb.draft.inputs = append(tb.draft.inputs, InputDescriptor{
		outpoint:        outpoint,
		isReferenceOnly: isReferenceOnly,
		unlockingProof:  unlockingProof,
	})

	return tb
}

// Transfer 添加转账意图（链式API）
//
// ⚠️ **内部接口**：仅供 helpers 层使用
//
// 🎯 **完整业务逻辑实现**：
// 1. 此方法仅添加转账意图到draft中
// 2. 当调用 Finalize() 时，会调用 host_build_transaction
// 3. host_build_transaction 会处理 intent，调用 txAdapter.AddTransfer
// 4. AddTransfer 会：
//    - 使用 Selector 选择 UTXO（基于 from 地址和 tokenID）
//    - 添加选中的 UTXO 作为输入
//    - 添加转账输出（to 地址）
//    - 计算并添加找零输出（如果有）
//
// ✅ **架构优势**：
// - SDK层只关心业务意图（谁给谁转多少钱）
// - Host层处理技术细节（UTXO选择、找零计算）
// - 符合WES"无业务语义"架构原则
func (tb *TransactionBuilder) Transfer(from, to framework.Address, tokenID framework.TokenID, amount framework.Amount) *TransactionBuilder {
	if tb.err != nil {
		return tb
	}

	tb.draft.intents = append(tb.draft.intents, IntentDescriptor{
		intentType: "transfer",
		from:       from.ToBytes(),
		to:         to.ToBytes(),
		tokenID:    []byte(tokenID),
		amount:     uint64(amount),
	})

	return tb
}

// Stake 添加质押意图（链式API）
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func (tb *TransactionBuilder) Stake(staker framework.Address, amount framework.Amount, validator framework.Address) *TransactionBuilder {
	if tb.err != nil {
		return tb
	}

	tb.draft.intents = append(tb.draft.intents, IntentDescriptor{
		intentType: "stake",
		from:       staker.ToBytes(),
		amount:     uint64(amount),
		validator:  validator.ToBytes(),
	})

	return tb
}

// WithFee 设置费用偏好（可选）
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func (tb *TransactionBuilder) WithFee(feeAmount framework.Amount) *TransactionBuilder {
	if tb.err != nil {
		return tb
	}

	// TODO: 在draft中添加费用偏好字段
	// 目前确定性模式会自动计算费用

	return tb
}

// Finalize 完成交易构建
//
// ⚠️ **内部接口**：仅供 helpers 层使用
//
// 🔄 **更新说明**：
//   - 使用新的 host_build_transaction 签名（4个参数）
//   - 返回 TxReceipt JSON，从中提取交易哈希
func (tb *TransactionBuilder) Finalize() (bool, []byte, uint32) {
	if tb.err != nil {
		return false, nil, framework.ERROR_EXECUTION_FAILED
	}

	// 序列化draft为JSON（添加 sign_mode 字段）
	draftJSON := tb.serializeDraft()
	if draftJSON == "" {
		return false, nil, framework.ERROR_EXECUTION_FAILED
	}

	// 调用宿主函数构建交易
	draftPtr, draftLen := framework.AllocateString(draftJSON)
	if draftPtr == 0 {
		return false, nil, framework.ERROR_EXECUTION_FAILED
	}

	// 分配 receipt 缓冲区（足够大以容纳 JSON 响应）
	receiptSize := uint32(4096) // 4KB 应该足够
	receiptPtr := framework.Malloc(receiptSize)
	if receiptPtr == 0 {
		return false, nil, framework.ERROR_EXECUTION_FAILED
	}

	// 调用宿主函数（新版本：4个参数）
	result := hostBuildTransaction(draftPtr, draftLen, receiptPtr, receiptSize)
	if result != framework.SUCCESS {
		return false, nil, result
	}

	// 读取 receipt JSON
	// 注意：需要找到实际的 JSON 结束位置，而不是使用整个缓冲区
	receiptBytes := framework.GetBytes(receiptPtr, receiptSize)
	if receiptBytes == nil || len(receiptBytes) == 0 {
		return false, nil, framework.ERROR_EXECUTION_FAILED
	}
	
	// 找到实际的 JSON 结束位置（查找最后一个 '}'）
	actualLen := findJSONEnd(receiptBytes)
	if actualLen == 0 {
		return false, nil, framework.ERROR_EXECUTION_FAILED
	}
	
	receiptJSON := string(receiptBytes[:actualLen])

	// 解析 receipt JSON 提取交易哈希
	txHash, errCode := parseTxHashFromReceipt(receiptJSON)
	if txHash == nil {
		return false, nil, errCode
	}

	return true, txHash, framework.SUCCESS
}

// parseTxHashFromReceipt 从 TxReceipt JSON 中解析交易哈希
//
// TxReceipt 结构：
//   {
//     "mode": "unsigned" | "delegated" | "threshold" | "paymaster",
//     "unsigned_tx_hash": "...",  // defer_sign/delegated/threshold/paymaster 模式
//     "signed_tx_hash": "...",     // 其他模式
//     "error": "..."               // 错误信息（如果失败）
//   }
func parseTxHashFromReceipt(receiptJSON string) ([]byte, uint32) {
	// 简单的 JSON 解析（TinyGo 环境）
	// 查找 "unsigned_tx_hash" 或 "signed_tx_hash" 字段
	
	// 先检查是否有错误
	if contains(receiptJSON, `"error"`) && !contains(receiptJSON, `"error":""`) && !contains(receiptJSON, `"error":null`) {
		return nil, framework.ERROR_EXECUTION_FAILED
	}

	// 尝试提取 unsigned_tx_hash
	if idx := indexOf(receiptJSON, `"unsigned_tx_hash":"`); idx >= 0 {
		start := idx + len(`"unsigned_tx_hash":"`)
		end := start
		for end < len(receiptJSON) && receiptJSON[end] != '"' {
			end++
		}
		if end > start {
			hashHex := receiptJSON[start:end]
			return hexDecode(hashHex), framework.SUCCESS
		}
	}

	// 尝试提取 signed_tx_hash
	if idx := indexOf(receiptJSON, `"signed_tx_hash":"`); idx >= 0 {
		start := idx + len(`"signed_tx_hash":"`)
		end := start
		for end < len(receiptJSON) && receiptJSON[end] != '"' {
			end++
		}
		if end > start {
			hashHex := receiptJSON[start:end]
			return hexDecode(hashHex), framework.SUCCESS
		}
	}

	// 如果都没有找到，返回错误
	return nil, framework.ERROR_EXECUTION_FAILED
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return indexOf(s, substr) >= 0
}

// indexOf 查找子串位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// hexDecode 解码十六进制字符串（简化实现）
func hexDecode(hexStr string) []byte {
	// 移除 0x 前缀
	if len(hexStr) >= 2 && hexStr[0:2] == "0x" {
		hexStr = hexStr[2:]
	}

	// 确保长度为偶数
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	result := make([]byte, len(hexStr)/2)
	for i := 0; i < len(hexStr); i += 2 {
		high := hexCharToByte(hexStr[i])
		low := hexCharToByte(hexStr[i+1])
		result[i/2] = (high << 4) | low
	}
	return result
}

// hexCharToByte 将十六进制字符转换为字节
func hexCharToByte(c byte) byte {
	if c >= '0' && c <= '9' {
		return c - '0'
	}
	if c >= 'a' && c <= 'f' {
		return c - 'a' + 10
	}
	if c >= 'A' && c <= 'F' {
		return c - 'A' + 10
	}
	return 0
}

// findJSONEnd 找到 JSON 的实际结束位置
// 通过查找最后一个 '}' 来确定 JSON 的实际长度
func findJSONEnd(data []byte) int {
	// 从后往前查找最后一个 '}'
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == '}' {
			// 检查前面是否有非空白字符
			// 简单实现：返回位置+1（包含 '}'）
			return i + 1
		}
		// 跳过空白字符
		if data[i] != ' ' && data[i] != '\t' && data[i] != '\n' && data[i] != '\r' && data[i] != 0 {
			// 遇到非空白字符但没有找到 '}'，说明 JSON 不完整
			break
		}
	}
	return 0
}

// serializeDraft 序列化draft为JSON
func (tb *TransactionBuilder) serializeDraft() string {
	json := `{"sign_mode":"defer_sign","outputs":[`

	// 序列化outputs
	for i, out := range tb.draft.outputs {
		if i > 0 {
			json += ","
		}

		json += `{"type":"` + out.outputType + `"`

		switch out.outputType {
		case "asset":
			json += `,"to":"` + base64Encode(out.to) + `"`
			json += `,"token_id":"` + base64Encode(out.tokenID) + `"`
			json += `,"amount":"` + uint64ToStr(out.amount) + `"`

		case "resource":
			json += `,"resource":"` + base64Encode(out.resource) + `"`

		case "state":
			json += `,"state_id":"` + base64Encode(out.stateID) + `"`
			json += `,"version":` + uint64ToStr(out.stateVer)
			json += `,"exec_hash":"` + base64Encode(out.execHash) + `"`
		}

		json += "}"
	}

	json += `],"intents":[`

	// 序列化intents
	for i, intent := range tb.draft.intents {
		if i > 0 {
			json += ","
		}

		json += `{"type":"` + intent.intentType + `"`

		switch intent.intentType {
		case "transfer":
			json += `,"from":"` + base64Encode(intent.from) + `"`
			json += `,"to":"` + base64Encode(intent.to) + `"`
			json += `,"token_id":"` + base64Encode(intent.tokenID) + `"`
			json += `,"amount":"` + uint64ToStr(intent.amount) + `"`

		case "stake":
			json += `,"staker":"` + base64Encode(intent.from) + `"`
			json += `,"amount":"` + uint64ToStr(intent.amount) + `"`
			json += `,"validator":"` + base64Encode(intent.validator) + `"`
		}

		json += "}"
	}

	json += "]}"

	return json
}

// ==================== 宿主函数声明 ====================

// hostBuildTransaction 构建交易（宿主函数）
//
// 🔄 **更新说明**：
//   - 新版本签名：4个参数（draftPtr, draftLen, receiptPtr, receiptSize）
//   - 返回 TxReceipt JSON 到 receiptPtr，而不是交易哈希
//   - receiptSize 是 receipt 缓冲区的最大容量
//
// 📋 **参数**：
//   - draftPtr: Draft JSON 指针（在 WASM 内存中）
//   - draftLen: Draft JSON 长度
//   - receiptPtr: TxReceipt JSON 写入指针（在 WASM 内存中）
//   - receiptSize: TxReceipt 缓冲区大小
//
// 🔧 **返回值**：
//   - 0: 成功
//   - 其他: 错误代码
//
//go:wasmimport env host_build_transaction
func hostBuildTransaction(draftPtr uint32, draftLen uint32, receiptPtr uint32, receiptSize uint32) uint32

